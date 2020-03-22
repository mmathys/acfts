package server

import (
	"github.com/mmathys/acfts/common"
	"testing"
)

func TestBasic(t *testing.T) {
	server := common.GetServers()[0]
	RunServer(server, true)
}