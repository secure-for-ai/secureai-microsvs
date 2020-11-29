package webCrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"template2/lib/util"
)

var ErrCipherText = errors.New("cipherText is invalid")
var ErrCipherKey = errors.New("cipher Key should be either encoded in base64 or a binary slice")

type Cipher interface {
	EncryptBase64(plaintext []byte) (string, error)
	DecryptBase64(cipherText string) ([]byte, error)
	EncryptByte(plaintext []byte) ([]byte, error)
	DecryptByte(cipherText []byte) ([]byte, error)
}

type aesGcmAEAD struct {
	cipher.AEAD
}

func NewAesGcm(encKey interface{}) (Cipher, error) {
	var encKeyBin []byte
	switch v := encKey.(type) {
	case string:
		k, err := util.Base64Decode(v)
		if err != nil {
			return nil, ErrCipherKey
		}
		encKeyBin = k
	case []byte:
		encKeyBin = v
	default:
		return nil, ErrCipherKey
	}

	c, err := aes.NewCipher(encKeyBin)
	if err != nil {
		return nil, err
	}
	ahead, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	return &aesGcmAEAD{
		ahead,
	}, nil
}

func (c *aesGcmAEAD) EncryptByte(plaintext []byte) ([]byte, error) {
	iv, err := util.GenerateRandomKey(12)

	if err != nil {
		// Todo add logger
		return []byte{}, err
	}

	cipherText := c.Seal(nil, iv, plaintext, nil)
	output := make([]byte, 12+len(cipherText))

	copy(output, iv)
	copy(output[12:], cipherText)

	return output, nil
}

func (c *aesGcmAEAD) DecryptByte(ciphertext []byte) ([]byte, error) {
	// check ciphertext format
	// ciphertext must contain 12 byte initial vector + 16 byte tag
	// Todo dynamic configurable initial vector size + tag size
	if len(ciphertext) < 28 {
		return nil, ErrCipherText
	}

	plainText, err := c.Open(nil, ciphertext[0:12], ciphertext[12:], nil)

	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func (c *aesGcmAEAD) EncryptBase64(plaintext []byte) (string, error) {

	output, err := c.EncryptByte(plaintext)

	if err != nil {
		return "", err
	}

	return util.Base64Encode(output), nil
}

func (c *aesGcmAEAD) DecryptBase64(ciphertext string) ([]byte, error) {
	cipherText, err := util.Base64Decode(ciphertext)

	if err != nil {
		return nil, err
	}

	return c.DecryptByte(cipherText)
}
