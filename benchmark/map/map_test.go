package _map

import (
	"fmt"
	"github.com/cornelk/hashmap"
	"github.com/mmathys/acfts/common"
	"os"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkSyncMap(b *testing.B) {
	limit := 32768 // 2^15

	for size := 8; size <= limit; size *= 2 {
		name := fmt.Sprintf("size: %d", size)
		b.Run(name, func(b *testing.B) {
			utxos := hashmap.New(uintptr(size))
			var identifiers []common.Identifier

			for i := 0; i < b.N; i++ {
				identifiers = append(identifiers, common.RandomIdentifier())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if newItem := utxos.Insert(identifiers[i], true); !newItem {
					b.Error("should not happen")
				}
			}
		})
	}
}

func BenchmarkParallelMap(b *testing.B) {
	utxos := hashmap.New(uintptr(b.N))
	lastParam := os.Args[len(os.Args)-1]
	numWorkers := 2
	if paramNum, err := strconv.Atoi(lastParam); err == nil {
		numWorkers = paramNum
	}

	fmt.Printf("numWorkers = %d\n", numWorkers)

	var identifiers [][]common.Identifier
	for i := 0; i < numWorkers; i++ {
		identifiers = append(identifiers, []common.Identifier{})
		for j := 0; j < b.N/numWorkers; j++ {
			identifiers[i] = append(identifiers[i], common.RandomIdentifier())
		}
	}

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(work []common.Identifier) {
			for j := 0; j < len(work); j++ {
				if newItem := utxos.Insert(work[j], true); !newItem {
					b.Error("should not happen")
				}
			}
			wg.Done()
		}(identifiers[i])
	}

	wg.Wait()
}

func BenchmarkParallelMapGo(b *testing.B) {
	utxos := sync.Map{}
	lastParam := os.Args[len(os.Args)-1]
	numWorkers := 4
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
				if _, loaded := utxos.LoadOrStore(work[j], true); loaded {
					b.Error("should not happen")
				}
			}
			wg.Done()
		}(identifiers[i])
	}

	wg.Wait()
}

func TestStuff(t *testing.T) {
	utxos := hashmap.HashMap{}
	id := common.RandomIdentifier()
	newItem := utxos.Insert(id, true)
	if !newItem {
		t.Error("fail: should not exist")
	}
	newItem = utxos.Insert(id, true)
	if newItem {
		t.Error("fail: should exist")
	}
}
