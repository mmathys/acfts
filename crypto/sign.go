package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"golang.org/x/crypto/sha3"
	"math"
)

func HashValue(value common.Value) []byte {
	d := sha3.New256()
	value.Signatures = nil                    // zero out signatures before hash
	d.Write([]byte(fmt.Sprintf("%v", value))) // this may be slow!
	return d.Sum(nil)
}

func SignValue(key *ecdsa.PrivateKey, value *common.Value) error {
	hash := HashValue(*value)

	if value.Signatures == nil {
		value.Signatures = []common.ECDSASig{}
	}

	r, s, err := ecdsa.Sign(rand.Reader, key, hash)
	if err != nil {
		return err
	}

	value.Signatures = append(value.Signatures, common.ECDSASig{R: r, S: s})
	return nil
}

func SignValues(key *ecdsa.PrivateKey, outputs []common.Value) ([]common.Value, error) {
	var signed []common.Value
	for _, i := range outputs {
		SignValue(key, &i)
		signed = append(signed, i)
	}
	return signed, nil
}

/**
Verifies single value
- Verifies all signatures
- Checks whether there are enough signatures to satisfy the validity constraint. (> 2/3 of all sigs)
*/
func VerifyValue(key *ecdsa.PrivateKey, value *common.Value) error {
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

	numServers := core.GetNumServers()
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
func VerifyTransaction(key *ecdsa.PrivateKey, value *common.Transaction) error {
	return nil
}
