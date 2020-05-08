package environment

import "github.com/mmathys/acfts/common"

func TestClient(index int) common.Address {
	return common.GetClients()[index]
}
