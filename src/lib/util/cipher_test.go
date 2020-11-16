package util_test

import (
	"github.com/stretchr/testify/assert"
	"template2/lib/util"
	"testing"
)

func TestCipher(t *testing.T) {
	key := "pR6kDdHYqNMRO74kUxFiGgpv3A6qKFeCY6IDHxDH8NY"
	cipher, err := util.NewAesGcmCipher(key)
	assert.NoError(t, err)

	t.Log("==============Test Encryption===========")
	plainText := "hello world"
	t.Log("PlainText:", plainText)
	cipherText, err := cipher.Encrypt(plainText)
	assert.NoError(t, err)
	t.Log("CipherText:", cipherText)
	decryptText, err := cipher.Decrypt(cipherText)
	t.Log("DecryptText:", decryptText)

	assert.Equal(t, plainText, decryptText)
}
