package common

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/oasislabs/ed25519"
	"sync"
	"time"
)

const (
	// EdDSA
	ModeEdDSA             = 1
	EdDSAPublicKeyLength  = 32 // address = public key
	EdDSAPrivateKeyLength = 64

	// BLS
	ModeBLS             = 2
	BLSPublicKeyLength  = 48 // encoded
	BLSPrivateKeyLength = 32

	// Merkle
	ModeMerkle             = 3
	MerklePublicKeyLength  = 32
	MerklePrivateKeyLength = 64

	IdentifierLength = 32 // used for UTXOs
	SignatureLength  = 64
)

// Keys
type EdDSAKey struct {
	Address    ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

type BLSKey struct {
	ID         bls.ID
	Address    bls.PublicKey
	PrivateKey bls.SecretKey
}

type Key struct {
	EdDSA *EdDSAKey
	BLS   *BLSKey
	Mode  int
}

// Signatures

type Signature struct {
	Address   Address
	Signature []byte
	Mode      int

	// BLS specific
	BLSID []byte

	// Merkle specific
	MerklePath    [][]byte
	MerkleIndexes []bool
}

type Address = []byte                    // len = EdDSAPublicKeyLength
type PrivateKey = []byte                 // len = EdDSAPrivateKeyLength
type Identifier = [IdentifierLength]byte // len = IdentifierLength

// Defines an Input / Output tuple; with extra fields
type Value struct {
	Address    Address     // The public key = Address (encoded)
	Amount     int         // The value itself
	Id         Identifier  // Unique identifier
	Signatures []Signature // Signatures
}

type Transaction struct {
	Inputs  []Value
	Outputs []Value
}

type TransactionSigReq struct {
	Transaction Transaction
	Signature   Signature
}

type TransactionSignRes struct {
	Outputs []Value
}

type Wallet struct {
	*Key
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
	Instance Instance
	Key      *Key
	Balance  int
}

type ServerNode struct {
	Instances []Instance
	Key       *Key
}

type Agent struct {
	NumTransactions int           // how many tx the agent completes before exiting
	StartDelay      time.Duration // delay before starting transactions in ns (waiting for other agents to launch)
	Address         Address       // reference to node
	Topology        []Address     // other nodes (excluding self)
}
