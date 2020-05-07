package _map

import (
	"fmt"
	"github.com/mmathys/acfts/common"
	"os"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkFunSet(b *testing.B) {
	utxos := NewFunSet()
	lastParam := os.Args[len(os.Args)-1]
	numWorkers := 8
	if paramNum, err := strconv.Atoi(lastParam); err == nil {
		numWorkers = paramNum
	}

	fmt.Printf("numWorkers = %d\n", numWorkers)

	var identifiers [][][common.IdentifierLength]byte
	for i := 0; i < numWorkers; i++ {
		identifiers = append(identifiers, [][common.IdentifierLength]byte{})
		for j := 0; j < b.N/numWorkers; j++ {
			id := common.RandomIdentifier()
			array := [common.IdentifierLength]byte{}
			copy(array[:], id[:common.IdentifierLength])
			identifiers[i] = append(identifiers[i], array)
		}
	}

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(work [][common.IdentifierLength]byte) {
			for j := 0; j < len(work); j++ {
				if !utxos.Insert(work[j]) {
					b.Error("should not happen")
				}
			}
			wg.Done()
		}(identifiers[i])
	}

	wg.Wait()
}
