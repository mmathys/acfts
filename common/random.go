package common

import (
	"crypto/rand"
	"encoding/binary"
)

// from https://stackoverflow.com/questions/35203635/golang-cryptographic-shuffle and modified

type CryptoRandSource struct{}

func NewCryptoRandSource() CryptoRandSource {
	return CryptoRandSource{}
}

func (_ CryptoRandSource) Int63() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	// mask off sign bit to ensure positive number
	return int64(binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1))
}

func (_ CryptoRandSource) Seed(_ int64) {}
