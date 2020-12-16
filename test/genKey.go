package main

import (
	"fmt"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: genKey.go [keyLength]")
		os.Exit(0)
	}
	keyLength, err := strconv.Atoi(os.Args[1])

	if err != nil || keyLength <= 0 {
		fmt.Println("keyLength must be an positiveinteger")
		os.Exit(0)
	}
	value, _ := util.GenerateRandomKey(keyLength)
	fmt.Println(util.Base64Encode(value))
}
