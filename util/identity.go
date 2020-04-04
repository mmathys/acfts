package util

import (
	"github.com/mmathys/acfts/common"
	"sync"
)

var utxo sync.Map

func GetIdentity(address common.Address) *common.Identity {
	key := common.GetKey(address)
	id := common.Identity{Key: key, Address: address}
	return &id
}

func NewWalletWithAmount(address common.Address, value int) *common.Wallet {
	utxoId := common.RandomIdentifier()
	id := GetIdentity(address)

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

	index := [common.IdentifierLength]byte{}
	copy(index[:], utxoId[:common.IdentifierLength])

	utxo.Delete(index)
	utxo.Store(index, v)

	return &common.Wallet{Identity: id, UTXO: &utxo}
}

// creates test wallet with 100 money
func NewWallet(address common.Address) *common.Wallet {
	return NewWalletWithAmount(address, 100)
}
