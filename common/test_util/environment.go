package test_util

import "github.com/mmathys/acfts/common"

func TestEnvironment() {
	common.InitAddresses("../../topologies/localSimple.json")
}

func TestClient(index int) common.Address {
	return common.GetClients()[index]
}
