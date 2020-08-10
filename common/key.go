package common

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/oasislabs/ed25519"
)

func edDSAKey(pub []byte, sk []byte, mode int) *Key {
	return &Key{
		EdDSA: &EdDSAKey{
			Address:    pub,
			PrivateKey: sk,
		},
		BLS:  nil,
		Mode: mode,
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
		Mode: ModeBLS,
	}
}

func GenerateKey(mode int, id int) *Key {
	if mode == ModeEdDSA || mode == ModeMerkle {
		pub, sk, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		return edDSAKey(pub, sk, mode)
	} else if mode == ModeBLS {
		var sec bls.SecretKey
		sec.SetByCSPRNG()
		pub := sec.GetPublicKey()
		return blsKey(*pub, sec, id)
	} else {
		panic("invalid mode")
	}
}

func DecodeKey(mode int, index int, pubS string, skS string) (*Key, error) {
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
		return edDSAKey(pub, sk, ModeEdDSA), nil
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
		return blsKey(blsPub, blsSk, index+1), nil
	} else if mode == ModeMerkle {
		if len(pub) != MerklePublicKeyLength {
			return nil, errors.New("topology: encountered wrong address length")
		}
		if len(sk) != MerklePrivateKeyLength {
			return nil, errors.New("topology: encountered wrong private key length")
		}
		return edDSAKey(pub, sk, ModeMerkle), nil
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) GetAddress() []byte {
	if key.Mode == ModeEdDSA || key.Mode == ModeMerkle {
		return key.EdDSA.Address
	} else if key.Mode == ModeBLS {
		return key.BLS.Address.Serialize()
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) GetPrivateKey() []byte {
	if key.Mode == ModeEdDSA || key.Mode == ModeMerkle {
		return key.EdDSA.PrivateKey
	} else if key.Mode == ModeBLS {
		return key.BLS.PrivateKey.Serialize()
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) SerializePublicKey() []byte {
	if key.Mode == ModeEdDSA || key.Mode == ModeMerkle {
		return key.EdDSA.Address
	} else if key.Mode == ModeBLS {
		return key.BLS.Address.Serialize()
	} else {
		panic("unsupported mode")
	}
}

func (key *Key) SerializePrivateKey() []byte {
	if key.Mode == ModeEdDSA || key.Mode == ModeMerkle {
		return key.EdDSA.PrivateKey
	} else if key.Mode == ModeBLS {
		return key.BLS.PrivateKey.Serialize()
	} else {
		panic("unsupported mode")
	}
}
