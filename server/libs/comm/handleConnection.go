package comm

import (
	"github.com/quic-go/quic-go"
	"github.com/rs/zerolog/log"
)

func HandleConnection(conn quic.Connection) {
	log.Debug().Msgf("Client connected from %s.", conn.RemoteAddr())
}
