package webCrypto_test

import (
	"bytes"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"github.com/secure-for-ai/secureai-microsvs/webCrypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCipherBase64(t *testing.T) {
	key, _ := util.GenerateRandomKey(32)
	cipher, err := webCrypto.NewAesGcm(util.Base64EncodeToString(key))
	assert.NoError(t, err)

	t.Log("==============Test Encryption Base64===========")
	plainText := []byte("hello world")
	t.Log("PlainText:", string(plainText))
	cipherText, err := cipher.EncryptBase64ToString(plainText)
	assert.NoError(t, err)
	t.Log("CipherText:", cipherText)
	decryptText, err := cipher.DecryptBase64ToString(cipherText)
	assert.NoError(t, err)
	t.Log("DecryptText:", string(decryptText))

	assert.Equal(t, plainText, decryptText)
}

func TestCipherByte(t *testing.T) {
	key, _ := util.GenerateRandomKey(32)
	cipher, err := webCrypto.NewAesGcm(util.Base64EncodeToString(key))
	assert.NoError(t, err)

	t.Log("===============Test Encryption Byte============")
	plainText := []byte("hello world")
	t.Log("PlainText:", string(plainText))
	cipherText, err := cipher.EncryptByte(plainText)
	assert.NoError(t, err)
	t.Log("CipherText:", cipherText)
	decryptText, err := cipher.DecryptByte(cipherText)
	assert.NoError(t, err)
	t.Log("DecryptText:", string(decryptText))

	assert.Equal(t, plainText, decryptText)
}

func benchmarkCipherByte(b *testing.B, keyLen, n int) {
	key, _ := util.GenerateRandomKey(keyLen)
	cipher, err := webCrypto.NewAesGcm(util.Base64EncodeToString(key))
	assert.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(n))

	b.RunParallel(func(pb *testing.PB) {
		plainText, _ := util.GenerateRandomKey(n)
		for pb.Next() {
			cipherText, _ := cipher.EncryptByte(plainText)
			decrypted, _ := cipher.DecryptByte(cipherText)
			if !bytes.Equal(plainText, decrypted) {
				assert.Equal(b, plainText, decrypted)
			}
		}
	})
}

// AES 128 N Bytes
func BenchmarkCipherByte128_16(b *testing.B)    { benchmarkCipherByte(b, 16, 16) }
func BenchmarkCipherByte128_128(b *testing.B)   { benchmarkCipherByte(b, 16, 128) }
func BenchmarkCipherByte128_1024(b *testing.B)  { benchmarkCipherByte(b, 16, 1024) }
func BenchmarkCipherByte128_8192(b *testing.B)  { benchmarkCipherByte(b, 16, 8192) }
func BenchmarkCipherByte128_65536(b *testing.B) { benchmarkCipherByte(b, 16, 65536) }

// AES 192 N Bytes
func BenchmarkCipherByte192_16(b *testing.B)    { benchmarkCipherByte(b, 24, 16) }
func BenchmarkCipherByte192_128(b *testing.B)   { benchmarkCipherByte(b, 24, 128) }
func BenchmarkCipherByte192_1024(b *testing.B)  { benchmarkCipherByte(b, 24, 1024) }
func BenchmarkCipherByte192_8192(b *testing.B)  { benchmarkCipherByte(b, 24, 8192) }
func BenchmarkCipherByte192_65536(b *testing.B) { benchmarkCipherByte(b, 24, 65536) }

// AES 256 N Bytes
func BenchmarkCipherByte256_16(b *testing.B)    { benchmarkCipherByte(b, 32, 16) }
func BenchmarkCipherByte256_128(b *testing.B)   { benchmarkCipherByte(b, 32, 128) }
func BenchmarkCipherByte256_1024(b *testing.B)  { benchmarkCipherByte(b, 32, 1024) }
func BenchmarkCipherByte256_8192(b *testing.B)  { benchmarkCipherByte(b, 32, 8192) }
func BenchmarkCipherByte256_65536(b *testing.B) { benchmarkCipherByte(b, 32, 65536) }

func benchmarkCipherBase64(b *testing.B, keyLen, n int) {
	key, _ := util.GenerateRandomKey(keyLen)
	cipher, err := webCrypto.NewAesGcm(util.Base64EncodeToString(key))
	assert.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(n))

	b.RunParallel(func(pb *testing.PB) {
		plainText, _ := util.GenerateRandomKey(n)
		for pb.Next() {
			cipherText, _ := cipher.EncryptBase64ToString(plainText)
			decrypted, _ := cipher.DecryptBase64ToString(cipherText)
			if !bytes.Equal(plainText, decrypted) {
				assert.Equal(b, plainText, decrypted)
			}
		}
	})
}

// AES 128 N Bytes
func BenchmarkCipherBase64_128_16(b *testing.B)    { benchmarkCipherBase64(b, 16, 16) }
func BenchmarkCipherBase64_128_128(b *testing.B)   { benchmarkCipherBase64(b, 16, 128) }
func BenchmarkCipherBase64_128_1024(b *testing.B)  { benchmarkCipherBase64(b, 16, 1024) }
func BenchmarkCipherBase64_128_8192(b *testing.B)  { benchmarkCipherBase64(b, 16, 8192) }
func BenchmarkCipherBase64_128_65536(b *testing.B) { benchmarkCipherBase64(b, 16, 65536) }

// AES 192 N Bytes
func BenchmarkCipherBase64_192_16(b *testing.B)    { benchmarkCipherBase64(b, 24, 16) }
func BenchmarkCipherBase64_192_128(b *testing.B)   { benchmarkCipherBase64(b, 24, 128) }
func BenchmarkCipherBase64_192_1024(b *testing.B)  { benchmarkCipherBase64(b, 24, 1024) }
func BenchmarkCipherBase64_192_8192(b *testing.B)  { benchmarkCipherBase64(b, 24, 8192) }
func BenchmarkCipherBase64_192_65536(b *testing.B) { benchmarkCipherBase64(b, 24, 65536) }

// AES 256 N Bytes
func BenchmarkCipherBase64_256_16(b *testing.B)    { benchmarkCipherBase64(b, 32, 16) }
func BenchmarkCipherBase64_256_128(b *testing.B)   { benchmarkCipherBase64(b, 32, 128) }
func BenchmarkCipherBase64_256_1024(b *testing.B)  { benchmarkCipherBase64(b, 32, 1024) }
func BenchmarkCipherBase64_256_8192(b *testing.B)  { benchmarkCipherBase64(b, 32, 8192) }
func BenchmarkCipherBase64_256_65536(b *testing.B) { benchmarkCipherBase64(b, 32, 65536) }
