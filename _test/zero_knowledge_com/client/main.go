package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/kyber/kyber768"
	"github.com/ezydark/zero_knowledge_com/app"
	"github.com/quic-go/quic-go"
	"golang.org/x/crypto/hkdf"

	"github.com/rs/zerolog/log"
)

// --- 1. Core Cryptographic Structures and Functions ---
// (This section remains unchanged)
type PrivateKeypair struct {
	Classical   *ecdh.PrivateKey
	PostQuantum kem.PrivateKey
}

type PublicKeypair struct {
	Classical   *ecdh.PublicKey
	PostQuantum kem.PublicKey
}

type AppMessage struct {
	Nonce uint64
	Data  []byte
}

func generateHybridKeyPair() (*PrivateKeypair, *PublicKeypair, error) {
	classicalPrivate, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	classicalPublic := classicalPrivate.PublicKey()
	pqScheme := kyber768.Scheme()
	pqPublic, pqPrivate, err := pqScheme.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}
	privKeys := &PrivateKeypair{Classical: classicalPrivate, PostQuantum: pqPrivate}
	pubKeys := &PublicKeypair{Classical: classicalPublic, PostQuantum: pqPublic}
	return privKeys, pubKeys, nil
}

func deriveFinalKey(classicalSecret, postQuantumSecret []byte) ([]byte, error) {
	combinedSecret := append(classicalSecret, postQuantumSecret...)
	hash := sha256.New
	kdf := hkdf.New(hash, combinedSecret, nil, []byte("e2ee-chat-key"))
	key := make([]byte, 32)
	if _, err := io.ReadFull(kdf, key); err != nil {
		return nil, err
	}
	return key, nil
}

func encrypt(key []byte, msg AppMessage) ([]byte, error) {
	plaintext, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(key, ciphertext []byte) (*AppMessage, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, actualCiphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	var msg AppMessage
	if err := json.Unmarshal(plaintext, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return &msg, nil
}

// WHY: Encapsulate the client logic to be called multiple times.
// It takes the shared key and configs as arguments to reuse them.
func runClientLogic(sharedKey []byte, tlsConf *tls.Config, quicConf *quic.Config) {
	conn, err := quic.DialAddr(context.Background(), "localhost:4242", tlsConf, quicConf)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to server: %v", err)
	}
	defer conn.CloseWithError(0, "connection closed normally")

	// WHY: We can check the connection state to see if the session was resumed.
	state := conn.ConnectionState().TLS
	log.Info().Msgf("🚀 Alice has connected to the server. Session Resumed: %t", state.DidResume)

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal().Msgf("Failed to open stream: %v", err)
	}

	var nonceBytes [8]byte
	if _, err := rand.Read(nonceBytes[:]); err != nil {
		log.Fatal().Msgf("Failed to generate nonce: %v", err)
	}
	sentNonce := binary.BigEndian.Uint64(nonceBytes[:])

	appMsg := AppMessage{
		Nonce: sentNonce,
		Data:  []byte("This is a secret message!"),
	}
	log.Info().Msgf("✉️ Alice is sending a secret message with nonce %d", appMsg.Nonce)

	encryptedMessage, err := encrypt(sharedKey, appMsg)
	if err != nil {
		log.Fatal().Msgf("Failed to encrypt message: %v", err)
	}

	_, err = stream.Write(encryptedMessage)
	if err != nil {
		log.Fatal().Msgf("Failed to send message: %v", err)
	}
	// Close the stream for writing
	stream.Close()

	response, err := io.ReadAll(stream)
	if err != nil {
		// This can error if the server closes the stream first, which is fine
		log.Info().Msgf("Could not read response: %v", err)
	} else {
		log.Info().Msgf("📬 Alice received %d bytes back.", len(response))
		decryptedMsg, err := decrypt(sharedKey, response)
		if err != nil {
			log.Fatal().Msgf("Failed to decrypt message: %v", err)
		}
		if decryptedMsg.Nonce != sentNonce {
			log.Fatal().Msgf("FATAL: Replay attack! Nonce mismatch. Sent %d, got %d", sentNonce, decryptedMsg.Nonce)
		}
		log.Info().Msgf("✅ Nonce verified. Decrypted message: '%s'", string(decryptedMsg.Data))
	}
}

func main() {
	app.InitLogger()
	// --- Key Exchange (Done once at the start) ---
	alicePrivate, _, _ := generateHybridKeyPair()
	_, bobPublic, _ := generateHybridKeyPair()
	classicalSecretAlice, _ := alicePrivate.Classical.ECDH(bobPublic.Classical)
	pqScheme := kyber768.Scheme()
	_, postQuantumSecretAlice, _ := pqScheme.Encapsulate(bobPublic.PostQuantum)
	sharedKeyAlice, _ := deriveFinalKey(classicalSecretAlice, postQuantumSecretAlice)
	log.Info().Msg("✅ Key exchange complete. Alice has her shared key.")

	// WHY: Define TLS and QUIC configs once to be reused.
	// The tls.Config holds the session cache, which is critical for resumption.
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"pq-chat-example"},
		ClientSessionCache: tls.NewLRUClientSessionCache(1),
	}
	quicConf := &quic.Config{
		Allow0RTT: true,
	}

	// --- First Connection ---
	log.Info().Msg("\n======================================")
	log.Info().Msg("      Attempting First Connection     ")
	log.Info().Msg("======================================")
	runClientLogic(sharedKeyAlice, tlsConf, quicConf)
	log.Info().Msg("✅ First connection finished.")

	// --- Second Connection ---
	log.Info().Msg("\n======================================")
	log.Info().Msg("  Attempting Second Connection (Resumed) ")
	log.Info().Msg("======================================")
	// WHY: Wait a moment to simulate a real disconnect.
	time.Sleep(1 * time.Second)
	// WHY: Call the logic again with the *same* configs. The session ticket from the
	// first connection is still in tlsConf.ClientSessionCache and will be used.
	runClientLogic(sharedKeyAlice, tlsConf, quicConf)
	log.Info().Msg("✅ Second connection finished.")
}
