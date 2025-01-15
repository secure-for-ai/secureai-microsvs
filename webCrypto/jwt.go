package webCrypto

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GetJWTToken(data jwt.MapClaims, secret string, expire time.Duration) (token string) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	claims := t.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(expire).Unix()
	token, _ = t.SignedString([]byte(secret))
	return
}

func ValidateJWT(token, secret string) (jwt.MapClaims, bool) {
	t, _ := jwt.Parse(token, func(jwtToken *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, t.Valid
	}
	if exp, ok := claims["exp"].(float64); ok {
		claims["exp"] = int64(exp)
	}
	return claims, t.Valid
}
