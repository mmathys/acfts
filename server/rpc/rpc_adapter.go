package rpc

import (
	"errors"
	"fmt"
	"github.com/cornelk/hashmap"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/checks"
	"github.com/mmathys/acfts/util"
	"log"
	"net"
	"net/http"
	"net/rpc"
)


var id *common.Identity
var debug bool
var benchmarkMode bool
var TxCounter *int32
var SignedUTXO *hashmap.HashMap

type Server struct {}
type RPCAdapter struct {}

func (s *Server) Sign(req common.TransactionSigReq, res *common.TransactionSignRes) error {
	//log.Printf("got sign request: %v", req)

	if benchmarkMode {
		defer util.CountTx(TxCounter)
	}

	if !debug {
		err := checks.CheckValidity(&req)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	tx := req.Transaction
	if !debug {
		for _, input := range tx.Inputs {
			found := SignedUTXO.Cas(input.Id, true, true)
			if found {
				err := errors.New("UTXO already exists: no double spending")
				fmt.Println(err)
				return err
			}
		}
	}

	// Sign the request
	var outputs []common.Value
	if debug {
		outputs = tx.Outputs
		for i, _ := range outputs {
			outputs[i].Signatures = []common.ECDSASig{}
		}
	} else {
		var err error = nil
		outputs, err = common.SignValues(id.Key, tx.Outputs)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// Form the response
	*res = common.TransactionSignRes{Outputs: outputs}
	return nil
}

func (a *RPCAdapter) Init(port int, _id *common.Identity, _debug bool, _benchmark bool, _TxCounter *int32, _SignedUTXO *hashmap.HashMap) {
	id = _id
	debug = _debug
	benchmarkMode  = _benchmark
	TxCounter = _TxCounter
	SignedUTXO = _SignedUTXO

	addr := fmt.Sprintf(":%d", port)
	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	http.Serve(lis, nil)
}
