package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ezydark/ezlog/log"
	"github.com/quic-go/quic-go"
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
		Allow0RTT:      true,
		MaxIdleTimeout: 30 * time.Second,
	},
}

func InitServer() error {
	tlsConf, err := GenerateTLSConfig()
	if err != nil {
		return fmt.Errorf("failed to generate TLS config:\n%v", err)
	}
	listener, err := quic.ListenAddr(serverConfig.Host+":"+strconv.Itoa(serverConfig.Port), tlsConf, serverConfig.QuicConf)
	if err != nil {
		return fmt.Errorf("failed to start server:\n%v", err)
	}
	defer listener.Close()
	log.Info().Msg("Server started and listening...")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			return fmt.Errorf("failed to accept connection:\n%v", err)
		}

		go HandleConnection(conn)
	}
}
