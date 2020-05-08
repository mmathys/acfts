package test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/mmathys/acfts/common"
	"sync"
	"testing"
)

/**
Benchmark for crypto ops
*/

func BenchmarkSign(b *testing.B) {
	key := common.GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)

	numWorkers := 8

	b.ResetTimer()

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < b.N/numWorkers; j++ {
				ecdsa.Sign(rand.Reader, key, hash)
			}
			wg.Done()
		}()
	}

	wg.Wait()
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
