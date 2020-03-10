package core

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/mmathys/acfts/common"
	"golang.org/x/crypto/sha3"
)

func doHash(value common.Value) []byte {
	d := sha3.New256()
	d.Write([]byte(fmt.Sprintf("%v", value))) // this may be slow!
	return d.Sum(nil)
}

func Sign(key *ecdsa.PrivateKey, outputs []common.Value) ([]common.Value, error) {
	var signed []common.Value

	for _, i := range outputs {
		if i.Signatures == nil {
			i.Signatures = []common.ECDSASig{}
		}

		hash := doHash(i)
		r, s, err := ecdsa.Sign(rand.Reader, key, hash)
		if err != nil {
			return nil, err
		}

		i.Signatures = append(i.Signatures, common.ECDSASig{R:r, S:s})

		signed = append(signed, i)
	}



	return signed, nil
}
