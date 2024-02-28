package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/pbkdf2"
)

const (
	aesKeySize       = 32
	saltSize         = 16     // NIST recommendation of 128 bits (16 bytes)
	pbkdf2Iterations = 600000 // OWASP Recommended
)

func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		log.Error().Err(err).Msg("Failed to generate salt")
		return nil, err
	}
	return salt, nil
}

func (c *cipherConfig) deriveKeyFromPassword(salt []byte) []byte {
	return pbkdf2.Key([]byte(c.passphrase), salt, pbkdf2Iterations, aesKeySize, sha512.New)
}

func (c *cipherConfig) Encrypt(dataByte []byte) ([]byte, error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}

	key := c.deriveKeyFromPassword(salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create AES cipher")
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GCM cipher mode")
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Error().Err(err).Msg("Failed to generate nonce")
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, dataByte, nil)
	ciphertext = append(salt, ciphertext...)
	log.Info().Msg("Data encrypted successfully")
	return ciphertext, nil
}

func (c *cipherConfig) Decrypt(cipherTextWithNonceAndSalt []byte) ([]byte, error) {
	if len(cipherTextWithNonceAndSalt) < saltSize {
		err := fmt.Errorf("ciphertext too short")
		log.Error().Err(err).Msg("Failed to decrypt data")
		return nil, err
	}

	salt := cipherTextWithNonceAndSalt[:saltSize]
	cipherTextWithNonce := cipherTextWithNonceAndSalt[saltSize:]
	key := c.deriveKeyFromPassword(salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create AES cipher for decryption")
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GCM cipher mode for decryption")
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherTextWithNonce) < nonceSize {
		err := fmt.Errorf("ciphertext too short for nonce")
		log.Error().Err(err).Msg("Failed to decrypt data")
		return nil, err
	}

	nonce, ciphertext := cipherTextWithNonce[:nonceSize], cipherTextWithNonce[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decrypt data with GCM")
		return nil, err
	}

	log.Info().Msg("Data decrypted successfully")
	return plaintext, nil
}
