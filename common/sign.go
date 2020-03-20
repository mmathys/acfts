package common

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"math"
)

func signHash(key *ecdsa.PrivateKey, hash []byte) (ECDSASig, error) {
	r, s, err := ecdsa.Sign(rand.Reader, key, hash)
	if err != nil {
		return ECDSASig{}, err
	}
	return ECDSASig{R: r, S: s}, nil
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
- Checks whether there are enough signatures to satisfy the validity constraint. (> 2/3 of all sigs)
*/
func VerifyValue(key *ecdsa.PrivateKey, value *Value) error {
	hash := HashValue(*value)
	origins := make(map[string]bool)
	numSigs := 0

	for _, sig := range value.Signatures {
		valid := ecdsa.Verify(&key.PublicKey, hash, sig.R, sig.S)
		if !valid {
			return errors.New("verification failed")
		}

		// look out for duplicates
		encoded := string(crypto.FromECDSAPub(&key.PublicKey))
		if origins[encoded] {
			return errors.New("duplicate signatures")
		}
		origins[encoded] = true
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
