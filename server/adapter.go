package server

import (
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/rest"
	"github.com/mmathys/acfts/server/rpc"
	"log"
	"sync"
)

type Adapter interface {
	Init(port int, id *common.Identity, debug bool, benchmark bool, SignedUTXO *sync.Map, TxCounter *int32)
}

var restAdapter = &rest.Adapter{}
var rpcAdapter = &rpc.Adapter{}
var currentAdapter Adapter = restAdapter

func SetAdapterMode(mode string) {
	if mode == "rest" {
		currentAdapter = restAdapter
	} else if mode == "rpc" {
		currentAdapter = rpcAdapter
	} else {
		log.Fatalf("unrecognized mode %s", mode)
	}
}

func Init(port int, id *common.Identity, debug bool, benchmark bool, SignedUTXO *sync.Map, TxCounter *int32) {
	currentAdapter.Init(port, id, debug, benchmark, SignedUTXO, TxCounter)
}