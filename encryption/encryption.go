package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"hash"
	"io/ioutil"
	"log"
)

// Service holder for the encryption package
type Service struct {
	Hash      hash.Hash
	PublicKey *rsa.PublicKey
}

// New returns a new instance of encryption
func New(publicKey string) Service {
	if publicKey == "" {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

		if err != nil {
			log.Fatalf("error generating random key: %v", err)
		}

		return Service{
			Hash:      sha256.New(),
			PublicKey: &privateKey.PublicKey,
		}
	}

	b, err := ioutil.ReadFile(publicKey)

	if err != nil {
		log.Fatalf("Could not read public key from file: %v", err)
	}

	block, _ := pem.Decode([]byte(b))

	key, err := x509.ParsePKCS1PublicKey(block.Bytes)

	if err != nil {
		log.Fatalf("unable to decode public key: %v", err)
	}

	return Service{
		Hash:      sha256.New(),
		PublicKey: key,
	}
}

func (s *Service) Encrypt(msg []byte, label []byte) string {
	// OAEP is the recommended padding choice for encryption for new protocols or applications.
	// PKCS1v15 should only be used to support legacy protocols.
	ciphertext, err := rsa.EncryptOAEP(s.Hash, rand.Reader, s.PublicKey, msg, label)

	if err != nil {
		log.Fatalf("error encrypting message: %v", err)
	}

	return string(ciphertext)
}
