package webCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"sync"
)

var ErrCipherText = errors.New("cipherText is invalid")
var ErrCipherKey = errors.New("cipher Key should be either encoded in base64 or a binary slice")

type Cipher interface {
	EncryptBase64ToString(plaintext []byte) (string, error)
	DecryptBase64ToString(cipherText string) ([]byte, error)
	EncryptByte(plaintext []byte) ([]byte, error)
	DecryptByte(cipherText []byte) ([]byte, error)
}

type aesGcmAEAD struct {
	cipher.AEAD
	minEncryptLen int
}

func NewAesGcm(encKey interface{}) (Cipher, error) {
	var encKeyBin []byte
	switch v := encKey.(type) {
	case string:
		k, err := util.Base64DecodeString(v)
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
		ahead.NonceSize() + ahead.Overhead(),
	}, nil
}

func (c *aesGcmAEAD) MinEncryptLen() int {
	return c.minEncryptLen
}

func (c *aesGcmAEAD) EncryptLen(plaintext []byte) int {
	return len(plaintext) + c.MinEncryptLen()
}

func (c *aesGcmAEAD) Encrypt(dst, plaintext []byte) ([]byte, error) {
	nonceSize := c.NonceSize()
	err := util.GenerateRandomKeyToBuf(dst, nonceSize)

	if err != nil {
		// Todo add logger
		return nil, err
	}

	iv := dst[0:nonceSize]
	ciphertext := c.Seal(dst[0:nonceSize], iv, plaintext, nil)

	return ciphertext, nil
}

func (c *aesGcmAEAD) EncryptByte(plaintext []byte) (cipherText []byte, err error) {
	cipherText = make([]byte, c.EncryptLen(plaintext))
	_, err = c.Encrypt(cipherText, plaintext)
	return
}

func (c *aesGcmAEAD) Decrypt(dst, ciphertext []byte) ([]byte, error) {
	// check ciphertext format
	// ciphertext must contain 12 byte initial vector + 16 byte tag, which is the MinEncryptLen
	// Todo dynamic configurable initial vector size + tag size
	if len(ciphertext) < c.MinEncryptLen() {
		return nil, ErrCipherText
	}

	nonceSize := c.NonceSize()
	plaintext, err := c.Open(dst[:0], ciphertext[0:nonceSize], ciphertext[nonceSize:], nil)

	if err != nil {
		return plaintext, err
	}

	return plaintext, err
}

func (c *aesGcmAEAD) DecryptByte(ciphertext []byte) ([]byte, error) {
	return c.Decrypt(nil, ciphertext)
}

var encryptBufPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func (c *aesGcmAEAD) EncryptBase64(dst, plaintext []byte) error {
	buf := encryptBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		encryptBufPool.Put(buf)
	}()

	bufCap := buf.Cap()
	bufCapMin := c.EncryptLen(plaintext)

	// detect the minimum cap need to increased, grow the capacity by bufCapMin - buf.Len()
	if bufCapMin > bufCap {
		buf.Grow(bufCapMin - buf.Len())
	}

	ciphertext, err := c.Encrypt(buf.Bytes(), plaintext)

	if err != nil {
		return err
	}

	util.Base64Encode(dst, ciphertext)

	return nil
}

func (c *aesGcmAEAD) EncryptBase64ToString(plaintext []byte) (string, error) {
	buf := make([]byte, base64.RawURLEncoding.EncodedLen(c.EncryptLen(plaintext)))
	err := c.EncryptBase64(buf, plaintext)
	if err != nil {
		return "", err
	}

	return util.FastBytesToString(buf), nil
}

func (c *aesGcmAEAD) DecryptBase64(dst []byte, ciphertext string) error {
	buf := encryptBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		encryptBufPool.Put(buf)
	}()

	bufCap := buf.Cap()
	bufCapMin := base64.RawURLEncoding.DecodedLen(len(ciphertext))

	// detect the minimum cap need to increased, grow the capacity by bufCapMin - buf.Len()
	if bufCapMin > bufCap {
		buf.Grow(bufCapMin - buf.Len())
	}

	cipherRaw := buf.Bytes()[:bufCapMin] //buf.Next(bufCapMin)
	_, err := util.Base64Decode(cipherRaw, util.FastStringToBytes(ciphertext))

	if err != nil {
		return err
	}

	_, err = c.Decrypt(dst, cipherRaw)

	return err
}

func (c *aesGcmAEAD) DecryptBase64ToString(ciphertext string) ([]byte, error) {
	bufLen := base64.RawURLEncoding.DecodedLen(len(ciphertext)) - c.MinEncryptLen()
	if bufLen < 0 {
		return nil, ErrCipherText
	}

	buf := make([]byte, bufLen)
	err := c.DecryptBase64(buf, ciphertext)

	if err != nil {
		return nil, err
	}

	return buf, err
}
