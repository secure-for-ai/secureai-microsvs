package webCrypto_test

import (
	"crypto"
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
	sessionInfo := util.Base64Encode(sessionInfoBin)
	expire := 60 * 15 * time.Second
	assert.NoError(t, err)

	t.Log("==============Test AES_GCM CSRF Valid Token===========")
	token, err := csrf.GetToken(sessionInfo, expire)
	t.Log("session:", sessionInfo)
	t.Log("expired:", expire)
	t.Log("token:", token)
	valid := csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, true, valid)

	t.Log("==============Test AES_GCM CSRF Token Expiring on Time===========")
	token, err = csrf.GetToken(sessionInfo, 0)
	t.Log("session:", sessionInfo)
	t.Log("expired:", 0)
	t.Log("token:", token)
	valid = csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, false, valid)

	t.Log("==============Test AES_GCM CSRF Expired Token===========")
	token, err = csrf.GetToken(sessionInfo, -expire)
	t.Log("session:", sessionInfo)
	t.Log("expired:", -expire)
	t.Log("token:", token)
	valid = csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, false, valid)
}

func TestNewHmacCSRF(t *testing.T) {
	key, _ := util.GenerateRandomKey(32)
	csrf, err := webCrypto.NewHmacCSRF(crypto.SHA256, key)
	sessionInfoBin, _ := util.GenerateRandomKey(100)
	sessionInfo := util.Base64Encode(sessionInfoBin)
	expire := 60 * 15 * time.Second
	assert.NoError(t, err)

	t.Log("==============Test HMAC CSRF Valid Token===========")
	token, err := csrf.GetToken(sessionInfo, expire)
	t.Log("session:", sessionInfo)
	t.Log("expired:", expire)
	t.Log("token:", token)
	valid := csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, true, valid)

	t.Log("==============Test HMAC CSRF Token Expiring on Time===========")
	token, err = csrf.GetToken(sessionInfo, 0)
	t.Log("session:", sessionInfo)
	t.Log("expired:", 0)
	t.Log("token:", token)
	valid = csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, false, valid)

	t.Log("==============Test HMAC CSRF Expired Token===========")
	token, err = csrf.GetToken(sessionInfo, -expire)
	t.Log("session:", sessionInfo)
	t.Log("expired:", -expire)
	t.Log("token:", token)
	valid = csrf.ValidToken(sessionInfo, token)
	assert.Equal(t, false, valid)
}

func BenchmarkNewAesGcmCSRF(b *testing.B) {
	key, _ := util.GenerateRandomKey(32)
	csrf, err := webCrypto.NewAesGcmCSRF(key)
	sessionInfoBin, _ := util.GenerateRandomKey(100)
	sessionInfo := util.Base64Encode(sessionInfoBin)
	assert.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			expire := time.Duration(rand.Intn(100)+-50) * time.Second
			token, _ := csrf.GetToken(sessionInfo, expire)
			valid := csrf.ValidToken(sessionInfo, token)
			assert.Equal(b, expire > 0, valid)
		}
	})
}

func BenchmarkNewHmacCSRF(b *testing.B) {
	key, _ := util.GenerateRandomKey(32)
	csrf, err := webCrypto.NewHmacCSRF(crypto.SHA256, key)
	sessionInfoBin, _ := util.GenerateRandomKey(100) // "sample_session"
	sessionInfo := util.Base64Encode(sessionInfoBin) // "sample_session"
	assert.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			expire := time.Duration(rand.Intn(100)+-50) * time.Second
			token, _ := csrf.GetToken(sessionInfo, expire)
			valid := csrf.ValidToken(sessionInfo, token)
			assert.Equal(b, expire > 0, valid)
		}
	})
}
