package util

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func GetNowTimestamp() int64 {
	return time.Now().Unix()
}

func GenerateRandomKey(len int) (dst []byte, err error) {
	dst = make([]byte, len)
	err = GenerateRandomKeyToBuf(dst, len)
	return
}

func GenerateRandomKeyToBuf(dst []byte, len int) error {
	_, err := rand.Read(dst[0:len])
	if err != nil {
		return err
	}
	return nil
}

// use ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_
// no padding
func Base64EncodeToString(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

// use ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_
// no padding
func Base64DecodeString(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// use ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_
// no padding
func Base64Encode(dst []byte, src []byte) {
	base64.RawURLEncoding.Encode(dst, src)
}

// use ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_
// no padding
func Base64Decode(dst, src []byte) (int, error) {
	return base64.RawURLEncoding.Decode(dst, src)
}

func FastStringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func FastBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetIP(req *http.Request) (ip net.IP) {
	var rmtAddr string
	/* Todo: handle real client IP if the service is behind load balancer such as nginx */
	//log.Println("X-Real-IP", req.Header.Get("X-Real-IP"))
	//log.Println("X-Forwarded-For", req.Header.Get("X-Forwarded-For"))
	//log.Println("Host", req.Header.Get("Host"))
	//log.Println("RemoteAddr", req.RemoteAddr)
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

// get the concrete value of
func ReflectValue(value interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(value))
}

func init() {
	gob.Register(net.IP{})
	gob.Register(map[string]interface{}{})
}

func HashString(str string, algo crypto.Hash) []byte {
	h := algo.New()
	h.Write([]byte(str))
	return h.Sum(nil)
}

func HashStringToHex(str string, algo crypto.Hash) string {
	return hex.EncodeToString(HashString(str, algo))
}

func HashStringToBase64(str string, algo crypto.Hash) string {
	return base64.RawURLEncoding.EncodeToString(HashString(str, algo))
}
