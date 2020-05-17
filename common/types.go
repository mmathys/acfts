package common

import (
	"github.com/oasislabs/ed25519"
	"sync"
	"time"
)

const (
	AddressLength    = 32 // address = public key
	PrivateKeyLength = 64
	IdentifierLength = 32 // used for UTXOs
	SignatureLength = 64
)

type Address = ed25519.PublicKey         // len = AddressLength
type PrivateKey = ed25519.PrivateKey     // len = PrivateKeyLength
type Identifier = [IdentifierLength]byte // len = IdentifierLength

type EdDSASig struct {
	Address		Address
	Signature 	[]byte
}

// Defines an Input / Output tuple; with extra fields
type Value struct {
	Address    Address    // The public key = Address (encoded)
	Amount     int        // The value itself
	Id         Identifier // Unique identifier
	Signatures []EdDSASig // Signatures
}

type Transaction struct {
	Inputs  []Value
	Outputs []Value
}

type TransactionSigReq struct {
	Transaction Transaction
	Signature   EdDSASig
}

type TransactionSignRes struct {
	Outputs []Value
}

type Identity struct {
	Address		Address
	PrivateKey  PrivateKey
}

type Wallet struct {
	*Identity
	UTXO *sync.Map // of type int --> Value
}

/**
Topology configuration
*/

type Instance struct {
	Net  string // network address, with http
	Port int    // port
}

type ClientNode struct {
	Instance	Instance
	Key      	PrivateKey
	Balance		int
}

type ServerNode struct {
	Instances	[]Instance
	Key       	PrivateKey
}

type Agent struct {
	NumTransactions int           // how many tx the agent completes before exiting
	StartDelay      time.Duration // delay before starting transactions in ns (waiting for other agents to launch)
	Address         Address       // reference to node
	Topology        []Address     // other nodes (excluding self)
}
