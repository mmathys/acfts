package core

import (
	"errors"
	"fmt"
	"github.com/mmathys/acfts/common"
)

type Entry struct {
	address common.Address
	network string
}

var m = map[common.Address]string{
	common.Address{0}: "http://localhost:5555", // A (client)
	common.Address{1}: "http://localhost:5556", // B (client)
	common.Address{2}: "http://localhost:5557", // C (client)
	common.Address{3}: "http://localhost:6666", // W (server)
	common.Address{4}: "http://localhost:6667", // X (server)
	common.Address{5}: "http://localhost:6668", // Y (server)
	common.Address{6}: "http://localhost:6669", // Z (server)
}

func LookupNetworkFromAddress(address common.Address) (string, error) {
	res, ok := m[address]
	if ok {
		return res, nil
	} else {
		msg := fmt.Sprintf("could not find address 0x%x\n", address)
		return "", errors.New(msg)
	}
}

func GetServers() []common.Address {
	return []common.Address{
		common.Address{3},
		//common.Address{4},
		//common.Address{5},
		//common.Address{6},
	}
}
