package common

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/oasislabs/ed25519"
)

func edDSAKey(pub []byte, sk []byte) *Key {
	return &Key{
		EdDSA: &EdDSAKey{
			Address:    pub,
			PrivateKey: sk,
		},
		BLS:   nil,
		Mode:  ModeEdDSA,
	}
}

func blsKey(pub []byte, sk []byte) *Key {
	panic("not implemented yet")
}

func GenerateKey(mode int) *Key {
	if mode == ModeEdDSA {
		pub, sk, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		return edDSAKey(pub, sk)
	} else if mode == ModeBLS {
		panic("bls not yet supported")
	} else {
		panic("invalid mode")
	}
}

func DecodeKey(mode int, pubS string, skS string) (*Key, error) {
	if mode == ModeEdDSA {
		pub, err := hex.DecodeString(pubS)
		if err != nil {
			return nil, err
		}

		sk, err := hex.DecodeString(skS)
		if err != nil {
			return nil, err
		}

		if len(pub) != EdDSAPublicKeyLength {
			return nil, errors.New("topology: encountered wrong address length")
		}
		if len(sk) != EdDSAPrivateKeyLength {
			return nil, errors.New("topology: encountered wrong private key length")
		}

		return edDSAKey(pub, sk), nil
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) GetAddress() []byte {
	if key.Mode == ModeEdDSA {
		return key.EdDSA.Address
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) GetPrivateKey() []byte {
	if key.Mode == ModeEdDSA {
		return key.EdDSA.PrivateKey
	} else {
		panic("unsupported mode")
	}
}