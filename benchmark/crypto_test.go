package benchmark

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"github.com/mmathys/acfts/crypto"
	"github.com/mmathys/acfts/util"
	"testing"
)

/**
Benchmark for crypto ops
 */

func BenchmarkSign(b *testing.B) {
	key := crypto.GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecdsa.Sign(rand.Reader, key, hash)
	}
}

func BenchmarkHashValue(b *testing.B) {
	addr := common.Address{0}
	w := util.GetWallet(addr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.DoHash(w.UTXO[0])
	}
}

func BenchmarkVerify(b *testing.B) {
	key := crypto.GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)
	r, s, _ := ecdsa.Sign(rand.Reader, key, hash)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecdsa.Verify(&key.PublicKey, hash, r, s)
	}
}
