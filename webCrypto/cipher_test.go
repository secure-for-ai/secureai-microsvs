package webCrypto_test

import (
	"github.com/secure-for-ai/secureai-microsvs/util"
	"github.com/secure-for-ai/secureai-microsvs/webCrypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCipherBase64(t *testing.T) {
	key, _ := util.GenerateRandomKey(32)
	cipher, err := webCrypto.NewAesGcm(util.Base64Encode(key))
	assert.NoError(t, err)

	t.Log("==============Test Encryption Base64===========")
	plainText := []byte("hello world")
	t.Log("PlainText:", string(plainText))
	cipherText, err := cipher.EncryptBase64(plainText)
	assert.NoError(t, err)
	t.Log("CipherText:", cipherText)
	decryptText, err := cipher.DecryptBase64(cipherText)
	t.Log("DecryptText:", string(decryptText))

	assert.Equal(t, plainText, decryptText)
}

func TestCipherByte(t *testing.T) {
	key, _ := util.GenerateRandomKey(32)
	cipher, err := webCrypto.NewAesGcm(util.Base64Encode(key))
	assert.NoError(t, err)

	t.Log("===============Test Encryption Byte============")
	plainText := []byte("hello world")
	t.Log("PlainText:", string(plainText))
	cipherText, err := cipher.EncryptByte(plainText)
	assert.NoError(t, err)
	t.Log("CipherText:", cipherText)
	decryptText, err := cipher.DecryptByte(cipherText)
	t.Log("DecryptText:", string(decryptText))

	assert.Equal(t, plainText, decryptText)
}

func BenchmarkCipherByte(b *testing.B) {
	key, _ := util.GenerateRandomKey(32)
	cipher, err := webCrypto.NewAesGcm(util.Base64Encode(key))
	assert.NoError(b, err)

	plainText, _ := util.GenerateRandomKey(1024)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cipherText, _ := cipher.EncryptByte(plainText)
			decrypted, _ := cipher.DecryptByte(cipherText)
			assert.Equal(b, plainText, decrypted)
		}
	})
}

func BenchmarkCipherBase64(b *testing.B) {
	key, _ := util.GenerateRandomKey(32)
	cipher, err := webCrypto.NewAesGcm(util.Base64Encode(key))
	assert.NoError(b, err)

	plainText, _ := util.GenerateRandomKey(1024)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cipherText, _ := cipher.EncryptBase64(plainText)
			decrypted, _ := cipher.DecryptBase64(cipherText)
			assert.Equal(b, plainText, decrypted)
		}
	})
}
