package server

import (
	"github.com/mmathys/acfts/common"
	"log"
	"sync"
)

var _mode = "rest"

func SetAdapterMode(mode string) {
	_mode = mode
}

func Init(port int, id *common.Identity, debug bool, benchmark bool, SignedUTXO *sync.Map, TxCounter *int32) {
	if _mode == "rest" {
		initREST(port, id, debug, benchmark, SignedUTXO, TxCounter)
	} else if _mode == "rpc" {
		initRPC(port, id, debug, benchmark, SignedUTXO, TxCounter)
	} else {
		log.Panic("unrecognized _mode")
	}
}