package comm

import (
	"context"
	"strconv"

	"github.com/quic-go/quic-go"
	"github.com/rs/zerolog/log"
)

type mainServer struct {
	Host string
	Port int

	QuicConf *quic.Config
}

var serverConfig = &mainServer{
	Host: "localhost",
	Port: 8080,

	QuicConf: &quic.Config{
		Allow0RTT: true,
	},
}

func InitServer() error {
	listener, err := quic.ListenAddr(serverConfig.Host+":"+strconv.Itoa(serverConfig.Port), GenerateTLSConfig(), serverConfig.QuicConf)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Debug().Msg("Server started successfully")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			return err
		}

		go HandleConnection(conn)
	}
}
