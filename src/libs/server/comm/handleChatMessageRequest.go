package server

import (
	"fmt"

	fb "github.com/ezydark/ezMsg/src/libs/flatbuffers/generated/ezMsg/Communication"
	"github.com/quic-go/quic-go"
	"github.com/rs/zerolog/log"
)

func HandleChatMessageRequest(conn *quic.Conn, stream *quic.Stream, req *fb.ChatMessageRequest) error {
	log.Debug().Msgf("Received chat message from '%s' to chat '%d'", conn.RemoteAddr(), req.ChatUid())

	log.Debug().Msg("Echoing the message back to the client")
	respBuf := req.Table().Bytes

	respBuf = append(respBuf, Delimiter...)

	_, err := stream.Write(respBuf)
	if err != nil {
		return fmt.Errorf("failed to echo message back:\n%v", err)
	}

	log.Debug().Msg("Message echoed successfully.")

	return nil
}
