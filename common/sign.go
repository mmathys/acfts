package common

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"
	"github.com/oasislabs/ed25519"
)

/**
Signing
*/

// Signs a hash
func (key *Key) SignHash(hash []byte) *Signature {
	if key.Mode == ModeEdDSA {
		opts := ed25519.Options{
			Hash: crypto.SHA512,
		}
		sig, err := key.EdDSA.PrivateKey.Sign(nil, hash, &opts)
		if err != nil {
			panic(err)
		}
		return &Signature{
			Address:        key.GetAddress(),
			EdDSASignature: &sig,
			Mode:           key.Mode,
		}
	} else if key.Mode == ModeBLS {
		sig := key.BLS.PrivateKey.SignHash(hash)
		return &Signature{
			Address:        key.GetAddress(),
			BLSSignature: 	sig,
			Mode:           key.Mode,
		}
	} else {
		panic("unsupported mode")
	}

}

// Signs a value
func (key *Key) SignValue(value *Value) error {
	hash := HashValue(key.Mode, *value)

	if value.Signatures == nil {
		value.Signatures = []Signature{}
	}

	sig := key.SignHash(hash)
	value.Signatures = append(value.Signatures, *sig)
	return nil
}

// Signs multiple values
func (key *Key) SignValues(outputs []Value) ([]Value, error) {
	var signed []Value

	for _, i := range outputs {
		key.SignValue(&i)
		signed = append(signed, i)
	}

	return signed, nil
}

// Signs transaction signature request, which is requested by a client
func (key *Key) SignTransactionSigRequest(request *TransactionSigReq) error {
	hash := HashTransactionSigRequest(key.Mode, *request)
	sig := key.SignHash(hash)
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
func Verify(sig *Signature, hash []byte) (bool, error) {
	if sig.Mode == ModeEdDSA {
		eddsaSig := *sig.EdDSASignature
		if len(eddsaSig) != SignatureLength {
			msg := fmt.Sprintf("invalid signature length. wanted: %d, got: %d", SignatureLength, len(eddsaSig))
			return false, errors.New(msg)
		}
		if len(hash) != crypto.SHA512.Size() {
			msg := fmt.Sprintf("invalid hash length. wanted: %d, got: %d", crypto.SHA512.Size(), len(hash))
			return false, errors.New(msg)
		}
		opts := ed25519.Options{
			Hash: crypto.SHA512,
		}
		return ed25519.VerifyWithOptions(sig.Address, hash, eddsaSig, &opts), nil
	} else if sig.Mode == ModeBLS {
		if len(hash) != crypto.SHA3_256.Size() {
			msg := fmt.Sprintf("invalid hash length. wanted: %d, got: %d", crypto.SHA3_256.Size(), len(hash))
			return false, errors.New(msg)
		}
		panic("bls is not implemented yet")
	} else {
		panic("mode not supported")
	}
}

// Performs batch verification (EdDSA mode only)
func VerifyBatch(sigs []Signature, hash []byte) (bool, error) {
	var pks []Address
	var sigsByte [][]byte
	for _, sig := range sigs {
		if sig.Mode != ModeEdDSA {
			return false, errors.New("batch verification is only available for EdDSA, but found other types of signatures")
		}

		pks = append(pks, sig.Address)

		eddsaSig := *sig.EdDSASignature
		if len(eddsaSig) != SignatureLength {
			msg := fmt.Sprintf("invalid signature length. wanted: %d, got: %d", SignatureLength, len(eddsaSig))
			return false, errors.New(msg)
		}
		sigsByte = append(sigsByte, eddsaSig)
	}

	var messages [][]byte
	for range sigs {
		messages = append(messages, hash)
	}

	opts := ed25519.Options{
		Hash: crypto.SHA512,
	}
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
	// check signature type
	var mode int
	if len(value.Signatures) > 0 {
		firstSig := value.Signatures[0]
		mode = firstSig.Mode
	} else {
		return errors.New("got value with no signatures")
	}

	// look out for duplicates signatures
	origins := make(map[[IndexLength]byte]bool)
	for _, sig := range value.Signatures {
		index := [IndexLength]byte{}
		copy(index[:], sig.Address[:])
		if origins[index] {
			return errors.New("duplicate signatures")
		}
		origins[index] = true
	}

	// check whether the signatures have all been made by valid servers
	for _, sig := range value.Signatures {
		if !IsValidServer(sig.Address) {
			return errors.New("encountered signature signed by invalid server")
		}
	}

	// check that there are enough signatures
	numRequiredSigs := QuorumSize()
	if len(value.Signatures) < numRequiredSigs {
		text := fmt.Sprintf("not enough signatures. need %d, have %d", numRequiredSigs, len(value.Signatures))
		return errors.New(text)
	}

	hash := HashValue(mode, *value)
	if mode == ModeEdDSA {
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
	} else {
		panic("verifying bls is not supported yet")
	}
}

// Verifies a signature request
// - checks if all inputs are owned by the same party
// - checks if party signed the request
func VerifyTransactionSigRequest(req *TransactionSigReq) error {
	ownerAddress := req.Signature.Address

	for _, input := range req.Transaction.Inputs {
		if !bytes.Equal(ownerAddress, input.Address) {
			return errors.New("inputs are not owned by the same party")
		}
	}

	hash := HashTransactionSigRequest(req.Signature.Mode, *req)

	valid, err := Verify(&req.Signature, hash)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("sig request verification failed")
	}

	return nil
}
