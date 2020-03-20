package benchmark

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/mmathys/acfts/common"
	"testing"
)

/**
Benchmark for crypto ops
 */

func BenchmarkSign(b *testing.B) {
	key := common.GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecdsa.Sign(rand.Reader, key, hash)
	}
}

func BenchmarkHashValue(b *testing.B) {
	//TODO
}

func BenchmarkVerify(b *testing.B) {
	key := common.GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)
	r, s, _ := ecdsa.Sign(rand.Reader, key, hash)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecdsa.Verify(&key.PublicKey, hash, r, s)
	}
}
