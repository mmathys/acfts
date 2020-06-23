package common

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/oasislabs/ed25519"
)

func edDSAKey(pub []byte, sk []byte) *Key {
	return &Key{
		EdDSA: &EdDSAKey{
			Address:    pub,
			PrivateKey: sk,
		},
		BLS:  nil,
		Mode: ModeEdDSA,
	}
}

func blsKey(pub bls.PublicKey, sk bls.SecretKey) *Key {
	return &Key{
		EdDSA: nil,
		BLS: &BLSKey{
			Address:    pub,
			PrivateKey: sk,
		},
		Mode: ModeBLS,
	}
}

func GenerateKey(mode int) *Key {
	if mode == ModeEdDSA {
		pub, sk, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		return edDSAKey(pub, sk)
	} else if mode == ModeBLS {
		var sec bls.SecretKey
		sec.SetByCSPRNG()
		pub := sec.GetPublicKey()
		return blsKey(*pub, sec)
	} else {
		panic("invalid mode")
	}
}

func DecodeKey(mode int, pubS string, skS string) (*Key, error) {
	pub, err := hex.DecodeString(pubS)
	if err != nil {
		return nil, err
	}

	sk, err := hex.DecodeString(skS)
	if err != nil {
		return nil, err
	}
	if mode == ModeEdDSA {
		if len(pub) != EdDSAPublicKeyLength {
			return nil, errors.New("topology: encountered wrong address length")
		}
		if len(sk) != EdDSAPrivateKeyLength {
			return nil, errors.New("topology: encountered wrong private key length")
		}
		return edDSAKey(pub, sk), nil
	} else if mode == ModeBLS {
		if len(pub) != BLSPublicKeyLength {
			return nil, errors.New("topology: encountered wrong address length")
		}
		if len(sk) != BLSPrivateKeyLength {
			return nil, errors.New("topology: encountered wrong private key length")
		}
		var blsPub bls.PublicKey
		blsPub.Deserialize(pub)
		var blsSk bls.SecretKey
		blsSk.Deserialize(sk)
		return blsKey(blsPub, blsSk), nil
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) GetAddress() []byte {
	if key.Mode == ModeEdDSA {
		return key.EdDSA.Address
	} else if key.Mode == ModeBLS {
		return key.BLS.Address.Serialize()
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) SerializePublicKey() []byte {
	if key.Mode == ModeEdDSA {
		panic("not yet implemented")
	} else if key.Mode == ModeBLS {
		return key.BLS.Address.Serialize()
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) SerializePrivateKey() []byte {
	if key.Mode == ModeEdDSA {
		return key.EdDSA.PrivateKey
	} else if key.Mode == ModeBLS {
		return key.BLS.PrivateKey.Serialize()
	} else {
		panic("unsupported mode")
	}
}
