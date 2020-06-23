package common

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/oasislabs/ed25519"
	"sync"
	"time"
)

const (
	EdDSAPublicKeyLength  = 32 // address = public key
	EdDSAPrivateKeyLength = 64
	BLSPublicKeyLength    = 48 // encoded
	BLSPrivateKeyLength   = 32
	IdentifierLength      = 32 // used for UTXOs
	SignatureLength       = 64
	ModeEdDSA             = 1
	ModeBLS               = 2
)

// Keys
type EdDSAKey struct {
	Address    ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

type BLSKey struct {
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
	Address        Address
	EdDSASignature *[]byte
	BLSSignature   *bls.Sign
	Mode           int
}

type Address = ed25519.PublicKey         // len = EdDSAPublicKeyLength
type PrivateKey = ed25519.PrivateKey     // len = EdDSAPrivateKeyLength
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
