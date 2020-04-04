package common

import (
	"crypto/ecdsa"
	"crypto/rand"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"log"
)



func GenerateKey() *ecdsa.PrivateKey {
	key, _ := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	return key
}

func MarshalPubkey(pub *ecdsa.PublicKey) Address {
	encoded := crypto2.FromECDSAPub(pub)

	if len(encoded) != AddressLength {
		log.Fatalln("key length does not match when marshalling")
	}

	return encoded
}

func UnmarshalPubkey(pub Address) *ecdsa.PublicKey {
	decoded, err := crypto2.UnmarshalPubkey(pub[:])
	if err != nil {
		log.Fatalln("could not unmarshal pubkey")
	}
	return decoded
}

func MarshalKey(key *ecdsa.PrivateKey) *PrivateKey {
	encoded := crypto2.FromECDSA(key)

	if len(encoded) != PrivateKeyLength {
		log.Fatalln("key length does not match when marshalling private key")
	}

	return &encoded
}

func UnmarshalPrivateKey(key *PrivateKey) *ecdsa.PrivateKey {
	res, err := crypto2.ToECDSA(*key)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func UnmarshalPrivateKeyHex(key string) *ecdsa.PrivateKey {
	res, err := crypto2.HexToECDSA(key)
	if err != nil {
		log.Fatal(err)
	}
	return res
}