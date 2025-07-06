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
	fb "github.com/ezydark/zero_knowledge_com/schema/ezMsgSchema/Payload"
	"github.com/quic-go/quic-go"

	"github.com/rs/zerolog/log"
)

const serverAddr = "localhost:4242"

var nonceStore = struct {
	sync.RWMutex
	m map[uint64]time.Time
}{m: make(map[uint64]time.Time)}

// isNonceSeen checks the nonce store for replays. This is the server's
// primary security function in this model.
func isNonceSeen(nonce uint64) bool {
	nonceStore.Lock()
	defer nonceStore.Unlock()
	if _, ok := nonceStore.m[nonce]; ok {
		return true // Nonce has been seen
	}
	nonceStore.m[nonce] = time.Now() // Record new nonce
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
	log.Info().Msgf("Accepted connection from %s. 0-RTT: %t", conn.RemoteAddr(), conn.ConnectionState().Used0RTT)
	for {
		stream, errAccept := conn.AcceptStream(context.Background())
		if errAccept != nil {
			log.Info().Msgf("Client %s closed the connection: %v", conn.RemoteAddr(), errAccept)
			return
		}

		go func(str quic.Stream) {
			defer str.Close()
			buf, errRead := io.ReadAll(str)
			if errRead != nil {
				log.Error().Err(errRead).Msg("Error reading from stream")
				return
			}

			// Parse OuterFrame from received bytes
			fbOuterFrame := fb.GetRootAsOuterFrame(buf, 0)
			nonce := fbOuterFrame.Nonce()

			// Core replay attack prevention
			if isNonceSeen(nonce) {
				log.Warn().Uint64("nonce", nonce).Msg("🚨 REPLAY ATTACK DETECTED! Rejecting frame.")
				conn.CloseWithError(2, "replay detected") // QUIC error code for replay
				return
			}

			payloadSize := fbOuterFrame.EncryptedPayloadLength()
			log.Info().Uint64("nonce", nonce).Int("payload_size", payloadSize).Msg("✅ Received valid frame, echoing back.")

			// Server echoes the exact same buffer 'buf' back.
			_, errWrite := str.Write(buf)
			if errWrite != nil {
				log.Error().Err(errWrite).Msg("Error writing response to stream")
			}
		}(stream)
	}
}

func main() {
	app.InitLogger()
	log.Info().Msg("✅ E2EE Server starting. I have no knowledge of user keys.")
	if err := runServer(); err != nil {
		log.Fatal().Msgf("Server error: %v", err)
	}
}
