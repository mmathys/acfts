package common

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	secp256k1 "github.com/ethereum/go-ethereum/crypto"
	"math"
)

/**
Signature Recovery
 */
// Recovers a ECDSA public key (bytes, uncompressed) from a hash and signature. Using secp256k1 C-bindings crypto.
func RecoverPubkeyBytes(hash []byte, sig []byte) ([]byte, error) {
	return secp256k1.Ecrecover(hash, sig)
}

// Recovers a ECDSA public key (*ecdsa.PublicKey) from a hash and signature. Using secp256k1 C-bindings crypto.
func recoverPubkey(hash []byte, sig []byte) (*ecdsa.PublicKey, error) {
	return secp256k1.SigToPub(hash, sig)
}

// Recovers a an address from a hash and signature. Using secp256k1 C-bindings crypto.
func recoverAddress(hash []byte, sig []byte) (Address, error) {
	owner, err := recoverPubkey(hash, sig)
	if err != nil {
		return nil, err
	}
	return MarshalPubkey(owner), nil
}

/**
Signing
 */

// Signs a hash
func SignHash(hash []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := secp256k1.Sign(hash, key)
	if err != nil {
		return []byte{}, err
	}
	return sig, nil
}

// Signs a value
func SignValue(key *ecdsa.PrivateKey, value *Value) error {
	hash := HashValue(*value)

	if value.Signatures == nil {
		value.Signatures = []ECDSASig{}
	}

	addr := MarshalPubkey(&key.PublicKey)

	sig, err := SignHash(hash, key)
	if err != nil {
		return err
	}

	value.Signatures = append(value.Signatures, ECDSASig{
		Address: addr,
		RS:      sig,
	})
	return nil
}

// Signs multiple values
func SignValues(key *ecdsa.PrivateKey, outputs []Value) ([]Value, error) {
	var signed []Value

	for _, i := range outputs {
		SignValue(key, &i)
		signed = append(signed, i)
	}

	return signed, nil
}

// Signs transaction signature request, which is requested by a client
func SignTransactionSigRequest(key *ecdsa.PrivateKey, request *TransactionSigReq) error {
	addr := MarshalPubkey(&key.PublicKey)
	hash := HashTransactionSigRequest(*request)
	sig, err := SignHash(hash, key)
	if err != nil {
		return err
	}
	request.Signature = ECDSASig{
		Address: addr,
		RS:      sig,
	}

	return nil
}

/**
Signature Verification
*/

// Verifies a signature. Using secp256k1 C bindings crypto
func Verify(pubkey []byte, hash []byte, sig []byte) (bool, error) {
	if len(sig) != secp256k1.SignatureLength { // 64 + 1
		msg := fmt.Sprintf("invalid signature length. wanted: %d, got: %d", secp256k1.SignatureLength, len(sig))
		return false, errors.New(msg)
	}
	sig = sig[:len(sig)-1] // remove recovery bit
	return secp256k1.VerifySignature(pubkey, hash, sig), nil
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
		valid, err := Verify(sig.Address, hash, sig.RS)
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

	valid, err := Verify(ownerAddress, hash, req.Signature.RS)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("sig request verification failed")
	}

	return nil
}
