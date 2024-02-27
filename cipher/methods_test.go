package cipher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionDecryption(t *testing.T) {
	passphrase := "very-secure-passphrase"
	plaintext := []byte("Is this working?")

	c := Init(passphrase)

	ciphertext, _, err := c.Encrypt(plaintext)
	assert.NoError(t, err)

	decrypted, err := c.Decrypt(ciphertext)
	assert.NoError(t, err)

	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptDecryptWithDifferentInstances(t *testing.T) {
	passphrase := "another-very-secure-passphrase"
	plaintext := []byte("Boring! Is this even working?")

	c1 := Init(passphrase)
	c2 := Init(passphrase)

	ciphertext, _, err := c1.Encrypt(plaintext)
	assert.NoError(t, err)

	decrypted, err := c2.Decrypt(ciphertext)
	assert.NoError(t, err)

	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptionWithDifferentPassphrases(t *testing.T) {
	plaintext := []byte("Might as well think of some new plaintext.")

	c1 := Init("passphrase1")
	c2 := Init("passphrase2")

	ciphertext, _, err := c1.Encrypt(plaintext)
	assert.NoError(t, err)

	decrypted, err := c2.Decrypt(ciphertext)
	// Decrypting using wrong key should definitely create an error
	assert.Error(t, err)

	assert.NotEqual(t, plaintext, decrypted)
}
