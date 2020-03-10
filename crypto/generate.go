package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func GenerateKey() *ecdsa.PrivateKey {
	key, _ := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	return key
}
