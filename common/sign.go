package common

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/oasislabs/ed25519"
	"math"
)

/**
Signing
*/

// Signs a hash
func SignHash(id *Identity, hash []byte) *EdDSASig {
	sig := ed25519.Sign(id.PrivateKey, hash)
	return &EdDSASig{
		Address:   id.Address,
		Signature: sig,
	}
}

// Signs a value
func SignValue(id *Identity, value *Value) error {
	hash := HashValue(*value)

	if value.Signatures == nil {
		value.Signatures = []EdDSASig{}
	}

	sig := SignHash(id, hash)
	value.Signatures = append(value.Signatures, *sig)
	return nil
}

// Signs multiple values
func SignValues(id *Identity, outputs []Value) ([]Value, error) {
	var signed []Value

	for _, i := range outputs {
		SignValue(id, &i)
		signed = append(signed, i)
	}

	return signed, nil
}

// Signs transaction signature request, which is requested by a client
func SignTransactionSigRequest(id *Identity, request *TransactionSigReq) error {
	hash := HashTransactionSigRequest(*request)
	sig := SignHash(id, hash)
	request.Signature = *sig

	return nil
}

/**
Signature Verification
*/

// Verifies a signature. Using secp256k1 C bindings crypto
func Verify(sig *EdDSASig, hash []byte) (bool, error) {
	if len(sig.Signature) != SignatureLength {
		msg := fmt.Sprintf("invalid signature length. wanted: %d, got: %d", SignatureLength, len(sig.Signature))
		return false, errors.New(msg)
	}
	return ed25519.Verify(sig.Address, hash, sig.Signature), nil
}

// Verifies single value
// - Verifies all signatures
// - Checks whether there are duplicate signatures
// - Checks whether the signature are from valid severs
// - Checks whether there are enough signatures to satisfy the validity constraint. (> 2/3 of all sigs)
func VerifyValue(value *Value) error {
	hash := HashValue(*value)
	origins := make(map[[AddressLength]byte]bool)
	numSigs := 0

	for _, sig := range value.Signatures {
		valid, err := Verify(&sig, hash)
		if err != nil {
			return err
		}

		if !valid {
			return errors.New("value verification failed")
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

// Verifies a signature request
// - checks if all inputs are owned by the same party
// - checks if party signed the request
func VerifyTransactionSigRequest(req *TransactionSigReq) error {
	hash := HashTransactionSigRequest(*req)

	ownerAddress := req.Signature.Address

	for _, input := range req.Transaction.Inputs {
		if !bytes.Equal(ownerAddress, input.Address) {
			return errors.New("inputs are not owned by the same party")
		}
	}

	valid, err := Verify(&req.Signature, hash)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("sig request verification failed")
	}

	return nil
}
