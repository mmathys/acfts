package client

import (
	"github.com/mmathys/acfts/common"
	"log"
	"sync"
)

var _mode = "rest"

func SetAdapterMode(mode string) {
	_mode = mode
}

func Init(port int, incoming chan common.Value) {
	if _mode == "rest" {
		initREST(port, incoming)
	} else if _mode == "rpc" {
		initRPC(port, incoming)
	} else {
		log.Panic("unrecognized _mode")
	}
}

func RequestSignature(serverAddr common.Address, id *common.Identity, t common.Transaction, wg *sync.WaitGroup, sigs *chan common.TransactionSignRes) {
	if _mode == "rest" {
		requestSignatureREST(serverAddr, id, t, wg, sigs)
	} else if _mode == "rpc" {
		requestSignatureRPC(serverAddr, id, t, wg, sigs)
	} else {
		log.Panic("unrecognized _mode")
	}
}

func ForwardValue(t common.Value) {
	if _mode == "rest" {
		forwardValueREST(t)
	} else if _mode == "rpc" {
		forwardValueRPC(t)
	} else {
		log.Panic("unrecognized _mode")
	}
}