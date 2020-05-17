package common

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/oasislabs/ed25519"
	"testing"
)

// test basic behavior.
func TestBasic(t *testing.T) {
	id := GenerateKey()

	msg := []byte("hello")
	sig := SignHash(id, msg)

	valid, err := Verify(sig, msg)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("invalid signature")
	}

	msg2 := []byte("world")
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
	id := GenerateKey()
	id2 :=GenerateKey()

	if bytes.Equal(id.Address, id2.Address) || bytes.Equal(id.PrivateKey, id2.PrivateKey) {
		t.Fatal("did not generate different keypairs")
	}
}

func TestKeylength(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	if len(pub) != AddressLength {
		t.Fatal("wrong public key length")
	}

	if len(priv) != PrivateKeyLength {
		t.Fatal("wrong private key length")
	}

	msg := []byte("hello")
	sig := ed25519.Sign(priv, msg)

	// TODO!
	if len(sig) != SignatureLength {
		t.Fatal("wrong private key length")
	}
}

func TestPrintGeneratedKey(t *testing.T) {
	for i := 0; i < 16; i++ {
		id := GenerateKey()
		fmt.Printf("{\"%x\",\"%x\"},\n", id.Address, id.PrivateKey)
	}
}