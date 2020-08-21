package _map

import (
	"fmt"
	"github.com/cornelk/hashmap"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/common/funset"
	"gotest.tools/assert"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
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
	set := funset.NewFunSet()
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
		set := funset.NewFunSet()
		var wg sync.WaitGroup
		wg.Add(numWorkers)

		anyInsert := false
		for i := 0; i < numWorkers; i++ {
			go func() {
				inserted := set.Insert(id)
				if inserted {
					if anyInsert {
						t.Errorf("saw more than one insertBenchmark")
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
func TestFunSetInserts(t *testing.T) {
	var N int = 10e8

	insertMemoryTest(t, 64, N)
}

// this function first prepares the transactions, then inserts them
func insertMemoryTest(t *testing.T, numWorkers int, overrideN int) {
	fmt.Println("initializing set...")
	utxos := funset.NewFunSet()
	fmt.Println("done initializing set. Sleeping for 5 seconds")
	time.Sleep(5 * time.Second)
	fmt.Printf("numWorkers = %d\n", numWorkers)

	N := overrideN
	fmt.Printf("N = %d\n", N)

	var wg sync.WaitGroup

	fmt.Println("inserting transactions...")
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(amount int) {
			for j := 0; j < amount; j++ {
				id := common.RandomIdentifier()
				if !utxos.Insert(id) {
					t.Fatal("should not happen")
				}
			}
			wg.Done()
		}(N/numWorkers)
	}

	wg.Wait()

	fmt.Println("finished. Sleeping 5 seconds before exiting")
	time.Sleep(5 * time.Second)
}

func BenchmarkFunSetSingleIdentifier(b *testing.B) {
	utxos := funset.NewFunSet()

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
			for j := 0; j < b.N/numWorkers; j++ {
				utxos.Insert(id)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

// this is likely not gonna work (OOM)
func BenchmarkFunSet(b *testing.B) {
	lastParam := os.Args[len(os.Args)-1]
	numWorkers := 64

	if paramNum, err := strconv.Atoi(lastParam); err == nil {
		numWorkers = paramNum
	}

	insertBenchmark(b, numWorkers, -1)
}

// this function first prepares the transactions, then inserts them
func insertBenchmark(b *testing.B, numWorkers int, overrideN int) {
	fmt.Println("initializing set...")
	utxos := funset.NewFunSet()
	fmt.Println("done initializing set.")
	fmt.Printf("numWorkers = %d\n", numWorkers)

	N := 0
	if b == nil {
		N = overrideN
	} else {
		N = b.N
	}

	fmt.Printf("N = %d\n", N)

	var ig sync.WaitGroup
	fmt.Println("initializing transactions...")
	var identifiers [][][common.IdentifierLength]byte
	for i := 0; i < numWorkers; i++ {
		identifiers = append(identifiers, [][common.IdentifierLength]byte{})
		ig.Add(1)
		go func(i int) {
			for j := 0; j < N/numWorkers; j++ {
				id := common.RandomIdentifier()
				identifiers[i] = append(identifiers[i], id)
			}
			//fmt.Printf("worker %d done\n", i)
			ig.Done()
		}(i)
	}

	ig.Wait()

	var wg sync.WaitGroup

	fmt.Println("inserting transactions...")
	start := time.Now()
	if b != nil {
		b.ResetTimer()
	}
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

	if b == nil {
		end := time.Now()
		elapsed := end.Sub(start)
		fmt.Printf("executed %v inserts in %v\n", N, float64(elapsed)/float64(time.Second))
		txps := float64(N) / (float64(elapsed) / float64(time.Second))
		fmt.Printf("= %v tx/s\n", txps)
	}
}
