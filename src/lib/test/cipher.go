package main

import (
	"fmt"
	"template2/lib/util"
)

var (
	key       = "pR6kDdHYqNMRO74kUxFiGgpv3A6qKFeCY6IDHxDH8NY"
	cipher, _ = util.NewAesGcmCipher(key)
	plainText = "hello world"
)

func main() {
	fmt.Println("PlainText:", plainText)
	cipherText, _ := cipher.Encrypt(plainText)
	fmt.Println("CipherText:", cipherText)
	decryptText, _ := cipher.Decrypt(cipherText)
	fmt.Println("DecryptText:", decryptText)
}
