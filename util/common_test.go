package util_test

import (
	"crypto"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestHashStr(t *testing.T) {
	digest := util.HashStringToHex("sha1 this string", crypto.SHA1)
	assert.Equal(t, "cf23df2207d99a74fbe169e3eba035e633b65d94", digest)
	digest = util.HashStringToBase64("sha1 this string", crypto.SHA1)
	assert.Equal(t, "zyPfIgfZmnT74Wnj66A15jO2XZQ", digest)
}

func benchmarkHashToHex(b *testing.B, n int) {
	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(n))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := strings.Repeat("A", n)
			util.HashStringToHex(s, crypto.SHA256)
		}
	})
}

func BenchmarkHashToHex5(b *testing.B)     { benchmarkHashToHex(b, 5) }
func BenchmarkHashToHex16(b *testing.B)    { benchmarkHashToHex(b, 16) }
func BenchmarkHashToHex64(b *testing.B)    { benchmarkHashToHex(b, 64) }
func BenchmarkHashToHex1024(b *testing.B)  { benchmarkHashToHex(b, 1024) }
func BenchmarkHashToHex65536(b *testing.B) { benchmarkHashToHex(b, 65536) }

func benchmarkHashToBase64(b *testing.B, n int) {
	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(n))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := strings.Repeat("A", n)
			util.HashStringToBase64(s, crypto.SHA256)
		}
	})
}

func BenchmarkHashToBase64_5(b *testing.B)     { benchmarkHashToBase64(b, 5) }
func BenchmarkHashToBase64_16(b *testing.B)    { benchmarkHashToBase64(b, 16) }
func BenchmarkHashToBase64_64(b *testing.B)    { benchmarkHashToBase64(b, 64) }
func BenchmarkHashToBase64_1024(b *testing.B)  { benchmarkHashToBase64(b, 1024) }
func BenchmarkHashToBase64_65536(b *testing.B) { benchmarkHashToBase64(b, 65536) }
