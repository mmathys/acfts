package common

import (
	"crypto/ecdsa"
	"math/big"
	"sync"
	"time"
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
	Key			*ecdsa.PrivateKey
}

type Wallet struct {
	*Identity
	UTXO 		*sync.Map // of type int --> Value
}

type NodeType string		// "server" | "client"

type Node struct {
	NodeType 	NodeType
	Address  	Address				// internal address
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