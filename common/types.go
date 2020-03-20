package common

import (
	"crypto/ecdsa"
	"math/big"
	"sync"
	"time"
)

const (
	AddressLength = 65 // address = public key
	PrivateKeyLength = 32
)

type Address = [AddressLength]byte
type PrivateKey = [PrivateKeyLength]byte

type ECDSASig struct {
	R		*big.Int
	S		*big.Int
}

// Defines an Input / Output tuple; with extra fields
type Value struct {
	Address 	Address			// The public key = Address (encoded)
	Amount  	int     		// The value itself
	Id      	int     		// Unique identifier
	Signatures	[]ECDSASig		// Signatures
}

type Transaction struct {
	Inputs 		[]Value
	Outputs 	[]Value
}

type TransactionSigReq struct {
	Transaction	Transaction
	Signature	ECDSASig
}

type TransactionSignRes struct {
	Outputs		[]Value
}

type Identity struct {
	Address		Address
	Key			*ecdsa.PrivateKey
}

type Wallet struct {
	*Identity
	UTXO 		*sync.Map 	// of type int --> Value
}

type NodeType string		// "server" | "client"

type Node struct {
	NodeType 	NodeType
	Net      	string				// network address, with http
	Port     	int					// port
	Key      	*ecdsa.PrivateKey
}

type Agent struct {
	NumTransactions	int				// how many tx the agent completes before exiting
	StartDelay		time.Duration	// delay before starting transactions in ns (waiting for other agents to launch)
	EndDelay		time.Duration	// delay finishing (for receiving stuff tx)
	Address			Address			// reference to node
	Topology		[]Address		// other nodes (exluding self)
}