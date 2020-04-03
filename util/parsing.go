package util

import (
	"encoding/hex"
	"errors"
	"github.com/mmathys/acfts/common"
	"strings"
)

func ReadAddress(s string) (common.Address, error) {
	split := strings.Split(s, "0x")
	if len(split) != 2 {
		return common.Address{}, errors.New("hex should look like 0x04\n")
	}

	addr, err := hex.DecodeString(split[1])
	if err != nil {
		return common.Address{}, errors.New("could not decode hex\n")
	}

	return addr, nil
}