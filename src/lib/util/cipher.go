package util

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var ErrCipherText = errors.New("ErrCipherText")

type AesGcmCipher struct {
	ahead cipher.AEAD
}

func NewAesGcmCipher(encKey string) (*AesGcmCipher, error) {
	encKeyBin, err := Base64Decode(encKey)
	if err != nil {
		return nil, err
	}
	c, err := aes.NewCipher(encKeyBin)
	if err != nil {
		return nil, err
	}
	ahead, err := cipher.NewGCMWithNonceSize(c, 12)
	if err != nil {
		return nil, err
	}
	return &AesGcmCipher{
		ahead: ahead,
	}, nil
}

func (c AesGcmCipher) Encrypt(str string) (string, error) {

	strByte := []byte(str)

	iv, err := GenerateRandomKey(12)

	if err != nil {
		// Todo add logger
		return "", err
	}

	cipherText := c.ahead.Seal(nil, iv, strByte, nil)
	output := make([]byte, 12+len(cipherText))

	copy(output, iv)
	copy(output[12:], cipherText)

	return Base64Encode(output), nil
}

func (c AesGcmCipher) Decrypt(str string) (string, error) {
	ciperText, err := Base64Decode(str)

	if err != nil {
		return "", err
	}

	// check cipherText format
	// cipherText must contain 12 byte initial vector + 16 byte tag
	// Todo dynamic configurable initial vector size + tag size
	if len(ciperText) < 28 {
		return "", ErrCipherText
	}

	plainText, err := c.ahead.Open(nil, ciperText[0:12], ciperText[12:], nil)

	if err != nil {
		return "", err
	}

	plainTextLen := len(plainText)

	// check plainText format
	// plaintext contains 8 bytes timestamp plus random session ID
	if plainTextLen < 9 {
		return "", ErrCipherText
	}

	output := string(plainText)

	return output, nil
}
