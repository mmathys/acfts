package common

import (
	"sync"
)

func NewWalletWithAmount(address Address, value int) *Wallet {
	utxoId := RandomIdentifier()
	key := GetKey(address)

	var utxo sync.Map

	v := Value{Address: key.GetAddress(), Amount: value, Id: utxoId}

	// every client gets valid 100 credits to their account.
	// this is for debugging. In production, there would be an origin output or something like that
	for _, server := range GetServers() {
		key := GetKey(server)
		err := key.SignValue(&v)
		if err != nil {
			panic(err)
		}
	}

	// calculate the shardIndex, which is static

	index := [IdentifierLength]byte{}
	copy(index[:], utxoId[:IdentifierLength])

	utxo.Store(index, v)

	return &Wallet{Key: key, UTXO: &utxo}
}