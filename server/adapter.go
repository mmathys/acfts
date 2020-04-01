package server

import (
	"github.com/cornelk/hashmap"
	"github.com/mmathys/acfts/common"
	"log"
)

type Adapter interface {
	Init(port int, id *common.Identity, debug bool, benchmark bool, TxCounter *int32)
}

var restAdapter = &RESTAdapter{}
var rpcAdapter = &RPCAdapter{}
var currentAdapter Adapter = restAdapter

var SignedUTXO hashmap.HashMap

func SetAdapterMode(mode string) {
	if mode == "rest" {
		currentAdapter = restAdapter
	} else if mode == "rpc" {
		currentAdapter = rpcAdapter
	} else {
		log.Fatalf("unrecognized mode %s", mode)
	}
}

func Init(port int, id *common.Identity, debug bool, benchmark bool, TxCounter *int32) {
	currentAdapter.Init(port, id, debug, benchmark, TxCounter)
}