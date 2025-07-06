package comm

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/ezydark/ezMsg/libs/client/flatbuffers/generated/ezMsg/Communication"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/quic-go/quic-go"
	"github.com/rs/zerolog/log"
)

type remoteServer struct {
	Host string
	Port int

	QuicConf *quic.Config
	TlsConf  *tls.Config
}

var serverConfig = &remoteServer{
	Host: "localhost",
	Port: 8080,

	QuicConf: &quic.Config{
		Allow0RTT: true,
	},

	TlsConf: &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"ez-msg-protocol"},
		ClientSessionCache: tls.NewLRUClientSessionCache(1),
	},
}

var Delimiter = []byte("\n\r\n\r")

func InitClient() error {
	_, _, err := GenerateKeypairs()
	if err != nil {
		return fmt.Errorf("failed to generate key pair:\n%v", err)
	}
	// TODO: Derive and get shared secret?

	conn, err := quic.DialAddrEarly(context.Background(),
		serverConfig.Host+":"+strconv.Itoa(serverConfig.Port),
		serverConfig.TlsConf,
		serverConfig.QuicConf,
	)
	if err != nil {
		return fmt.Errorf("failed to connect to server:\n%v", err)
	}
	defer conn.CloseWithError(0, "connection closed by client")

	<-conn.HandshakeComplete() // Wait for handshake
	if !conn.ConnectionState().Used0RTT {
		log.Debug().Msg("Server did not used 0-RTT.")
	}

	stream, err := conn.OpenStream()
	if err != nil {
		return fmt.Errorf("failed to open stream:\n%v", err)
	}
	log.Debug().Msg("Stream opened successfully.")

	builder := flatbuffers.NewBuilder(1024)

	sessionToken := builder.CreateString("dummy-session-token")
	Communication.UnencryptedClientMetadataStart(builder)
	Communication.UnencryptedClientMetadataAddSessionToken(builder, sessionToken)
	Communication.UnencryptedClientMetadataAddNonce(builder, 1234567890) // Used to prevent replay attacks.
	Communication.UnencryptedClientMetadataAddTimestamp(builder, time.Now().Unix())
	metadataOffset := Communication.UnencryptedClientMetadataEnd(builder)

	encryptedContent := builder.CreateByteVector([]byte("this would be encrypted message content"))
	Communication.ChatMessageRequestStart(builder)
	Communication.ChatMessageRequestAddChatUid(builder, 9876543210) // The unique ID of the chat.
	Communication.ChatMessageRequestAddEncryptedContent(builder, encryptedContent)
	payloadOffset := Communication.ChatMessageRequestEnd(builder)

	Communication.ClientFrameStart(builder)
	Communication.ClientFrameAddMetadata(builder, metadataOffset)
	Communication.ClientFrameAddPayloadType(builder, Communication.ClientPayloadChatMessageRequest)
	Communication.ClientFrameAddPayload(builder, payloadOffset)
	clientFrameOffset := Communication.ClientFrameEnd(builder)

	builder.Finish(clientFrameOffset)
	buf := builder.FinishedBytes()

	log.Debug().Msg("Sending client frame...")
	buf = append(buf, Delimiter...)
	_, err = stream.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write to stream:\n%v", err)
	}

	log.Debug().Msg("Client frame sent successfully.")

	scanner := bufio.NewScanner(stream)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if i := bytes.Index(data, Delimiter); i >= 0 {
			return i + len(Delimiter), data[0:i], nil
		}
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		return 0, nil, nil
	})

	if scanner.Scan() {
		responseBytes := scanner.Bytes()
		log.Debug().Msg("Response received successfully.")
		fbFrame := Communication.GetRootAsClientFrame(responseBytes, 0)
		metadataTable := new(Communication.UnencryptedClientMetadata)
		metadata := fbFrame.Metadata(metadataTable)
		log.Info().Msgf("Received frame with nonce '%d' and sessionToken '%s'", metadata.Nonce(), metadata.SessionToken())

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from stream: %w", err)
	}

	return nil
}
