package common

import (
	"encoding/gob"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"reflect"
	"testing"
)

func TestPubkeyMarshal(t *testing.T) {
	key := GenerateKey()
	encoded := crypto2.FromECDSAPub(&key.PublicKey)

	msg := "hello"
	d := sha3.New256()
	enc := gob.NewEncoder(d)
	enc.Encode(msg)
	hash := d.Sum(nil)

	sig, err := crypto2.Sign(hash, key)
	if err != nil {
		panic(err)
	}

	sig2 := sig[:len(sig)-1]
	valid := crypto2.VerifySignature(encoded, hash, sig2)
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

	sig, err := crypto2.Sign(hash, key)
	if err != nil {
		t.Errorf("Sign error: %s", err)
	}

	recoveredPub, err := crypto2.Ecrecover(hash, sig)
	if err != nil {
		t.Errorf("ECRecover error: %s", err)
	}
	pubKey, _ := crypto2.UnmarshalPubkey(recoveredPub)

	// should be equal to SigToPub
	recoveredPub2, err := crypto2.SigToPub(hash, sig)
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
