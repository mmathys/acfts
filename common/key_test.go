package common

import (
	"crypto/rand"
	"encoding/gob"
	"golang.org/x/crypto/sha3"
	"reflect"
	"testing"
)

func TestPubkeyMarshal(t *testing.T) {
	key := GenerateKey()
	encoded := MarshalPubkey(&key.PublicKey)

	msg := "hello"
	d := sha3.New256()
	enc := gob.NewEncoder(d)
	enc.Encode(msg)
	hash := d.Sum(nil)

	sig, err := SignHash(hash, key)
	if err != nil {
		panic(err)
	}

	valid, err := Verify(encoded, hash, sig)
	if err != nil {
		panic(err)
	}
	if !valid {
		t.Error("not valid")
	}
}

func TestRecoverPubkey(t *testing.T) {
	key := GenerateKey()

	msg := "hello"
	d := sha3.New256()
	enc := gob.NewEncoder(d)
	enc.Encode(msg)
	hash := d.Sum(nil)

	sig, err := SignHash(hash, key)
	if err != nil {
		t.Errorf("Sign error: %s", err)
	}

	recoveredPub, err := RecoverPubkeyBytes(hash, sig)
	if err != nil {
		t.Errorf("ECRecover error: %s", err)
	}
	pubKey := UnmarshalPubkey(recoveredPub)

	// should be equal to SigToPub
	recoveredPub2, err := recoverPubkey(hash, sig)
	if err != nil {
		t.Errorf("ECRecover error: %s", err)
	}

	if !reflect.DeepEqual(MarshalPubkey(pubKey), recoveredPub) {
		t.Errorf("pubkey mismatch #0")
	}

	if !reflect.DeepEqual(MarshalPubkey(pubKey), MarshalPubkey(recoveredPub2)) {
		t.Errorf("pubkey mismatch #1")
	}

	if !reflect.DeepEqual(MarshalPubkey(&key.PublicKey), MarshalPubkey(recoveredPub2)) {
		t.Errorf("pubkey mismatch #2")
	}
}


func TestGenerateParseKey(t *testing.T) {
	key := GenerateKey()

	hash := make([]byte, 32) // random hash
	rand.Read(hash)

	sig, err := SignHash(hash, key)
	if err != nil {
		t.Error(err)
	}

	encoded := MarshalKey(key)
	decoded := UnmarshalKey(encoded)
	encodedPub2 := MarshalPubkey(&decoded.PublicKey)

	valid, err := Verify(encodedPub2, hash, sig)
	if err != nil {
		panic(err)
	}
	if !valid {
		t.Error("verification failed")
	}
}
