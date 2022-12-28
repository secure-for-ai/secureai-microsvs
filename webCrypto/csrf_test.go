package webCrypto_test

import (
	"github.com/minio/sha256-simd"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"github.com/secure-for-ai/secureai-microsvs/webCrypto"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestNewAesGcmCSRF(t *testing.T) {
	key, _ := util.GenerateRandomKey(32)
	csrf, err := webCrypto.NewAesGcmCSRF(key)
	sessionInfoBin, _ := util.GenerateRandomKey(100)
	sessionInfo := util.Base64EncodeToString(sessionInfoBin)
	expire := 60 * 15 * time.Second
	assert.NoError(t, err)

	t.Log("==============Test AES_GCM CSRF Valid Token===========")
	token, err := csrf.GetToken(sessionInfo, expire)
	assert.NoError(t, err)
	t.Log("session:", sessionInfo)
	t.Log("expired:", expire)
	t.Log("token:", token)
	valid := csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, true, valid)

	t.Log("==============Test AES_GCM CSRF Token Expiring on Time===========")
	token, err = csrf.GetToken(sessionInfo, 0)
	assert.NoError(t, err)
	t.Log("session:", sessionInfo)
	t.Log("expired:", 0)
	t.Log("token:", token)
	valid = csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, false, valid)

	t.Log("==============Test AES_GCM CSRF Expired Token===========")
	token, err = csrf.GetToken(sessionInfo, -expire)
	assert.NoError(t, err)
	t.Log("session:", sessionInfo)
	t.Log("expired:", -expire)
	t.Log("token:", token)
	valid = csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, false, valid)
}

func TestNewHmacCSRF(t *testing.T) {
	key, _ := util.GenerateRandomKey(32)
	csrf, err := webCrypto.NewHmacCSRF(sha256.New, key)
	sessionInfoBin, _ := util.GenerateRandomKey(100)
	sessionInfo := util.Base64EncodeToString(sessionInfoBin)
	expire := 60 * 15 * time.Second
	assert.NoError(t, err)

	t.Log("==============Test HMAC CSRF Valid Token===========")
	token, err := csrf.GetToken(sessionInfo, expire)
	assert.NoError(t, err)
	t.Log("session:", sessionInfo)
	t.Log("expired:", expire)
	t.Log("token:", token)
	valid := csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, true, valid)

	t.Log("==============Test HMAC CSRF Token Expiring on Time===========")
	token, err = csrf.GetToken(sessionInfo, 0)
	assert.NoError(t, err)
	t.Log("session:", sessionInfo)
	t.Log("expired:", 0)
	t.Log("token:", token)
	valid = csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, false, valid)

	t.Log("==============Test HMAC CSRF Expired Token===========")
	token, err = csrf.GetToken(sessionInfo, -expire)
	assert.NoError(t, err)
	t.Log("session:", sessionInfo)
	t.Log("expired:", -expire)
	t.Log("token:", token)
	valid = csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, false, valid)
}

func benchmarkNewAesGcmCSRF(b *testing.B, keyLen, msgLen int) {
	key, _ := util.GenerateRandomKey(keyLen)
	csrf, err := webCrypto.NewAesGcmCSRF(key)
	sessionInfoBin, _ := util.GenerateRandomKey(msgLen)
	sessionInfo := util.Base64EncodeToString(sessionInfoBin)
	assert.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(msgLen))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			expire := time.Duration(rand.Intn(100)+-50) * time.Second * 2
			token, _ := csrf.GetToken(sessionInfo, expire)
			valid := csrf.ValidToken(sessionInfo, token)
			assert.Equal(b, expire > 0, valid)
		}
	})
}

// AES 128 GCM N Bytes
func BenchmarkNewAes128GcmCSRF_16(b *testing.B)    { benchmarkNewAesGcmCSRF(b, 16, 16) }
func BenchmarkNewAes128GcmCSRF_128(b *testing.B)   { benchmarkNewAesGcmCSRF(b, 16, 128) }
func BenchmarkNewAes128GcmCSRF_1024(b *testing.B)  { benchmarkNewAesGcmCSRF(b, 16, 1024) }
func BenchmarkNewAes128GcmCSRF_8192(b *testing.B)  { benchmarkNewAesGcmCSRF(b, 16, 8192) }
func BenchmarkNewAes128GcmCSRF_65336(b *testing.B) { benchmarkNewAesGcmCSRF(b, 16, 65336) }

