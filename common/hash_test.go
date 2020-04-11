package common

import (
	"testing"
)

func prepareValue() Value {
	// just fill some data
	addr := [AddressLength]byte{}
	id := [IdentifierLength]byte{}
	idk := [64]byte{}
	sig := ECDSASig{
		Address: addr[:],
		R:       idk[:],
		S:       idk[:],
	}

	return Value{
		Address:    addr[:],
		Amount:     123,
		Id:         id[:],
		Signatures: []ECDSASig{sig, sig, sig, sig},
	}
}

func BenchmarkHashValueSprintf(b *testing.B) {
	value := prepareValue()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashValueSprintf(value)
	}
}

func BenchmarkHashValueGob(b *testing.B) {
	value := prepareValue()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashValue(value)
	}
}