package comm

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/kyber/kyber768"
	"golang.org/x/crypto/hkdf"
)

type PrivateKeypair struct {
	Classical   *ecdh.PrivateKey
	PostQuantum kem.PrivateKey
}

type PublicKeypair struct {
	Classical   *ecdh.PublicKey
	PostQuantum kem.PublicKey
}

func GenerateKeypairs() (*PrivateKeypair, *PublicKeypair, error) {
	classicalPrivate, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate classical key pair:\n%v", err)
	}
	classicalPublic := classicalPrivate.PublicKey()

	pqScheme := kyber768.Scheme()
	pqPublic, pqPrivate, err := pqScheme.GenerateKeyPair()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate post-quantum key pair:\n%v", err)
	}

	privKeys := &PrivateKeypair{Classical: classicalPrivate, PostQuantum: pqPrivate}
	pubKeys := &PublicKeypair{Classical: classicalPublic, PostQuantum: pqPublic}
	return privKeys, pubKeys, nil
}

func DeriveFinalKey(classicalSecret, postQuantumSecret []byte) ([]byte, error) {
	combinedSecret := append(classicalSecret, postQuantumSecret...)

	hash := sha256.New
	kdf := hkdf.New(hash, combinedSecret, nil, []byte("e2ee-chat-key"))
	key := make([]byte, 32)

	_, err := io.ReadFull(kdf, key)
	if err != nil {
		return nil, fmt.Errorf("failed to derive final key:\n%v", err)
	}
	return key, nil
}
