package util

import (
	"github.com/mmathys/acfts/common"
	"math/rand"
	"sync"
	"time"
)

func GetIdentity(address common.Address) *common.Identity {
	key := common.GetKey(address)
	id := common.Identity{Key: key, Address: address}
	return &id
}

func NewWalletWithAmount(address common.Address, value int) *common.Wallet {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	utxoId := r1.Int()
	id := GetIdentity(address)

	var utxo sync.Map

	addr := common.MarshalPubkey(&id.Key.PublicKey)
	v := common.Value{Address: addr, Amount: value, Id: utxoId}

	// this is for debugging. In production, there would be an origin output or something like that
	for _, server := range common.GetServers() {
		key := common.GetKey(server)
		err := common.SignValue(key, &v)
		if err != nil {
			panic(err)
		}
	}

	utxo.Store(utxoId, v)

	return &common.Wallet{Identity: id, UTXO: &utxo}
}

// creates test wallet with 100 money
func NewWallet(address common.Address) *common.Wallet {
	return NewWalletWithAmount(address, 100)
}
