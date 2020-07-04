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
	if key.Mode == ModeEdDSA ||key.Mode == ModeMerkle {
		// EdDSA/Merkle: sign the 100 credits by each server
		for _, server := range GetServers() {
			key := GetKey(server)
			err := key.SignValue(&v)
			if err != nil {
				panic(err)
			}
		}
	} else if key.Mode == ModeBLS {
		// BLS: sign with master key
		key := GetBLSMasterKey()
		err := key.SignValue(&v)
		if err != nil {
			panic(err)
		}
	} else {
		panic("unsupported mode")
	}

	index := [IdentifierLength]byte{}
	copy(index[:], utxoId[:IdentifierLength])

	utxo.Store(index, v)

	return &Wallet{Key: key, UTXO: &utxo}
}
