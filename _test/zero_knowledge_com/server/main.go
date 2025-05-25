package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"math/big"
	"sync"
	"time"

	"github.com/ezydark/zero_knowledge_com/app"
	"github.com/quic-go/quic-go"

	"github.com/rs/zerolog/log"
)

const serverAddr = "localhost:4242"

var nonceStore = struct {
	sync.RWMutex
	m map[uint64]bool
}{m: make(map[uint64]bool)}

type AppMessage struct {
	Nonce uint64
	Data  []byte
}

func isNonceSeen(nonce uint64) bool {
	nonceStore.Lock()
	defer nonceStore.Unlock()
	if _, ok := nonceStore.m[nonce]; ok {
		return true
	}
	nonceStore.m[nonce] = true
	return false
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates:           []tls.Certificate{tlsCert},
		NextProtos:             []string{"pq-chat-example"},
		MinVersion:             tls.VersionTLS13,
		SessionTicketsDisabled: false,
	}
}

func runServer() error {
	quicConf := &quic.Config{
		Allow0RTT: true,
	}
	listener, err := quic.ListenAddr(serverAddr, generateTLSConfig(), quicConf)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Info().Msgf("Server listening on %s", serverAddr)

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			return err
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn quic.Connection) {
	log.Info().Msgf("Accepted connection from %s", conn.RemoteAddr())
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			// WHY: Check if this is the "normal closure" error from the client.
			if appErr, ok := err.(*quic.ApplicationError); ok && appErr.ErrorCode == 0 {
				log.Info().Msgf("Client %s closed the connection gracefully.", conn.RemoteAddr())
			} else {
				// It's a different, unexpected error.
				log.Info().Msgf("Error accepting stream: %v", err)
			}
			return // Exit the loop and end the handler for this connection.
		}
		go func(str quic.Stream) {
			defer str.Close()
			buf, err := io.ReadAll(str)
			if err != nil {
				log.Info().Msgf("Error reading from stream: %v", err)
				return
			}
			log.Info().Msgf("Server received %d bytes of encrypted data, echoing back.", len(buf))
			_, err = str.Write(buf)
			if err != nil {
				log.Info().Msgf("Error writing to stream: %v", err)
			}
		}(stream)
	}
}

func main() {
	app.InitLogger()
	if err := runServer(); err != nil {
		log.Fatal().Msgf("Server error: %v", err)
	}
}
