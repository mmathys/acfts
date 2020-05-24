package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	secp256k1 "github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func BenchmarkSignSecp256k1(b *testing.B) {
	key, err := secp256k1.GenerateKey()
	if err != nil {
		panic(err)
	}

	hash := make([]byte, 32) // random hash
	rand.Read(hash)

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		secp256k1.Sign(hash, key)
	}
}

func BenchmarkVerifySecp256k1NoRecovery(b *testing.B) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	hash := make([]byte, 32) // random hash
	rand.Read(hash)
	sig, err := secp256k1.Sign(hash, key)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	pubkey, err := secp256k1.Ecrecover(hash, sig)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		valid := secp256k1.VerifySignature(pubkey, hash, sig[:64])
		if !valid {
			b.Fatal(errors.New("validation failed"))
		}
	}
}

func BenchmarkVerifySecp256k1Recovery(b *testing.B) {
	key, err := secp256k1.GenerateKey()
	if err != nil {
		panic(err)
	}

	hash := make([]byte, 32) // random hash
	rand.Read(hash)
	sig, err := secp256k1.Sign(hash, key)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pubkey, err := secp256k1.Ecrecover(hash, sig)
		if err != nil {
			b.Fatal(err)
		}
		valid := secp256k1.VerifySignature(pubkey, hash, sig[:64])
		if !valid {
			b.Fatal(errors.New("validation failed"))
		}
	}
}

func BenchmarkRecoverSecp256k1(b *testing.B) {
	key, err := secp256k1.GenerateKey()
	if err != nil {
		panic(err)
	}

	hash := make([]byte, 32) // random hash
	rand.Read(hash)
	sig, err := secp256k1.Sign(hash, key)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := secp256k1.Ecrecover(hash, sig)
		if err != nil {
			b.Fatal(err)
		}
	}
}
