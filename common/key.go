package common

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/oasislabs/ed25519"
)

func edDSAKey(pub []byte, sk []byte, keyType int) *Key {
	return &Key{
		EdDSA: &EdDSAKey{
			Address:    pub,
			PrivateKey: sk,
		},
		BLS:  nil,
		Type: keyType,
	}
}

func blsKey(pub bls.PublicKey, sk bls.SecretKey, id int) *Key {
	var blsId bls.ID
	blsId.SetLittleEndian([]byte{uint8(id)})
	return &Key{
		EdDSA: nil,
		BLS: &BLSKey{
			ID:         blsId,
			Address:    pub,
			PrivateKey: sk,
		},
		Type: TypeBLS,
	}
}

func GenerateKey(keyType int, id int) *Key {
	if keyType == TypeEdDSA {
		pub, sk, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		return edDSAKey(pub, sk, keyType)
	} else if keyType == TypeBLS {
		var sec bls.SecretKey
		sec.SetByCSPRNG()
		pub := sec.GetPublicKey()
		return blsKey(*pub, sec, id)
	} else {
		panic("invalid type")
	}
}

func DecodeKey(keyType int, index int, pubS string, skS string) (*Key, error) {
	pub, err := hex.DecodeString(pubS)
	if err != nil {
		return nil, err
	}

	sk, err := hex.DecodeString(skS)
	if err != nil {
		return nil, err
	}
	if keyType == TypeEdDSA {
		if len(pub) != EdDSAPublicKeyLength {
			return nil, errors.New("topology: encountered wrong address length")
		}
		if len(sk) != EdDSAPrivateKeyLength {
			return nil, errors.New("topology: encountered wrong private key length")
		}
		return edDSAKey(pub, sk, ModeNaive), nil
	} else if keyType == TypeBLS {
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
		return blsKey(blsPub, blsSk, index+1), nil
	} else {
		panic("unsupported type")
	}
}

func (key *Key) GetAddress() []byte {
	if key.Type == TypeEdDAS || key.Mode == ModeMerkle {
		return key.EdDSA.Address
	} else if key.Mode == ModeBLS {
		return key.BLS.Address.Serialize()
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) GetPrivateKey() []byte {
	if key.Mode == ModeNaive || key.Mode == ModeMerkle {
		return key.EdDSA.PrivateKey
	} else if key.Mode == ModeBLS {
		return key.BLS.PrivateKey.Serialize()
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) SerializePublicKey() []byte {
	if key.Mode == ModeNaive || key.Mode == ModeMerkle {
		return key.EdDSA.Address
	} else if key.Mode == ModeBLS {
		return key.BLS.Address.Serialize()
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) SerializePrivateKey() []byte {
	if key.Mode == ModeNaive || key.Mode == ModeMerkle {
		return key.EdDSA.PrivateKey
	} else if key.Mode == ModeBLS {
		return key.BLS.PrivateKey.Serialize()
	} else {
		panic("unsupported mode")
	}
}
