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

const (
	// if num sigs >= BatchVerificationThreshold, then batch verification is more efficient.
	BatchVerificationThreshold = 4
)

// Verifies a signature.
func Verify(sig *EdDSASig, hash []byte) (bool, error) {
	if len(sig.Signature) != SignatureLength {
		msg := fmt.Sprintf("invalid signature length. wanted: %d, got: %d", SignatureLength, len(sig.Signature))
		return false, errors.New(msg)
	}
	return ed25519.Verify(sig.Address, hash, sig.Signature), nil
}

// Performs batch verification
func VerifyBatch(sigs []EdDSASig, hash []byte) (bool, error) {
	var pks []Address
	var sigsByte [][]byte
	for _, sig := range sigs {
		pks = append(pks, sig.Address)

		if len(sig.Signature) != SignatureLength {
			msg := fmt.Sprintf("invalid signature length. wanted: %d, got: %d", SignatureLength, len(sig.Signature))
			return false, errors.New(msg)
		}
		sigsByte = append(sigsByte, sig.Signature)
	}

	var messages [][]byte
	for range sigs {
		messages = append(messages, hash)
	}

	var opts ed25519.Options
	ok, _, err := ed25519.VerifyBatch(nil, pks[:], messages[:], sigsByte[:], &opts)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// Verifies single value
// - Checks whether there are duplicate signatures
// - Checks whether the signature are from valid servers
// - Checks whether there are enough signatures to satisfy the validity constraint. (> 2/3 of all sigs)
// - Verifies all signatures
func VerifyValue(value *Value, enableBatchVerification bool) error {
	hash := HashValue(*value)
	origins := make(map[[AddressLength]byte]bool)

	// look out for duplicates signatures
	for _, sig := range value.Signatures {
		index := [AddressLength]byte{}
		copy(index[:], sig.Address[:AddressLength])
		if origins[index] {
			return errors.New("duplicate signatures")
		}
		origins[index] = true
	}

	// check whether the signatures have all been made by valid servers
	for _, sig := range value.Signatures {
		if !IsValidServer(sig.Address) {
			return errors.New("encountered signature signed by valid server")
		}
	}

	// check that there are enough signatures
	numServers := GetNumServers()
	numRequiredSigs := int(math.Ceil(2.0 / 3.0 * float64(numServers)))
	if len(value.Signatures) < numRequiredSigs {
		text := fmt.Sprintf("not enough signatures. need %d, have %d", numRequiredSigs, len(value.Signatures))
		return errors.New(text)
	}

	// verify all signatures, either with batch verification or single verification
	if enableBatchVerification && len(value.Signatures) >= BatchVerificationThreshold {
		// batch verification
		valid, err := VerifyBatch(value.Signatures, hash)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("value verification failed (batch mode)")
		}
	} else {
		// single verification
		for _, sig := range value.Signatures {
			valid, err := Verify(&sig, hash)
			if err != nil {
				return err
			}
			if !valid {
				return errors.New("value verification failed (single mode)")
			}
		}
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
