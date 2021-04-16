package webCrypto

import (
	"crypto"
	"crypto/hmac"
	"encoding/binary"
	"errors"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"strconv"
	"strings"
	"time"
)

var ErrCSRFTokenGetFail = errors.New("webCrypto/csrf token get failed")
var ErrCSRFTokenInvalid = errors.New("webCrypto/csrf token is invalid")

type CSRF interface {
	GetToken(session string, expire time.Duration) (string, error)
	ValidToken(session, token string) bool
	// ValidateTokenFromHTTP(r *http.Request, cookie, csrfHeader string) bool
}

// Encryption based CSRF using AES GCM
type aesGcmCSRF struct {
	cipher Cipher
}

func NewAesGcmCSRF(key interface{}) (CSRF, error) {
	cipher, err := NewAesGcm(key)
	if err != nil {
		return nil, err
	}
	csrf := &aesGcmCSRF{
		cipher: cipher,
	}
	return csrf, nil
}

func (c *aesGcmCSRF) GetToken(session string, expire time.Duration) (string, error) {
	sessLen := len(session)
	plaintext := make([]byte, sessLen+8)
	timestamp := time.Now().Add(expire).Unix()
	binary.LittleEndian.PutUint64(plaintext, uint64(timestamp))
	copy(plaintext[8:], session)
	token, err := c.cipher.EncryptBase64(plaintext)
	if err != nil {
		return "", ErrCSRFTokenInvalid
	}
	return token, nil
}

func (c *aesGcmCSRF) ValidToken(session, token string) bool {
	plaintext, err := c.cipher.DecryptBase64(token)
	if err != nil {
		return false
	}

	// check length
	sessLen := len(session)
	if len(plaintext) != sessLen+8 {
		return false
	}

	tkSess := string(plaintext[8:])
	timestamp := int64(binary.LittleEndian.Uint64(plaintext))
	currentTime := time.Now().Unix()
	// check whether the token is expired, or the extracted session
	// is not equal to the current session
	return currentTime < timestamp && tkSess == session
}

// HMAC based CSRF
type hmacCSRF struct {
	hash crypto.Hash
	key  []byte
}

var ErrHmacCSRFNotAvailable = errors.New("webCrypto/csrf this HMAC CSRF algorithm is not available")

type ErrHmacCSRFKeySize int

func (k ErrHmacCSRFKeySize) Error() string {
	return "webCrypto/csrf: invalid hmac key size " + strconv.Itoa(int(k))
}

func NewHmacCSRF(hash crypto.Hash, key interface{}) (CSRF, error) {
	if !hash.Available() {
		return nil, ErrHmacCSRFNotAvailable
	}

	var keyBin []byte
	switch v := key.(type) {
	case string:
		k, err := util.Base64Decode(v)
		if err != nil {
			return nil, ErrCipherKey
		}
		keyBin = k
	case []byte:
		keyBin = v
	default:
		return nil, ErrCipherKey
	}

	// check key length, at least hash output size
	keyLen := len(keyBin)
	if keyLen < hash.Size() {
		return nil, ErrHmacCSRFKeySize(keyLen)
	}

	csrf := &hmacCSRF{
		hash: hash,
		key:  keyBin,
	}
	return csrf, nil
}

func (c *hmacCSRF) GetToken(session string, expire time.Duration) (string, error) {
	// generate the nonce of token
	mac := hmac.New(c.hash.New, c.key)
	macSize := mac.Size()
	nonce, err := util.GenerateRandomKey(macSize)
	if err != nil {
		return "", err
	}

	// plaintext = nonce || timestamp
	plaintext := make([]byte, macSize+8)
	timestamp := time.Now().Add(expire).Unix()
	copy(plaintext, nonce)
	binary.LittleEndian.PutUint64(plaintext[macSize:], uint64(timestamp))

	// token = hmac(nonce || timestamp || session)
	mac.Write(plaintext)
	mac.Write([]byte(session))

	tokenMAC := mac.Sum(nil)
	token := util.Base64Encode(plaintext) + "." + util.Base64Encode(tokenMAC)
	return token, nil
}

func (c *hmacCSRF) ValidToken(session, token string) bool {
	tokenArr := strings.Split(token, ".")
	// check token format
	if len(tokenArr) != 2 {
		return false
	}
	// base64 validation
	plaintext, err := util.Base64Decode(tokenArr[0])
	if err != nil {
		return false
	}
	tokenMAC, err := util.Base64Decode(tokenArr[1])
	if err != nil {
		return false
	}

	mac := hmac.New(c.hash.New, c.key)
	macSize := mac.Size()
	// invalid plaintext length
	if len(plaintext) != macSize+8 {
		return false
	}
	timestamp := int64(binary.LittleEndian.Uint64(plaintext[macSize:]))
	currentTime := time.Now().Unix()
	// check whether token is expired
	if currentTime >= timestamp {
		return false
	}

	// mac check
	mac.Write(plaintext)
	mac.Write([]byte(session))
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(tokenMAC, expectedMAC)
}
