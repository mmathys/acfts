package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestPrintGeneratedKey(t *testing.T) {
	key := GenerateKey()
	priv := crypto.FromECDSA(key)
	pub := crypto.FromECDSAPub(&key.PublicKey)
	fmt.Printf("%x\n%x\n", priv, pub)
}

func TestGenerateParseKey(t *testing.T) {
	key := GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)

	r, s, err := ecdsa.Sign(rand.Reader, key, hash)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%x%x\n", r, s)

	encoded := crypto.FromECDSA(key)
	decoded, err := crypto.ToECDSA(encoded)
	if err != nil {
		t.Error(err)
	}

	valid := ecdsa.Verify(&decoded.PublicKey, hash, r, s)

	if !valid {
		t.Error("verification failed")
	}
}

func TestGenerateParsePubkey(t *testing.T) {
	key := GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)

	r, s, err := ecdsa.Sign(rand.Reader, key, hash)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%x%x\n", r, s)

	encoded := crypto.FromECDSAPub(&key.PublicKey)
	decoded, err := crypto.UnmarshalPubkey(encoded)
	if err != nil {
		t.Error(err)
	}

	valid := ecdsa.Verify(decoded, hash, r, s)

	if !valid {
		t.Error("verification failed")
	}
}


func TestSignVerify(t *testing.T) {
	key := GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)

	r, s, err := ecdsa.Sign(rand.Reader, key, hash)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%x%x\n", r, s)

	valid := ecdsa.Verify(&key.PublicKey, hash, r, s)

	if !valid {
		t.Error("verification failed")
	}
}
