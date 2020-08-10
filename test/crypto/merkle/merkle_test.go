package merkle

import (
	"crypto/rand"
	"fmt"
	"github.com/mmathys/acfts/common"
	"os"
	"strconv"
	"testing"
)

const defaultPoolSize = 32

func hashes(numLeaves int) [][]byte {
	// prepare hashes
	var hashes [][]byte
	for i := 0; i < numLeaves; i++ {
		hash := make([]byte, 64)
		rand.Read(hash)
		hashes = append(hashes, hash)
	}
	return hashes
}

func TestSignMerkle(t *testing.T) {
	const NumLeaves = 32
	key := common.GenerateKey(common.ModeMerkle, 0)

	hashes := hashes(NumLeaves)
	sigs := key.SignMultipleMerkle(hashes)
	fmt.Println(len(sigs))
}

func TestVerifyMerkle(t *testing.T) {
	const NumLeaves = 32
	key := common.GenerateKey(common.ModeMerkle, 0)
	hashes := hashes(NumLeaves)
	sigs := key.SignMultipleMerkle(hashes)

	for i := range sigs {
		valid, err := common.Verify(sigs[i], hashes[i])
		if err != nil {
			t.Fatal(err)
		}
		if !valid {
			t.Fatal("signature not valid")
		}
	}
}

func getPoolSize() int {
	poolSize := defaultPoolSize
	arg := os.Args[len(os.Args)-1]
	if parsed, err := strconv.Atoi(arg); err == nil {
		poolSize = parsed
	}

	return poolSize
}

func BenchmarkPoolSize(b *testing.B) {
	fmt.Println("pool size:", getPoolSize())
}

func BenchmarkSignMerkle(b *testing.B) {
	poolSize := getPoolSize()
	hashes := hashes(poolSize)
	key := common.GenerateKey(common.ModeMerkle, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = key.SignMultipleMerkle(hashes)
	}
}

// only verify a single merkle signature.
func BenchmarkVerifyMerkle(b *testing.B) {
	poolSize := getPoolSize()
	hashes := hashes(poolSize)
	key := common.GenerateKey(common.ModeMerkle, 0)
	sigs := key.SignMultipleMerkle(hashes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		common.Verify(sigs[0], hashes[0])
	}
}
