package rpc

import (
	"errors"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/checks"
	"github.com/mmathys/acfts/util"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)


var Id *common.Identity
var Debug bool
var BenchmarkMode bool
var TxCounter *int32
var SignedUTXO *sync.Map
var AllowDoublespend = false
var UseUTXOMap = true
var CheckTransactions = true

type Server struct {}
type RPCAdapter struct {}

func (s *Server) Sign(req common.TransactionSigReq, res *common.TransactionSignRes) error {
	if BenchmarkMode {
		defer util.CountTx(TxCounter)
	}

	if !Debug && CheckTransactions {
		err := checks.CheckValidity(&req)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	tx := req.Transaction
	if !Debug && UseUTXOMap {
		for _, input := range tx.Inputs {
			_ , spent := SignedUTXO.LoadOrStore(input.Id, true)
			if spent && !AllowDoublespend {
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

func (a *RPCAdapter) Init(port int, _id *common.Identity, debug bool, benchmark bool, txCounter *int32, signedUTXO *sync.Map) {
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
