package common

import (
	"encoding/hex"
	"errors"
	"strings"
)

func ReadAddress(s string) (Address, error) {
	split := strings.Split(s, "0x")
	if len(split) != 2 {
		return Address{}, errors.New("hex should look like 0x04\n")
	}

	addr, err := hex.DecodeString(split[1])
	if err != nil {
		return Address{}, errors.New("could not decode hex\n")
	}

	return addr, nil
}
