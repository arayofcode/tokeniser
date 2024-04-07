package cipher

import (
	"testing"

	"github.com/arayofcode/tokeniser/src/common"
	"github.com/stretchr/testify/assert"
)

func TestEncryptionDecryption(t *testing.T) {
	passphrase := "very-secure-passphrase"
	plaintext := []byte("Is this working?")

	c := Init(passphrase)

	ciphertext, err := c.Encrypt(plaintext)
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

	ciphertext, err := c1.Encrypt(plaintext)
	assert.NoError(t, err)

	decrypted, err := c2.Decrypt(ciphertext)
	assert.NoError(t, err)

	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptionWithDifferentPassphrases(t *testing.T) {
	plaintext := []byte("Might as well think of some new plaintext.")

	c1 := Init("passphrase1")
	c2 := Init("passphrase2")

	ciphertext, err := c1.Encrypt(plaintext)
	assert.NoError(t, err)

	decrypted, err := c2.Decrypt(ciphertext)
	// Decrypting using wrong key should definitely create an error
	assert.Error(t, err)

	assert.NotEqual(t, plaintext, decrypted)
}

func TestEncryptionDecryptionEmptyPlaintext(t *testing.T) {
	passphrase := "so-many-insecure-passphrases-but-no-twist"
	plaintext := []byte("")

	c := Init(passphrase)

	ciphertext, err := c.Encrypt(plaintext)
	assert.NoError(t, err, "Encryption should not error on empty plaintext")

	decrypted, err := c.Decrypt(ciphertext)
	assert.NoError(t, err, "Decryption should not error on ciphertext from empty plaintext")

	common.AssertByteSliceEqual(t, plaintext, decrypted)
}
