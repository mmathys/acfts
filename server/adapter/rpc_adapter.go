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

//var SignedUTXO *funset.FunSet
var SignedUTXO *sync.Map
var AllowDoublespend = false
var UseUTXOMap = true
var CheckTransactions = true
var BatchVerification = true

// struct for RPC
type Server struct{}

// Signs a Transaction Request
func (s *Server) Sign(req common.TransactionSigReq, res *common.TransactionSignRes) error {
	if Benchmark {
		defer util.CountTx(TxCounter)
	}

	// Perform checks
	if !NoSigning && CheckTransactions {
		err := checks.CheckValidity(&req, BatchVerification)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// Sign transaction
	tx := req.Transaction
	if !NoSigning && UseUTXOMap {
		for _, input := range tx.Inputs {
			//inserted := SignedUTXO.Insert(input.Id)
			_, spent := SignedUTXO.LoadOrStore(input.Id, true)
			//if !inserted && !AllowDoublespend {
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
		outputs, err = common.SignValues(Id, tx.Outputs)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		outputs = tx.Outputs
		for i, _ := range outputs {
			outputs[i].Signatures = []common.EdDSASig{}
		}
	}

	// Respond
	*res = common.TransactionSignRes{Outputs: outputs}
	return nil
}

// Initialises the adapter
//func Init(port int, id *common.Identity, noSigning bool, benchmark bool, txCounter *int32, signedUTXO *funset.FunSet) {
func Init(port int, id *common.Identity, noSigning bool, benchmark bool, txCounter *int32, signedUTXO *sync.Map, batchVerification bool) {
	Id = id
	NoSigning = noSigning
	Benchmark = benchmark
	TxCounter = txCounter
	SignedUTXO = signedUTXO
	BatchVerification = batchVerification

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
