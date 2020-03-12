package util

import (
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"math/rand"
	"time"
)

func GetIdentity(addr common.Address) *common.Identity {
	key := core.GetKey(addr)
	id := common.Identity{Address: addr, Key: key}
	return &id
}

func NewWalletWithAmount(addr common.Address, value int) *common.Wallet {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	key := r1.Int()

	utxo := map[int]common.Value{
		key: {addr, value, key, nil},
	}

	id := GetIdentity(addr)
	return &common.Wallet{Identity: id, UTXO: utxo}
}

// creates test wallet with 100 money
func GetWallet(addr common.Address) *common.Wallet {
	return NewWalletWithAmount(addr, 100)
}