// AES 192 GCM N Bytes
func BenchmarkNewAes192GcmCSRF_16(b *testing.B)    { benchmarkNewAesGcmCSRF(b, 24, 16) }
func BenchmarkNewAes192GcmCSRF_128(b *testing.B)   { benchmarkNewAesGcmCSRF(b, 24, 128) }
func BenchmarkNewAes192GcmCSRF_1024(b *testing.B)  { benchmarkNewAesGcmCSRF(b, 24, 1024) }
func BenchmarkNewAes192GcmCSRF_8192(b *testing.B)  { benchmarkNewAesGcmCSRF(b, 24, 8192) }
func BenchmarkNewAes192GcmCSRF_65336(b *testing.B) { benchmarkNewAesGcmCSRF(b, 24, 65336) }

// AES 256 GCM N Bytes
func BenchmarkNewAes256GcmCSRF_16(b *testing.B)    { benchmarkNewAesGcmCSRF(b, 32, 16) }
func BenchmarkNewAes256GcmCSRF_128(b *testing.B)   { benchmarkNewAesGcmCSRF(b, 32, 128) }
func BenchmarkNewAes256GcmCSRF_1024(b *testing.B)  { benchmarkNewAesGcmCSRF(b, 32, 1024) }
func BenchmarkNewAes256GcmCSRF_8192(b *testing.B)  { benchmarkNewAesGcmCSRF(b, 32, 8192) }
func BenchmarkNewAes256GcmCSRF_65336(b *testing.B) { benchmarkNewAesGcmCSRF(b, 32, 65336) }

func benchmarkNewHmacCSRF(b *testing.B, keyLen, msgLen int) {
	key, _ := util.GenerateRandomKey(keyLen)
	csrf, err := webCrypto.NewHmacCSRF(sha256.New, key)
	sessionInfoBin, _ := util.GenerateRandomKey(msgLen)      // "sample_session"
	sessionInfo := util.Base64EncodeToString(sessionInfoBin) // "sample_session"
	assert.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(msgLen))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			expire := time.Duration(rand.Intn(100)+-50) * time.Second * 2
			token, _ := csrf.GetToken(sessionInfo, expire)
			valid := csrf.ValidToken(sessionInfo, token)
			assert.Equal(b, expire > 0, valid)
		}
	})
}

// HMAC SHA256 N Bytes
func BenchmarkNewHmacCSRF_sha256_16(b *testing.B)    { benchmarkNewHmacCSRF(b, 32, 16) }
func BenchmarkNewHmacCSRF_sha256_128(b *testing.B)   { benchmarkNewHmacCSRF(b, 32, 128) }
func BenchmarkNewHmacCSRF_sha256_1024(b *testing.B)  { benchmarkNewHmacCSRF(b, 32, 1024) }
func BenchmarkNewHmacCSRF_sha256_8192(b *testing.B)  { benchmarkNewHmacCSRF(b, 32, 8192) }
func BenchmarkNewHmacCSRF_sha256_65636(b *testing.B) { benchmarkNewHmacCSRF(b, 32, 65336) }

// HMAC SHA384 N Bytes
func BenchmarkNewHmacCSRF_sha384_16(b *testing.B)    { benchmarkNewHmacCSRF(b, 48, 16) }
func BenchmarkNewHmacCSRF_sha384_128(b *testing.B)   { benchmarkNewHmacCSRF(b, 48, 128) }
func BenchmarkNewHmacCSRF_sha384_1024(b *testing.B)  { benchmarkNewHmacCSRF(b, 48, 1024) }
func BenchmarkNewHmacCSRF_sha384_8192(b *testing.B)  { benchmarkNewHmacCSRF(b, 48, 8192) }
func BenchmarkNewHmacCSRF_sha384_65636(b *testing.B) { benchmarkNewHmacCSRF(b, 48, 65336) }

// HMAC SHA512 N Bytes
func BenchmarkNewHmacCSRF_sha512_16(b *testing.B)    { benchmarkNewHmacCSRF(b, 64, 16) }
func BenchmarkNewHmacCSRF_sha512_128(b *testing.B)   { benchmarkNewHmacCSRF(b, 64, 128) }
func BenchmarkNewHmacCSRF_sha512_1024(b *testing.B)  { benchmarkNewHmacCSRF(b, 64, 1024) }
func BenchmarkNewHmacCSRF_sha512_8192(b *testing.B)  { benchmarkNewHmacCSRF(b, 64, 8192) }
func BenchmarkNewHmacCSRF_sha512_65636(b *testing.B) { benchmarkNewHmacCSRF(b, 64, 65336) }
