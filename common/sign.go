package common

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
	"log"
	"math"
)

func signHash(key *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	sig, err := crypto2.Sign(hash, key)
	if err != nil {
		return []byte{}, err
	}
	return sig, nil
}

func SignValue(key *ecdsa.PrivateKey, value *Value) error {
	hash := HashValue(*value)

	if value.Signatures == nil {
		value.Signatures = [][]byte{}
	}

	sig, err := signHash(key, hash)
	if err != nil {
		return err
	}

	value.Signatures = append(value.Signatures, sig)
	return nil
}

func SignValues(key *ecdsa.PrivateKey, outputs []Value) ([]Value, error) {
	var signed []Value

	for _, i := range outputs {
		SignValue(key, &i)
		signed = append(signed, i)
	}

	return signed, nil
}

func SignTransactionSigRequest(key *ecdsa.PrivateKey, request *TransactionSigReq) error {
	hash := HashTransactionSigRequest(*request)
	sig, err := signHash(key, hash)
	if err != nil {
		return err
	}

	request.Signature = sig
	return nil
}

/**
Verifies single value
- Verifies all signatures
- Checks whether there are duplicate signatures
- Checks whether the signature are from valid severs
- Checks whether there are enough signatures to satisfy the validity constraint. (> 2/3 of all sigs)
*/
func VerifyValue(value *Value) error {
	hash := HashValue(*value)
	origins := make(map[[AddressLength]byte]bool)
	numSigs := 0

	for _, sig := range value.Signatures {
		pubkey, err := crypto2.Ecrecover(hash, sig)
		if err != nil {
			return err
		}

		sig = sig[:len(sig)-1] // remove recovery bit
		valid := crypto2.VerifySignature(pubkey, hash, sig)

		if !valid {
			return errors.New("value verification failed")
		}

		// look out for duplicates signatures
		index := [AddressLength]byte{}
		copy(index[:], pubkey[:AddressLength])
		if origins[index] {
			return errors.New("duplicate signatures")
		}
		origins[index] = true
		numSigs++
	}

	numServers := GetNumServers()
	numRequiredSigs := int(math.Ceil(2.0 / 3.0 * float64(numServers)))

	if numSigs < numRequiredSigs {
		text := fmt.Sprintf("not enough signatures. need %d, have %d", int(numRequiredSigs), numSigs)
		return errors.New(text)
	}

	return nil
}

/*
Checks validity of a *completed* transaction. It's only used in the client.
- verifies inputs and outputs
*/
func VerifyTransaction(value *Transaction) error {
	return nil
}

/*
Verifies a signature request
- checks if all inputs are owned by the same party
- checks if party signed the request
*/
func VerifyTransactionSigRequest(req *TransactionSigReq) error {
	hash := HashTransactionSigRequest(*req)

	owner, err := crypto2.Ecrecover(hash, req.Signature)
	if err != nil {
		panic(err)
	}

	for _, input := range req.Transaction.Inputs {
		if !bytes.Equal(owner, input.Address) {
			return errors.New("inputs are not owned by the same party")
		}
	}

	sig := req.Signature[:len(req.Signature)-1] // remove recovery bit
	valid := crypto2.VerifySignature(owner, hash, sig)

	if !valid {
		return errors.New("sig request verification failed")
	}

	return nil
}

func RecoverPubkeyBytes(hash []byte, sig []byte) []byte {
	recoveredPub, err := crypto2.Ecrecover(hash, sig)
	if err != nil {
		log.Panicf("ECRecover error: %s", err)
	}
	return recoveredPub
}

func RecoverPubkey(hash []byte, sig []byte) *ecdsa.PublicKey {
	recoveredPub := RecoverPubkeyBytes(hash, sig)
	return UnmarshalPubkey(recoveredPub)
}