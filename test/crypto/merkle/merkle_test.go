package merkle

import (
	"crypto/rand"
	"fmt"
	"github.com/mmathys/acfts/common"
	"testing"
)

func sign(key *common.Key, numLeaves int) (hashes [][]byte, sigs []*common.Signature) {

	// prepare hashes
	hashes = make([][]byte, numLeaves, 64)
	for i := 0; i < numLeaves; i++ {
		rand.Read(hashes[i])
	}

	// sign
	sigs = key.SignMultipleMerkle(hashes)
	return
}

func TestSignMerkle(t *testing.T) {
	const NumLeaves = 32
	key := common.GenerateKey(common.ModeMerkle, 0)

	_, sigs := sign(key, NumLeaves)
	fmt.Println(len(sigs))
}

func TestVerifyMerkle(t *testing.T) {
	const NumLeaves = 32
	key := common.GenerateKey(common.ModeMerkle, 0)
	hashes, sigs := sign(key, NumLeaves)

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
