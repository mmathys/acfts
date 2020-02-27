package util

import (
	"github.com/mmathys/acfts/common"
	"math/rand"
	"time"
)

// creates test wallet with 100 money
func NewWallet(addr common.Address) *common.Wallet {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	key := r1.Int()

	utxo := map[int]common.Tuple{
		key: {common.Address{0}, 100, key},
	}

	return &common.Wallet{addr, utxo}
}
