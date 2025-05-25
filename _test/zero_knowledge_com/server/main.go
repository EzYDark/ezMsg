package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
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

// OuterFrame is the message structure visible to the server.
// It contains the nonce for replay protection and the opaque E2EE payload.
type OuterFrame struct {
	Nonce            uint64
	EncryptedPayload []byte
}

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
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Info().Msgf("Client %s closed the connection: %v", conn.RemoteAddr(), err)
			return
		}

		go func(str quic.Stream) {
			defer str.Close()
			buf, err := io.ReadAll(str)
			if err != nil {
				log.Error().Err(err).Msg("Error reading from stream")
				return
			}

			// WHY: The server unmarshals the outer frame, but has NO KNOWLEDGE
			// of the key needed to decrypt the EncryptedPayload.
			var frame OuterFrame
			if err := json.Unmarshal(buf, &frame); err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal outer frame from client")
				conn.CloseWithError(1, "unmarshal failed")
				return
			}

			// WHY: This is the core of replay attack prevention, now done without
			// compromising E2EE.
			if isNonceSeen(frame.Nonce) {
				log.Warn().Uint64("nonce", frame.Nonce).Msg("🚨 REPLAY ATTACK DETECTED! Rejecting frame.")
				conn.CloseWithError(2, "replay detected")
				return
			}

			log.Info().Uint64("nonce", frame.Nonce).Int("payload_size", len(frame.EncryptedPayload)).Msg("✅ Received valid frame, echoing back.")

			// WHY: The server echoes the exact same buffer back. It cannot construct
			// a new one because it cannot read or re-encrypt the payload.
			_, err = str.Write(buf)
			if err != nil {
				log.Error().Err(err).Msg("Error writing response to stream")
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
