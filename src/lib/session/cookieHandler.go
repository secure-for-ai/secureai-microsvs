package session

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"github.com/gorilla/sessions"
	"net/http"
	"template2/lib/util"
)

type AesGcmCookieHandler struct {
	//signer hash.Hash
	//cipher cipher.Block
	ahead cipher.AEAD
}

func NewAesGcmCookieHandler(encKey []byte /*, authKey []byte*/) (*AesGcmCookieHandler, error) {
	c, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}
	ahead, err := cipher.NewGCMWithNonceSize(c, 12)
	if err != nil {
		return nil, err
	}
	return &AesGcmCookieHandler{
		//signer: hmac.New(sha256.New, authKey),
		//cipher: c,
		ahead: ahead,
	}, nil
}
func (h AesGcmCookieHandler) Encode(sess *sessions.Session, store *HybridStore) (string, error) {

	sessIDLen := len(sess.ID)
	inputLen := sessIDLen + 8
	input := make([]byte, inputLen)
	// append current timestamp
	timestamp := util.GetNowTimestamp()
	copy(input, sess.ID)
	binary.LittleEndian.PutUint64(input[sessIDLen:], uint64(timestamp))

	iv, err := util.GenerateRandomKey(12)

	if err != nil {
		// Todo add logger
		return "", err
	}

	cipherText := h.ahead.Seal(nil, iv, input, nil)
	output := make([]byte, 12+len(cipherText))

	copy(output, iv)
	copy(output[12:], cipherText)

	return util.Base64Encode(output), nil
}

func (h AesGcmCookieHandler) Decode(c *http.Cookie, sess *sessions.Session, store *HybridStore) error {
	//return securecookie.DecodeMulti(c.Name, c.Value, &sess.ID, h.Codecs...)

	ciperText, err := util.Base64Decode(c.Value)

	if err != nil {
		// Todo add logger
		return err
	}

	// check cipherText format
	// cipherText must contain 12 byte initial vector + 16 byte tag
	// Todo dynamic configurable initial vector size + tag size
	if len(ciperText) < 28 {
		return ErrInvalidCookie
	}

	plainText, err := h.ahead.Open(nil, ciperText[0:12], ciperText[12:], nil)

	if err != nil {
		// Todo add logger
		return err
	}

	plainTextLen := len(plainText)

	// check plainText format
	// plaintext contains 8 bytes timestamp plus random session ID
	if plainTextLen < 9 {
		return ErrInvalidCookie
	}

	timestamp := int64(binary.LittleEndian.Uint64(plainText[plainTextLen-8:]))

	currentTime := util.GetNowTimestamp()
	expireTime := timestamp + int64(sess.Options.MaxAge)

	if timestamp > currentTime || expireTime < currentTime {
		// too new or too old
		return ErrInvalidCookie
	}

	sess.ID = string(plainText[:plainTextLen-8])

	return nil
}
