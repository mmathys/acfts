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
	utxos := hashmap.HashMap{}
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
}

func BenchmarkParallelMap(b *testing.B) {
	utxos := hashmap.HashMap{}
	lastParam := os.Args[len(os.Args) - 1]
	numWorkers := 1
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