package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	ethereum "github.com/ethereum/go-ethereum/crypto"
	"github.com/mmathys/acfts/common"
	"testing"
)

/**
Benchmark for crypto ops
*/

func BenchmarkSign(b *testing.B) {
	key := common.GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		ecdsa.Sign(rand.Reader, key, hash)
	}
}

func BenchmarkSignEth(b *testing.B) {
	key := common.GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		ethereum.Sign(hash, key)
	}
}

func BenchmarkVerify(b *testing.B) {
	key := common.GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)
	r, s, _ := ecdsa.Sign(rand.Reader, key, hash)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecdsa.Verify(&key.PublicKey, hash, r, s)
	}
}

func BenchmarkVerifyEth(b *testing.B) {
	key := common.GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)
	sig, err := common.SignHash(hash, key)
	if err != nil {
		b.Fatal(err)
	}
	pubkey, err := common.RecoverPubkeyBytes(hash, sig)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		valid, err := common.Verify(pubkey, hash, sig)
		if err != nil {
			b.Fatal(err)
		}
		if !valid {
			b.Fatal(errors.New("validation failed"))
		}
	}
}
