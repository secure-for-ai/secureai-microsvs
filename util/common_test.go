package util_test

import (
	"crypto"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashStr(t *testing.T) {
	digest := util.HashStringToHex("sha1 this string", crypto.SHA1)
	assert.Equal(t, "cf23df2207d99a74fbe169e3eba035e633b65d94", digest)
	digest = util.HashStringToBase64("sha1 this string", crypto.SHA1)
	assert.Equal(t, "zyPfIgfZmnT74Wnj66A15jO2XZQ", digest)
}
