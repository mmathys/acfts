package _map

import (
	"fmt"
	"github.com/cornelk/hashmap"
	"github.com/mmathys/acfts/common"
	"gotest.tools/assert"
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
			var identifiers [][]byte

			for i := 0; i < b.N; i++ {
				id := common.RandomIdentifier()
				identifiers = append(identifiers, id[:])
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
	numWorkers := 4
	if paramNum, err := strconv.Atoi(lastParam); err == nil {
		numWorkers = paramNum
	}

	fmt.Printf("numWorkers = %d\n", numWorkers)

	var identifiers [][][]byte
	for i := 0; i < numWorkers; i++ {
		identifiers = append(identifiers, [][]byte{})
		for j := 0; j < b.N/numWorkers; j++ {
			id := common.RandomIdentifier()
			identifiers[i] = append(identifiers[i], id[:])
		}
	}

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(work [][]byte) {
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
	numWorkers := 64
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
	newItem := utxos.Insert(id[:], true)
	if !newItem {
		t.Error("fail: should not exist")
	}
	newItem = utxos.Insert(id[:], true)
	if newItem {
		t.Error("fail: should exist")
	}
}
/**
Fun Set
 */

/**
Basic Test
*/
func TestFunSetBasic(t *testing.T) {
	set := NewFunSet()
	var data [32]byte

	inserted := set.Insert(data)
	assert.Assert(t, inserted)
	inserted = set.Insert(data)
	assert.Assert(t, !inserted)
	inserted = set.Insert(data)
	assert.Assert(t, !inserted)
}

/**
This tests that no race conditions exist when inserting an element.
*/
func TestFunSetRaceCondition(t *testing.T) {
	numWorkers := 10
	id := common.RandomIdentifier()

	for i := 0; i < 10; i++ {
		set := NewFunSet()
		var wg sync.WaitGroup
		wg.Add(numWorkers)

		anyInsert := false
		for i := 0; i < numWorkers; i++ {
			go func() {
				inserted := set.Insert(id)
				if inserted {
					if anyInsert {
						t.Errorf("saw more than one insert")
					}
					anyInsert = true
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}


/**
This tests whether 1 million identifiers can be inserted into Fun Set
*/
func TestFunSet100MillionInserts(t *testing.T) {
	set := NewFunSet()

	for i := 0; i < 100e6; i++ {
		id := common.RandomIdentifier()
		inserted := set.Insert(id)
		if !inserted {
			t.Error("failed to insert")
		}
	}
}


func BenchmarkFunSetSingleIdentifier(b *testing.B) {
	utxos := NewFunSet()

	lastParam := os.Args[len(os.Args)-1]
	numWorkers := 64

	if paramNum, err := strconv.Atoi(lastParam); err == nil {
		numWorkers = paramNum
	}

	fmt.Printf("numWorkers = %d\n", numWorkers)

	id := common.RandomIdentifier()

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < b.N / numWorkers; j++ {
				utxos.Insert(id)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}


func BenchmarkFunSet(b *testing.B) {
	utxos := NewFunSet()

	lastParam := os.Args[len(os.Args)-1]
	numWorkers := 64

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

