package comm

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	fb "github.com/ezydark/ezMsg/libs/client/flatbuffers/generated/ezMsg/Communication"
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
		Allow0RTT:      true,
		MaxIdleTimeout: 30 * time.Second,
	},

	TlsConf: &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"ez-msg-protocol"},
		ClientSessionCache: tls.NewLRUClientSessionCache(1),
	},
}

var Delimiter = []byte("\n\r\n\r")

// sendHeartbeats periodically pings the server to keep the connection alive.
func sendHeartbeats(conn *quic.Conn) {
	ticker := time.NewTicker(serverConfig.QuicConf.MaxIdleTimeout - 5*time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stream, err := conn.OpenStream()
		if err != nil {
			log.Error().Err(err).Msg("Failed to open stream for heartbeat")
			return
		}

		_, err = stream.Write([]byte{0})
		if err != nil {
			log.Error().Err(err).Msg("Failed to write heartbeat")
		} else {
			log.Debug().Msg("Heartbeat sent.")
		}
		stream.Close()
	}
}

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

	<-conn.HandshakeComplete()
	log.Debug().Msg("Handshake complete.")

	go sendHeartbeats(conn)

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
	fb.UnencryptedClientMetadataStart(builder)
	fb.UnencryptedClientMetadataAddSessionToken(builder, sessionToken)
	fb.UnencryptedClientMetadataAddNonce(builder, 1234567890) // Used to prevent replay attacks.
	fb.UnencryptedClientMetadataAddTimestamp(builder, time.Now().Unix())
	metadataOffset := fb.UnencryptedClientMetadataEnd(builder)

	encryptedContent := builder.CreateByteVector([]byte("this would be encrypted message content"))
	fb.ChatMessageRequestStart(builder)
	fb.ChatMessageRequestAddChatUid(builder, 9876543210) // The unique ID of the chat.
	fb.ChatMessageRequestAddEncryptedContent(builder, encryptedContent)
	payloadOffset := fb.ChatMessageRequestEnd(builder)

	fb.ClientFrameStart(builder)
	fb.ClientFrameAddMetadata(builder, metadataOffset)
	fb.ClientFrameAddPayloadType(builder, fb.ClientPayloadChatMessageRequest)
	fb.ClientFrameAddPayload(builder, payloadOffset)
	clientFrameOffset := fb.ClientFrameEnd(builder)

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

		// TODO: switch to serverFrame instead
		clientFrame := fb.GetRootAsClientFrame(responseBytes, 0)
		payloadTable := new(flatbuffers.Table)
		if !clientFrame.Payload(payloadTable) {
			return fmt.Errorf("failed to get payload from ClientFrame")
		}

		metadataTable := new(fb.UnencryptedClientMetadata)
		metadata := clientFrame.Metadata(metadataTable)

		log.Info().Msgf("Received frame with nonce '%d' and sessionToken '%s'",
			metadata.Nonce(), metadata.SessionToken())

		switch clientFrame.PayloadType() {
		case fb.ClientPayloadChatMessageRequest:
			req := new(fb.ChatMessageRequest)
			req.Init(payloadTable.Bytes, payloadTable.Pos)
			log.Debug().Msg("[client] Chat message received successfully.")
			// return HandleChatMessageRequest(conn, stream, req)
		default:
			return fmt.Errorf("received unknown frame type:\n%v", clientFrame.PayloadType())
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from stream: %w", err)
	}

	for {
		time.Sleep(time.Second)
	}

	return nil
}
