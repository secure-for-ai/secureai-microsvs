package util

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

func GetNowTimestamp() int64 {
	return time.Now().Unix()
}

func GenerateRandomKey(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// use ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_
// no padding
func Base64Encode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

// use ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_
// no padding
func Base64Decode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}
