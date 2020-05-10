package common

import (
	"encoding/hex"
	"errors"
)

// reads a hex address
func ReadAddress(s string) (Address, error) {
	addr, err := hex.DecodeString(s)
	if err != nil {
		return Address{}, errors.New("could not decode hex\n")
	}

	return addr, nil
}
