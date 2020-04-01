package benchmark

import (
	"sync"
	"testing"
)

var SignedUTXO sync.Map

func BenchmarkSyncMap(b *testing.B) {
	//index := [common.IdentifierLength]byte{}

	for i := 0; i < b.N; i++ {
		SignedUTXO.LoadOrStore(i, true)
	}
}