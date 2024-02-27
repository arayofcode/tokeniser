package cipher

type cipherConfig struct {
	passphrase string
}

type Cipher interface {
	Encrypt(plaintextByte []byte) (ciphertext []byte, salt []byte, err error)
	Decrypt(cipherTextWithNonceAndSalt []byte) (plaintext []byte, err error)
}

func Init(passphrase string) Cipher {
	return &cipherConfig{
		passphrase: passphrase,
	}
}
