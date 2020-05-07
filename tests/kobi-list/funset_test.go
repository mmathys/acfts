package kobi_list

import (
	"fmt"
	"github.com/mmathys/acfts/common"
	"gotest.tools/assert"
	"os"
	"strconv"
	"sync"
	"testing"
)

/**
Basic Test
*/
func TestBasic(t *testing.T) {
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
func TestRaceCondition(t *testing.T) {
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
This tests whether 1 million identifiers can be inserted (not a benchmark)
*/
func Test100MillionInserts(t *testing.T) {
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
