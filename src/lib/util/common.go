package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"net"
	"net/http"
	"strings"
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

func GetIP(req *http.Request) (ip net.IP) {
	var rmtAddr string
	fwdAddr := req.Header.Get("X-Forwarded-For") // capitalisation doesn't matter
	if fwdAddr != "" {
		// Got X-Forwarded-For
		rmtAddr = fwdAddr // If it's a single IP, then awesome!

		// If we got an array... grab the first IP
		ips := strings.Split(fwdAddr, ", ")
		if len(ips) > 1 {
			rmtAddr = ips[0]
		}
		ip = net.ParseIP(rmtAddr)
		if ip != nil {
			return ip
		}
	}

	// no fwd Addr or the first IP in the fwd addr is invalid
	rmtAddr = req.RemoteAddr
	// assume the http set the correct format of ip:port
	rmtAddr, _, _ = net.SplitHostPort(req.RemoteAddr)

	return net.ParseIP(rmtAddr)
}

func init() {
	gob.Register(net.IP{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[string]interface{}{})
}
