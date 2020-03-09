package util

import (
	"github.com/mmathys/acfts/common"
	"math/rand"
	"time"
)

func NewWalletWithAmount(addr common.Address, value int) *common.Wallet {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	key := r1.Int()

	utxo := map[int]common.Tuple{
		key: {addr, value, key},
	}

	return &common.Wallet{Address: addr, UTXO: utxo}
}

// creates test wallet with 100 money
func NewWallet(addr common.Address) *common.Wallet {
	return NewWalletWithAmount(addr, 100)
}
