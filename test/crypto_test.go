package test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"testing"
)

func generateKey() *ecdsa.PrivateKey {
	key, _ := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	return key
}

func TestSignVerify(t *testing.T) {

	key := generateKey()
	privkey := crypto.FromECDSA(key)
	pubkey := crypto.FromECDSAPub(&key.PublicKey)

	msg := make([]byte, 32) // 32 bytes of random data for example, SHOULD BE HASHED.
	//rand.Read(msg)

	sig, err := secp256k1.Sign(msg, privkey)
	if err != nil {
		t.Error(err)
	}

	valid := secp256k1.VerifySignature(pubkey, msg, sig)

	if !valid {
		t.Error("verification failed")
	}
}
