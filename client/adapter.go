package main

import (
	"github.com/mmathys/acfts/client/rest"
	"github.com/mmathys/acfts/client/rpc"
	"github.com/mmathys/acfts/common"
	"log"
	"sync"
)

type Adapter interface {
	Init(port int, incoming chan common.Value)
	RequestSignature(serverAddr common.Address, id *common.Identity, t common.Transaction, wg *sync.WaitGroup, sigs *chan common.TransactionSignRes)
	ForwardValue(t common.Value)
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
		log.Panic("unrecognized mode")
	}
}

func Init(port int, incoming chan common.Value) {
	currentAdapter.Init(port, incoming)
}

func RequestSignature(serverAddr common.Address, id *common.Identity, t common.Transaction, wg *sync.WaitGroup, sigs *chan common.TransactionSignRes) {
	currentAdapter.RequestSignature(serverAddr, id, t, wg, sigs)
}

func ForwardValue(t common.Value) {
	currentAdapter.ForwardValue(t)
}
