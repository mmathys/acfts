package benchmark

import (
	"math/rand"
	"testing"
	"time"
)

/**
Benchmark for crypto ops
 */

func BenchmarkGenerateKeys(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	for i := 1; i < b.N; i++ {

	}
}