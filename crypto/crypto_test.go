package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"testing"
)



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
