package main

import (
	"fmt"
	"github.com/secure-for-ai/secureai-microsvs/webCrypto"
)

var (
	key       = "pR6kDdHYqNMRO74kUxFiGgpv3A6qKFeCY6IDHxDH8NY"
	cipher, _ = webCrypto.NewAesGcm(key)
	plainText = []byte("hello world")
)

func main() {
	fmt.Println("PlainText:", string(plainText))
	cipherText, _ := cipher.EncryptBase64(plainText)
	fmt.Println("CipherText:", cipherText)
	decryptText, _ := cipher.DecryptBase64(cipherText)
	fmt.Println("DecryptText:", string(decryptText))
}
