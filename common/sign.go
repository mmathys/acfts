package common

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
)

func signHash(key *ecdsa.PrivateKey, hash []byte) (ECDSASig, error) {
	r, s, err := ecdsa.Sign(rand.Reader, key, hash)
	if err != nil {
		return ECDSASig{}, err
	}
	addr := MarshalPubkey(&key.PublicKey)
	return ECDSASig{R: r, S: s, Address: addr}, nil
}

func SignValue(key *ecdsa.PrivateKey, value *Value) error {
	hash := HashValue(*value)

	if value.Signatures == nil {
		value.Signatures = []ECDSASig{}
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
		pubkey := UnmarshalPubkey(sig.Address)
		valid := ecdsa.Verify(pubkey, hash, sig.R, sig.S)
		if !valid {
			return errors.New("verification failed")
		}

		// look out for duplicates signatures
		index := [AddressLength]byte{}
		copy(index[:], sig.Address[:AddressLength])
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
func VerifyTransaction(key *ecdsa.PrivateKey, value *Transaction) error {
	return nil
}
