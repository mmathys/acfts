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


var Id *common.Identity
var Debug bool
var BenchmarkMode bool
var TxCounter *int32
var SignedUTXO *hashmap.HashMap
var AllowDoublespend = false

type Server struct {}
type RPCAdapter struct {}

func (s *Server) Sign(req common.TransactionSigReq, res *common.TransactionSignRes) error {
	//log.Printf("got sign request: %v", req)

	if BenchmarkMode {
		defer util.CountTx(TxCounter)
	}

	if !Debug {
		err := checks.CheckValidity(&req)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	tx := req.Transaction
	if !Debug {
		for _, input := range tx.Inputs {
			notSpent := SignedUTXO.Insert(input.Id, true)
			if !notSpent && !AllowDoublespend {
				err := errors.New("UTXO already exists: no double spending")
				fmt.Println(err)
				return err
			}
		}
	}

	// Sign the request
	var outputs []common.Value
	if Debug {
		outputs = tx.Outputs
		for i, _ := range outputs {
			outputs[i].Signatures = []common.ECDSASig{}
		}
	} else {
		var err error = nil
		outputs, err = common.SignValues(Id.Key, tx.Outputs)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// Form the response
	*res = common.TransactionSignRes{Outputs: outputs}
	return nil
}

func (a *RPCAdapter) Init(port int, _id *common.Identity, debug bool, benchmark bool, txCounter *int32, signedUTXO *hashmap.HashMap) {
	Id = _id
	Debug = debug
	BenchmarkMode = benchmark
	TxCounter = txCounter
	SignedUTXO = signedUTXO

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
