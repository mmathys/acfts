package common

import (
	"crypto"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"golang.org/x/crypto/sha3"
	"testing"
)

func TestSha512(t *testing.T) {
	msg := "hello"
	// crypto/sha512
	h := crypto.SHA512.New()
	enc1 := gob.NewEncoder(h)
	enc1.Encode(msg)
	hash1 := h.Sum(nil)
	fmt.Println(len(hash1))

	// sha3-512
	d := sha3.New512()
	enc2 := gob.NewEncoder(d)
	enc2.Encode(msg)
	hash2 := d.Sum(nil)
	fmt.Println(len(hash2))
}

func BenchmarkCryptoSha256(b *testing.B) {
	msg := make([]byte, 64) // random hash
	rand.Read(msg)

	b.ResetTimer()
	// crypto/sha256
	for i := 0; i < b.N; i++ {
		h := crypto.SHA256.New()
		h.Write(msg)
		h.Sum(nil)
	}
}

func BenchmarkCryptoSha512(b *testing.B) {
	msg := make([]byte, 64) // random hash
	rand.Read(msg)

	b.ResetTimer()
	// crypto/sha512
	for i := 0; i < b.N; i++ {
		h := crypto.SHA512.New()
		h.Write(msg)
		h.Sum(nil)
	}
}

func BenchmarkSha3256(b *testing.B) {
	msg := make([]byte, 64) // random hash
	rand.Read(msg)

	b.ResetTimer()
	// sha3-512
	for i := 0; i < b.N; i++ {
		d := sha3.New256()
		d.Write(msg)
		d.Sum(nil)
	}
}

func BenchmarkSha3512(b *testing.B) {
	msg := make([]byte, 64) // random hash
	rand.Read(msg)

	b.ResetTimer()
	// sha3-512
	for i := 0; i < b.N; i++ {
		d := sha3.New512()
		d.Write(msg)
		d.Sum(nil)
	}
}
