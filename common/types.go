package common

import (
	"crypto/ecdsa"
	"math/big"
)

const (
	AddressLength = 1
)

type Address [AddressLength]byte

type ECDSASig struct {
	R		*big.Int
	S		*big.Int
}

// Defines an Input / Output tuple; with extra fields
type Value struct {
	Address 	Address 	// The address
	Amount  	int     	// The value itself
	Id      	int     	// Unique identifier
	Signatures	[]ECDSASig	// Signatures
}

type Transaction struct {
	Inputs 		[]Value
	Outputs 	[]Value
}

type TransactionSigRequest struct {
	Inputs 		[]Value
	Outputs 	[]Value
}

type TransactionSignRes struct {
	Outputs		[]Value
}

type Identity struct {
	Address		Address
	Key			ecdsa.PrivateKey
}

type Wallet struct {
	Identity
	UTXO 		map[int]Value
}