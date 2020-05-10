package common

import (
	"sync"
)

func GetIdentity(address Address) *Identity {
	key := GetKey(address)
	id := Identity{Key: key, Address: address}
	return &id
}

func NewWalletWithAmount(address Address, value int) *Wallet {
	utxoId := RandomIdentifier()
	id := GetIdentity(address)

	var utxo sync.Map

	addr := MarshalPubkey(&id.Key.PublicKey)
	v := Value{Address: addr, Amount: value, Id: utxoId}

	// every client gets valid 100 credits to their account.
	// this is for debugging. In production, there would be an origin output or something like that
	for _, server := range GetServers() {
		key := GetKey(server)
		err := SignValue(key, &v)
		if err != nil {
			panic(err)
		}
	}

	// calculate the shardIndex, which is static

	index := [IdentifierLength]byte{}
	copy(index[:], utxoId[:IdentifierLength])

	utxo.Store(index, v)

	return &Wallet{Identity: id, UTXO: &utxo}
}