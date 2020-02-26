package core

import (
	"errors"
	"github.com/mmathys/acfts/common"
	"reflect"
)

type Entry struct {
	address common.Address
	network string
}

var m = map[string] Entry {
	"A": {common.Address{0}, "http://localhost:5555"},
	//"B": {common.Address{1}, "http://localhost:5556"},
	//"C": {common.Address{2}, "http://localhost:5557"},
	"W": {common.Address{3}, "http://localhost:6666"},
	//"X": {common.Address{4}, "http://localhost:6667"},
	//"Y": {common.Address{5}, "http://localhost:6668"},
	//"Z": {common.Address{6}, "http://localhost:6669"},
}

func LookupAddress(alias string) common.Address {
	return m[alias].address
}

func LookupNetwork(alias string) string {
	return m[alias].network
}

func LookupNetworkFromAddress(address common.Address) (string, error) {
	for _, e := range m {
		if reflect.DeepEqual(e.address, address) {
			return e.network, nil
		}
	}

	return "", errors.New("could not find address")
}

func GetServers() []common.Address {
	return []common.Address {
		common.Address{3},
		//common.Address{4},
		//common.Address{5},
		//common.Address{6},
	}
}

func GetOwnAddress() common.Address {
	return common.Address{0}
}