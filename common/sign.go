package common

import (
	"bytes"
	"crypto"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/oasislabs/ed25519"
	"sync"
)

var merkleSigCache sync.Map

/**
Signing
*/

// Signs a single hash with default mode
func (key *Key) SignHash(hash []byte) *Signature {
	return key.signHashWithMode(hash, key.Mode)
}

// Signs a single hash with a certain mode
func (key *Key) signHashWithMode(hash []byte, mode int) *Signature {
	if mode == ModeEdDSA {
		opts := ed25519.Options{
			Hash: crypto.SHA512,
		}
		sig, err := key.EdDSA.PrivateKey.Sign(nil, hash, &opts)
		if err != nil {
			panic(err)
		}
		return &Signature{
			Address:   key.GetAddress(),
			Signature: sig,
			Mode:      mode,
		}
	} else if mode == ModeBLS {
		sig := key.BLS.PrivateKey.SignHash(hash)
		id := key.BLS.ID
		return &Signature{
			BLSID:     id.Serialize(),
			Address:   key.GetAddress(),
			Signature: sig.Serialize(),
			Mode:      mode,
		}
	} else if mode == ModeMerkle {
		// this is only for a SINGLE merkle signature. should not be used
		fmt.Println("warning: using merkle signing for debugging purposes")
		sigs := key.SignMultipleMerkle([][]byte{hash})
		return sigs[0]
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

	// When signing a transaction sig request, use EdDSA only, even if the default mode is Merkle
	var sig *Signature
	if key.Mode == ModeEdDSA || key.Mode == ModeMerkle {
		sig = key.signHashWithMode(hash, ModeEdDSA)
	} else if key.Mode == ModeBLS {
		sig = key.signHashWithMode(hash, ModeBLS)
	} else {
		panic("unrecognized mode")
	}

	request.Signature = *sig

	return nil
}

// Signs (multiple) signatures efficiently into a Merkle Signature
func (key *Key) SignMultipleMerkle(hashes [][]byte) []*Signature {
	if key.Mode != ModeMerkle {
		panic("when using SignMerkle, the key mode must be merkle")
	}

	t, err := NewTreeWithHashStrategy(hashes, sha512.New)
	if err != nil {
		panic(err)
	}
	root := t.MerkleRoot()
	opts := ed25519.Options{
		Hash: crypto.SHA512,
	}

	rootSig, err := key.EdDSA.PrivateKey.Sign(nil, root, &opts)
	if err != nil {
		panic(err)
	}

	sigs := make([]*Signature, len(hashes))
	for i, item := range hashes {
		path, indexesInt64, err := t.GetMerklePath(item)
		if err != nil {
			panic(err)
		}
		indexes := make([]bool, len(indexesInt64))
		for j, index := range indexesInt64 {
			indexes[j] = index == 1
		}

		sigs[i] = &Signature{
			Address:       key.EdDSA.Address,
			Signature:     rootSig,
			Mode:          ModeMerkle,
			MerklePath:    path,
			MerkleIndexes: indexes,
		}
	}

	return sigs
}

/**
Signature Verification
*/

const (
	// if num sigs >= BatchVerificationThreshold, then batch verification is more efficient.
	BatchVerificationThreshold = 4
)

var UseMerkleSignatureCaching = true

// Verifies a signature. for EdDSA, BLS or Merkle.
func Verify(sig *Signature, hash []byte) (bool, error) {
	if sig.Mode == ModeEdDSA {
		if len(hash) != crypto.SHA512.Size() {
			msg := fmt.Sprintf("invalid hash length. wanted: %d, got: %d", crypto.SHA512.Size(), len(hash))
			return false, errors.New(msg)
		}

		eddsaSig := sig.Signature
		if len(eddsaSig) != SignatureLength {
			msg := fmt.Sprintf("invalid signature length. wanted: %d, got: %d", SignatureLength, len(eddsaSig))
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

		var blsSig bls.Sign
		blsSig.Deserialize(sig.Signature)
		var blsPub bls.PublicKey
		blsPub.Deserialize(sig.Address)

		return blsSig.VerifyHash(&blsPub, hash), nil
	} else if sig.Mode == ModeMerkle {
		// calculate master
		current := hash
		for i := range sig.MerklePath {
			h := sha512.New()
			hash := sig.MerklePath[i]
			index := sig.MerkleIndexes[i]
			var msg []byte
			if index == false {
				// hash is left
				msg = append(hash, current...)
			} else {
				// hash is right
				msg = append(current, hash...)
			}
			if _, err := h.Write(msg); err != nil {
				return false, err
			}

			current = h.Sum(nil)
		}

		// `current` should now be the merkle root.

		// use caching: find out whether we previously already checked that
		// signature is ok. for this, use hash(addr || merkle root || sig)
		h := crypto.SHA256.New()
		h.Write(sig.Address)
		h.Write(current)
		h.Write(sig.Signature)
		sigHash := h.Sum(nil)
		sigHashIndex := [32]byte{}
		copy(sigHashIndex[:], sigHash[:])

		// lookup cache and return if cached
		if UseMerkleSignatureCaching {
			cachedValid, ok := merkleSigCache.Load(sigHashIndex)
			if ok && cachedValid == true {
				return true, nil
			}
		}

		// there is no cache entry, or entry was false.
		opts := ed25519.Options{
			Hash: crypto.SHA512,
		}
		valid := ed25519.VerifyWithOptions(sig.Address, current, sig.Signature, &opts)
		if valid {
			merkleSigCache.Store(sigHashIndex, true)
		}
		return valid, nil
	} else {
		return false, errors.New("mode not supported")
	}
}

// Performs batch verification (for EdDSA mode only)
func VerifyEdDSABatch(sigs []Signature, hash []byte) (bool, error) {
	var pks []ed25519.PublicKey
	var sigsByte [][]byte
	for _, sig := range sigs {
		if sig.Mode != ModeEdDSA {
			return false, errors.New("batch verification is only available for EdDSA, but found other types of signatures")
		}

		pks = append(pks, sig.Address)

		eddsaSig := sig.Signature
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

	hash := HashValue(mode, *value)
	if mode == ModeEdDSA || mode == ModeMerkle {
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

		// check that there are enough signatures
		numRequiredSigs := QuorumSize()
		if len(value.Signatures) < numRequiredSigs {
			text := fmt.Sprintf("not enough signatures. need %d, have %d", numRequiredSigs, len(value.Signatures))
			return errors.New(text)
		}

		// check whether the signatures have all been made by valid servers
		for _, sig := range value.Signatures {
			if !IsValidServer(sig.Address) {
				return errors.New("encountered signature signed by invalid server (eddsa)")
			}
		}

		// verify all signatures, either with batch verification or single verification
		if mode == ModeEdDSA && enableBatchVerification && len(value.Signatures) >= BatchVerificationThreshold {
			// batch verification
			valid, err := VerifyEdDSABatch(value.Signatures, hash)
			if err != nil {
				return err
			}
			if !valid {
				return errors.New("value verification failed (batch mode)")
			}
		} else {
			// single verification for EdDSA and Merkle
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

	} else if mode == ModeBLS {
		// we want exactly one signature
		if len(value.Signatures) != 1 {
			return errors.New("in BLS mode, exactly one combined signature must be sent")
		}

		sig := value.Signatures[0]
		masterPub := GetBLSMasterPublicKey()
		if !bytes.Equal(masterPub.Serialize(), sig.Address) {
			return errors.New("signature address does not match master")
		}

		var masterSig bls.Sign
		masterSig.Deserialize(value.Signatures[0].Signature)
		valid := masterSig.VerifyHash(&masterPub, hash)
		if !valid {
			return errors.New("BLS signature is not valid")
		}
	} else {
		return errors.New("unrecognized mode")
	}

	return nil
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
