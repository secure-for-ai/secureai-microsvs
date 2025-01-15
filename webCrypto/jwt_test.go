package webCrypto_test

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"github.com/secure-for-ai/secureai-microsvs/webCrypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetJWTToken(t *testing.T) {
	claims := jwt.MapClaims{
		"foo": "bar",
	}
	key, _ := util.GenerateRandomKey(32)
	secret := util.Base64EncodeToString(key)
	token := webCrypto.GetJWTToken(claims, secret, 60)

	fmt.Println(token)

	claimsNew, valid := webCrypto.ValidateJWT(token, secret)

	assert.Equal(t, claims, claimsNew)
	assert.Equal(t, valid, true)
}
