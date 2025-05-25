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
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/hkdf"
)

// --- Core Cryptographic Structures and Functions (unchanged) ---
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

func runClientLogic(sharedKey []byte, tlsConf *tls.Config, quicConf *quic.Config) {
	conn, err := quic.DialAddr(context.Background(), "localhost:4242", tlsConf, quicConf)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to server")
	}
	defer conn.CloseWithError(0, "connection closed normally")

	state := conn.ConnectionState().TLS
	log.Info().Bool("resumed", state.DidResume).Msg("🚀 Alice has connected to the server")

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open stream")
	}

	var nonceBytes [8]byte
	if _, err := rand.Read(nonceBytes[:]); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate nonce")
	}
	sentNonce := binary.BigEndian.Uint64(nonceBytes[:])

	appMsg := AppMessage{
		Nonce: sentNonce,
		Data:  []byte("This is a secret message!"),
	}
	log.Info().Uint64("nonce", appMsg.Nonce).Msg("✉️ Alice is sending a secret message")

	encryptedMessage, err := encrypt(sharedKey, appMsg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to encrypt message")
	}

	_, err = stream.Write(encryptedMessage)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to send message")
	}
	stream.Close()

	response, err := io.ReadAll(stream)
	if err != nil {
		log.Error().Err(err).Msg("Could not read response")
	} else {
		log.Info().Int("bytes", len(response)).Msg("📬 Alice received data back")
		decryptedMsg, err := decrypt(sharedKey, response)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to decrypt message")
		}
		if decryptedMsg.Nonce != sentNonce {
			log.Fatal().Uint64("sent_nonce", sentNonce).Uint64("received_nonce", decryptedMsg.Nonce).Msg("Replay attack! Nonce mismatch")
		}
		log.Info().Str("message", string(decryptedMsg.Data)).Msg("✅ Nonce verified. Decrypted message successfully")
	}
}

func main() {
	app.InitLogger()

	// Key Exchange
	alicePrivate, _, _ := generateHybridKeyPair()
	_, bobPublic, _ := generateHybridKeyPair()
	classicalSecretAlice, _ := alicePrivate.Classical.ECDH(bobPublic.Classical)
	pqScheme := kyber768.Scheme()
	_, postQuantumSecretAlice, _ := pqScheme.Encapsulate(bobPublic.PostQuantum)
	sharedKeyAlice, err := deriveFinalKey(classicalSecretAlice, postQuantumSecretAlice)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to derive shared key")
	}
	log.Info().Msg("✅ Key exchange complete. Alice has her shared key.")

	// Configs
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"pq-chat-example"},
		ClientSessionCache: tls.NewLRUClientSessionCache(1),
	}
	quicConf := &quic.Config{
		Allow0RTT: true,
	}

	// First Connection
	log.Info().Msg("\n======================================\n      Attempting First Connection     \n======================================")
	runClientLogic(sharedKeyAlice, tlsConf, quicConf)
	log.Info().Msg("✅ First connection finished.")

	// Second Connection
	log.Info().Msg("\n======================================\n  Attempting Second Connection (Resumed) \n======================================")
	time.Sleep(1 * time.Second)
	runClientLogic(sharedKeyAlice, tlsConf, quicConf)
	log.Info().Msg("✅ Second connection finished.")
}
