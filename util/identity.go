package util

import (
	crypto2 "github.com/ethereum/go-ethereum/crypto"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"github.com/mmathys/acfts/crypto"
	"math/rand"
	"sync"
	"time"
)

func GetIdentity(alias common.Alias) *common.Identity {
	key := core.GetKey(alias)
	id := common.Identity{Alias: alias, Key: key}
	return &id
}

func NewWalletWithAmount(alias common.Alias, value int) *common.Wallet {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	key := r1.Int()
	id := GetIdentity(alias)

	var utxo sync.Map

	addr := crypto2.FromECDSAPub(&id.Key.PublicKey)
	v := common.Value{Address: addr, Amount: value, Id: key}

	// this is for debugging. In production, there would be an origin output or something like that
	for _, server := range core.GetServers() {
		key := core.GetKey(server)
		err := crypto.SignValue(key, &v)
		if err != nil {
			panic(err)
		}
	}

	utxo.Store(key, v)

	return &common.Wallet{Identity: id, UTXO: &utxo}
}

// creates test wallet with 100 money
func NewWallet(alias common.Alias) *common.Wallet {
	return NewWalletWithAmount(alias, 100)
}
