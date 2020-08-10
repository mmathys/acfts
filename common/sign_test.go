package common

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/oasislabs/ed25519"
	"testing"
)

// test basic behavior.
func TestBasic(t *testing.T) {
	key := GenerateKey(ModeEdDSA, 0)

	msg := make([]byte, 64) // random hash
	rand.Read(msg)
	sig := key.SignHash(msg)

	valid, err := Verify(sig, msg)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("invalid signature")
	}

	msg2 := make([]byte, 64) // random hash
	rand.Read(msg2)
	valid2, err := Verify(sig, msg2)
	if err != nil {
		t.Fatal(err)
	}
	if valid2 {
		t.Fatal("signature should be invalid")
	}

	sig.Address[0]++
	valid3, err := Verify(sig, msg)
	if err != nil {
		t.Fatal(err)
	}
	if valid3 {
		t.Fatal("Signature should be invalid")
	}
}

func TestGenerateKeypair(t *testing.T) {
	key := GenerateKey(ModeEdDSA, 0)
	key2 := GenerateKey(ModeEdDSA, 0)

	if bytes.Equal(key.GetAddress(), key2.GetAddress()) || bytes.Equal(key.GetPrivateKey(), key2.GetPrivateKey()) {
		t.Fatal("did not generate different keypairs")
	}
}

func TestKeylength(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	if len(pub) != EdDSAPublicKeyLength {
		t.Fatal("wrong public key length")
	}

	if len(priv) != EdDSAPrivateKeyLength {
		t.Fatal("wrong private key length")
	}

	msg := make([]byte, 64) // random hash
	rand.Read(msg)
	opts := ed25519.Options{
		Hash: crypto.SHA512,
	}
	sig, err := priv.Sign(nil, msg, &opts)
	if err != nil {
		panic(err)
	}

	// TODO!
	if len(sig) != SignatureLength {
		t.Fatal("wrong private key length")
	}
}

func TestPrintGeneratedKey(t *testing.T) {
	mode := ModeBLS
	if mode == ModeBLS {
		bls.Init(bls.BLS12_381)
		bls.SetETHmode(bls.EthModeDraft07)
	}
	for i := 0; i < 64; i++ {
		key := GenerateKey(mode, i)
		fmt.Printf("{\"%x\",\"%x\"},\n", key.SerializePublicKey(), key.SerializePrivateKey())
	}
}
