package server

import (
	"fmt"

	fb "github.com/ezydark/ezMsg/flatbuffers/generated/ezMsg/Communication"
	"github.com/ezydark/ezlog/log"
	"github.com/quic-go/quic-go"
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
