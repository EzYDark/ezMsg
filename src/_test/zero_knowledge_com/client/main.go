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
	fb "github.com/ezydark/zero_knowledge_com/schema/ezMsgSchema/Payload"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/quic-go/quic-go"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/hkdf"
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
func encrypt(key []byte, innerMessageData []byte) ([]byte, error) {
	builder := flatbuffers.NewBuilder(1024) // Initial builder size
	dataOffset := builder.CreateByteVector(innerMessageData)

	fb.InnerMessageStart(builder)
	fb.InnerMessageAddData(builder, dataOffset)
	innerOffset := fb.InnerMessageEnd(builder)
	builder.Finish(innerOffset)
	plaintextFBBytes := builder.FinishedBytes() // Serialized InnerMessage

	// AES-GCM encryption (logic remains the same, operates on plaintextFBBytes)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// Use crypto/rand for nonce generation
	nonce := make([]byte, gcm.NonceSize())
	if _, errReader := rand.Read(nonce); errReader != nil {
		return nil, errReader
	}
	return gcm.Seal(nonce, nonce, plaintextFBBytes, nil), nil
}

// decrypt now takes the payload and returns the InnerMessage
func decrypt(key, encryptedPayloadWithInnerFB []byte) (*fb.InnerMessage, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(encryptedPayloadWithInnerFB) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, actualCiphertext := encryptedPayloadWithInnerFB[:gcm.NonceSize()], encryptedPayloadWithInnerFB[gcm.NonceSize():]
	plaintextFBBytes, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	innerMsgFB := fb.GetRootAsInnerMessage(plaintextFBBytes, 0)
	return innerMsgFB, nil
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

func run0RTTLogic(sharedKey []byte, outerFrameNonce uint64, encryptedInnerPayload []byte, tlsConf *tls.Config, quicConf *quic.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Build OuterFrame Flatbuffer
	builder := flatbuffers.NewBuilder(1024)
	payloadOffset := builder.CreateByteVector(encryptedInnerPayload)

	fb.OuterFrameStart(builder)
	fb.OuterFrameAddNonce(builder, outerFrameNonce)
	fb.OuterFrameAddEncryptedPayload(builder, payloadOffset)
	outerOffset := fb.OuterFrameEnd(builder)
	builder.Finish(outerOffset) // This prepares the buffer for sending
	wirePayload := builder.FinishedBytes()

	conn, errDial := quic.DialAddrEarly(ctx, "localhost:4242", tlsConf, quicConf)
	if errDial != nil {
		log.Error().Err(errDial).Msg("Failed to dial early")
		return
	}
	defer conn.CloseWithError(0, "connection closed by client")

	<-conn.HandshakeComplete() // Wait for handshake
	if !conn.ConnectionState().Used0RTT {
		log.Warn().Msg("⚠️ Server did not accept 0-RTT. Aborting test logic for this run.")
		return
	}
	log.Info().Msg("🤝 Handshake complete and 0-RTT was accepted by the server.")

	stream, errStream := conn.OpenStream()
	if errStream != nil {
		log.Error().Err(errStream).Msg("Failed to open 0-RTT stream")
		return
	}

	_, errWrite := stream.Write(wirePayload)
	if errWrite != nil {
		log.Error().Err(errWrite).Msg("Failed to write to 0-RTT stream")
		stream.CancelWrite(0) // Or handle error appropriately
		return
	}
	log.Info().Msgf("✉️  Alice sent frame with nonce %d", outerFrameNonce)
	stream.Close() // Close the stream for writing

	// Read response from the same stream
	responseBytes, errRead := io.ReadAll(stream)
	if errRead != nil {
		log.Warn().Err(errRead).Msg("Could not read response from stream. Expected on replay if server closes.")
		return
	}

	receivedFbFrame := fb.GetRootAsOuterFrame(responseBytes, 0)
	log.Info().Msgf("📬 Alice received frame back with nonce %d.", receivedFbFrame.Nonce())

	decryptedInnerFbMsg, errDecrypt := decrypt(sharedKey, receivedFbFrame.EncryptedPayloadBytes())
	if errDecrypt != nil {
		log.Error().Err(errDecrypt).Msg("Failed to decrypt received payload")
	} else {
		log.Info().Msgf("✅ E2EE Payload decrypted successfully: '%s'", string(decryptedInnerFbMsg.DataBytes()))
	}
}

func main() {
	app.InitLogger() // Assuming app.InitLogger() is correctly set up
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
	runFullHandshakeConnection(sharedKeyAlice, tlsConf, quicConf) // Assuming this function remains largely the same for establishing session tickets
	log.Info().Msg("Waiting a moment before starting 0-RTT attempts...")
	time.Sleep(1 * time.Second)

	log.Info().Msg("\n======================================")
	log.Info().Msg("     Step 2: 0-RTT Replay Attack      ")
	log.Info().Msg("======================================")

	innerMsgData := []byte("This is a secret 0-RTT message!")
	encryptedPayloadBytes, err := encrypt(sharedKeyAlice, innerMsgData)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to encrypt inner message")
	}

	var nonceBytes [8]byte
	if _, err := rand.Read(nonceBytes[:]); err != nil { // Use crypto/rand
		log.Fatal().Err(err).Msg("Failed to generate nonce bytes for OuterFrame")
	}
	replayNonce := binary.BigEndian.Uint64(nonceBytes[:])

	log.Info().Msgf("Crafted a single malicious frame data with nonce %d to be replayed.", replayNonce)

	log.Info().Msg("\n--- Attempting LEGITIMATE 0-RTT connection... ---")
	run0RTTLogic(sharedKeyAlice, replayNonce, encryptedPayloadBytes, tlsConf, quicConf)
	log.Info().Msg("✅ Finished legitimate 0-RTT attempt.")

	time.Sleep(1 * time.Second)

	log.Info().Msg("\n--- Attempting REPLAY of the same 0-RTT data... ---")
	run0RTTLogic(sharedKeyAlice, replayNonce, encryptedPayloadBytes, tlsConf, quicConf) // Re-use the same data for replay
	log.Info().Msg("✅ Finished replay attack attempt.")
}
