package cipher

import (
	"fmt"
	"io"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"

	"golang.org/x/crypto/pbkdf2"
)

const (
	aesKeySize       = 32
	bcryptCostFactor = 12
	saltSize         = 16     // NIST recommendation of 128 bits (16 bytes)
	pbkdf2Iterations = 600000 // OWASP Recommended
)

func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func (c *cipherConfig) deriveKeyFromPassword(salt []byte) (key []byte) {
	return pbkdf2.Key([]byte(c.passphrase), salt, pbkdf2Iterations, aesKeySize, sha512.New)
}

func (c *cipherConfig) Encrypt(dataByte []byte) (ciphertext []byte, err error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}

	key := c.deriveKeyFromPassword(salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext = gcm.Seal(nonce, nonce, dataByte, nil)
	ciphertext = append(salt, ciphertext...)
	return ciphertext, nil
}

func (c *cipherConfig) Decrypt(cipherTextWithNonceAndSalt []byte) (plaintext []byte, err error) {
	if len(cipherTextWithNonceAndSalt) < saltSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	salt := cipherTextWithNonceAndSalt[:saltSize]
	cipherTextWithNonce := cipherTextWithNonceAndSalt[saltSize:]

	key := c.deriveKeyFromPassword(salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherTextWithNonce) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := cipherTextWithNonce[:nonceSize], cipherTextWithNonce[nonceSize:]

	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
