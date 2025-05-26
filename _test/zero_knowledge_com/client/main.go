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
	"fmt"
	"io"
	"time"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/kyber/kyber768"
	"github.com/ezydark/zero_knowledge_com/app"
	pb "github.com/ezydark/zero_knowledge_com/protobuf"
	"github.com/quic-go/quic-go"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/hkdf"
	"google.golang.org/protobuf/proto"
)

// --- Core Cryptographic Structures and Functions ---
// (This section remains unchanged)
type PrivateKeypair struct {
	Classical   *ecdh.PrivateKey
	PostQuantum kem.PrivateKey
}

type PublicKeypair struct {
	Classical   *ecdh.PublicKey
	PostQuantum kem.PublicKey
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

// encrypt now takes the InnerMessage and returns the encrypted payload
func encrypt(key []byte, msg *pb.InnerMessage) ([]byte, error) {
	plaintext, err := proto.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal inner message: %w", err)
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

// decrypt now takes the payload and returns the InnerMessage
func decrypt(key, encryptedPayload []byte) (*pb.InnerMessage, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(encryptedPayload) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, actualCiphertext := encryptedPayload[:gcm.NonceSize()], encryptedPayload[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	var msg pb.InnerMessage
	if err := proto.Unmarshal(plaintext, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inner message: %w", err)
	}
	return &msg, nil
}

func runFullHandshakeConnection(sharedKey []byte, tlsConf *tls.Config, quicConf *quic.Config) {
	conn, err := quic.DialAddr(context.Background(), "localhost:4242", tlsConf, quicConf)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to server: %v", err)
	}
	defer conn.CloseWithError(0, "connection closed normally")
	log.Info().Msg("🚀 Alice has connected to the server for the first time (1-RTT).")
	log.Info().Msg("✅ First connection finished, session ticket should now be cached.")
}

func run0RTTLogic(sharedKey []byte, frameToSend *pb.OuterFrame, tlsConf *tls.Config, quicConf *quic.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	payload, err := proto.Marshal(frameToSend)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal outer frame for sending")
		return
	}

	conn, err := quic.DialAddrEarly(ctx, "localhost:4242", tlsConf, quicConf)
	if err != nil {
		log.Error().Err(err).Msg("Failed to dial early")
		return
	}
	defer conn.CloseWithError(0, "connection closed by client")

	<-conn.HandshakeComplete()
	if !conn.ConnectionState().Used0RTT {
		log.Warn().Msg("⚠️ Server did not accept 0-RTT. Aborting test logic for this run.")
		return
	}
	log.Info().Msg("🤝 Handshake complete and 0-RTT was accepted by the server.")

	stream, err := conn.OpenStream()
	if err != nil {
		log.Error().Err(err).Msg("Failed to open 0-RTT stream")
		return
	}

	_, err = stream.Write(payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write to 0-RTT stream")
		stream.CancelWrite(0)
		return
	}
	log.Info().Msgf("✉️  Alice sent frame with nonce %d", frameToSend.GetNonce())
	stream.Close()

	response, err := io.ReadAll(stream)
	if err != nil {
		log.Warn().Err(err).Msg("Could not read response from stream")
		log.Warn().Msg("⚠️  This is expected on the replay attempt if the server rejects the data and closes the connection.")
		return
	}

	var receivedFrame pb.OuterFrame
	if err := proto.Unmarshal(response, &receivedFrame); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal received frame")
		return
	}

	log.Info().Msgf("📬 Alice received frame back with nonce %d.", receivedFrame.GetNonce())
	decryptedMsg, err := decrypt(sharedKey, receivedFrame.GetEncryptedPayload())
	if err != nil {
		log.Error().Err(err).Msg("Failed to decrypt received payload")
	} else {
		log.Info().Msgf("✅ E2EE Payload decrypted successfully: '%s'", string(decryptedMsg.GetData()))
	}
}

func main() {
	app.InitLogger()
	alicePrivate, _, _ := generateHybridKeyPair()
	_, bobPublic, _ := generateHybridKeyPair()
	classicalSecretAlice, _ := alicePrivate.Classical.ECDH(bobPublic.Classical)
	pqScheme := kyber768.Scheme()
	_, postQuantumSecretAlice, _ := pqScheme.Encapsulate(bobPublic.PostQuantum)
	sharedKeyAlice, _ := deriveFinalKey(classicalSecretAlice, postQuantumSecretAlice)
	log.Info().Msg("✅ Key exchange complete. Alice has her shared key.")

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"pq-chat-example"},
		ClientSessionCache: tls.NewLRUClientSessionCache(1),
	}
	quicConf := &quic.Config{
		Allow0RTT: true,
	}

	log.Info().Msg("\n======================================")
	log.Info().Msg("   Step 1: Initial 1-RTT Connection   ")
	log.Info().Msg("======================================")
	runFullHandshakeConnection(sharedKeyAlice, tlsConf, quicConf)
	log.Info().Msg("Waiting a moment before starting 0-RTT attempts...")
	time.Sleep(1 * time.Second)

	log.Info().Msg("\n======================================")
	log.Info().Msg("     Step 2: 0-RTT Replay Attack      ")
	log.Info().Msg("======================================")

	innerMsg := &pb.InnerMessage{}
	innerMsg.SetData([]byte("This is a secret 0-RTT message!"))
	encryptedPayload, err := encrypt(sharedKeyAlice, innerMsg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to encrypt inner message")
	}

	var nonceBytes [8]byte
	rand.Read(nonceBytes[:])
	replayNonce := binary.BigEndian.Uint64(nonceBytes[:])

	frameToSend := &pb.OuterFrame{}
	frameToSend.SetNonce(replayNonce)
	frameToSend.SetEncryptedPayload(encryptedPayload)
	log.Info().Msgf("Crafted a single malicious frame with nonce %d to be replayed.", replayNonce)

	log.Info().Msg("\n--- Attempting LEGITIMATE 0-RTT connection... ---")
	run0RTTLogic(sharedKeyAlice, frameToSend, tlsConf, quicConf)
	log.Info().Msg("✅ Finished legitimate 0-RTT attempt.")

	time.Sleep(1 * time.Second)

	log.Info().Msg("\n--- Attempting REPLAY of the same 0-RTT data... ---")
	run0RTTLogic(sharedKeyAlice, frameToSend, tlsConf, quicConf)
	log.Info().Msg("✅ Finished replay attack attempt.")
}
