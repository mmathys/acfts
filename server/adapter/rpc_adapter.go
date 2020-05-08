package adapter

import (
	"errors"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/checks"
	"github.com/mmathys/acfts/server/util"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

var Id *common.Identity
var NoSigning bool
var Benchmark bool
var TxCounter *int32
var SignedUTXO *sync.Map
var AllowDoublespend = false
var UseUTXOMap = true
var CheckTransactions = true

// struct for RPC
type Server struct{}

func (s *Server) Sign(req common.TransactionSigReq, res *common.TransactionSignRes) error {
	if Benchmark {
		defer util.CountTx(TxCounter)
	}

	// Perform checks
	if !NoSigning && CheckTransactions {
		err := checks.CheckValidity(&req)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// Sign transaction
	tx := req.Transaction
	if !NoSigning && UseUTXOMap {
		for _, input := range tx.Inputs {
			_, spent := SignedUTXO.LoadOrStore(input.Id, true)
			if spent && !AllowDoublespend {
				err := errors.New("UTXO already exists: no double spending")
				fmt.Println(err)
				return err
			}
		}
	}

	// Sign the transaction request
	var outputs []common.Value
	if !NoSigning {
		var err error = nil
		outputs, err = common.SignValues(Id.Key, tx.Outputs)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		outputs = tx.Outputs
		for i, _ := range outputs {
			outputs[i].Signatures = [][]byte{}
		}
	}

	// Respond
	*res = common.TransactionSignRes{Outputs: outputs}
	return nil
}

func Init(port int, id *common.Identity, noSigning bool, benchmark bool, txCounter *int32, signedUTXO *sync.Map) {
	Id = id
	NoSigning = noSigning
	Benchmark = benchmark
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
