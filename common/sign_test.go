package common

import (
	"fmt"
	"testing"
)

func TestPrintGeneratedKey(t *testing.T) {
	for i := 0; i < 10; i++ {
		key := GenerateKey()
		priv := MarshalKey(key)
		pub := MarshalPubkey(&key.PublicKey)
		fmt.Printf("\"%x\",\"%x\"\n", *priv, pub)
	}
}