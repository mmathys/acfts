package adapter

import (
	"errors"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/checks"
	"github.com/mmathys/acfts/server/merkle"
	"github.com/mmathys/acfts/server/store"
	"github.com/mmathys/acfts/server/util"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"runtime"
	"sync"
)

var Key *common.Key
var NoSigning bool
var Benchmark bool
var TxCounter *int32

var UTXOMap *store.UTXOMap
var AllowDoublespend = false
var UseUTXOMap = true
var CheckTransactions = true
var BatchVerification = true
var MerklePooling = false
var MerkleRequests = make(chan *merkle.PoolMsg)
var MerkleDispatches = make(chan []*merkle.PoolMsg)

// struct for RPC
type Server struct{}
type AdapterOpt struct {
	Port              int
	Key               *common.Key
	NoSigning         bool
	Benchmark         bool
	TxCounter         *int32
	UTXOMap           *store.UTXOMap
	BatchVerification bool
	MerklePooling     bool
	MerklePoolSize    int
}

// Signs UTXOs
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

	// Check for UTXO Map entries
	tx := req.Transaction
	if !NoSigning && UseUTXOMap {
		for _, input := range tx.Inputs {
			spent := UTXOMap.Store(input.Id)
			if spent && !AllowDoublespend {
				err := errors.New("UTXO already exists: no double spending")
				fmt.Println(err)
				return err
			}
		}
	}

	var outputs []common.Value

	if !MerklePooling && !NoSigning {
		// Non-merkle: Immediately sign the transaction request
		var err error = nil
		outputs, err = Key.SignValues(tx.Outputs)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Respond
		*res = common.TransactionSignRes{Outputs: outputs}
	} else if MerklePooling && !NoSigning {
		// Merkle: Pool the transaction request and wait
		var wg sync.WaitGroup
		wg.Add(1)
		MerkleRequests <- &merkle.PoolMsg{
			Req:       &req,
			Res:       res,
			WaitGroup: &wg,
		}
		wg.Wait()
	} else {
		// No signing at all: for debugging purposes only
		outputs = tx.Outputs
		for i, _ := range outputs {
			outputs[i].Signatures = []common.Signature{}
		}

		// Respond
		*res = common.TransactionSignRes{Outputs: outputs}
	}

	return nil
}

// Initialises the adapter
func Init(opt AdapterOpt) {
	Key = opt.Key
	NoSigning = opt.NoSigning
	Benchmark = opt.Benchmark
	TxCounter = opt.TxCounter
	UTXOMap = opt.UTXOMap
	BatchVerification = opt.BatchVerification
	MerklePooling = opt.MerklePooling

	UTXOMap.Init()

	if MerklePooling {
		initMerklePooling()
	}

	addr := fmt.Sprintf(":%d", opt.Port)
	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	http.Serve(lis, nil)
}

// Merkle pooling
func initMerklePooling() {
	// initialize a single collector with threshold 2 (TODO)
	go merkle.CollectAndDispatch(2, MerkleRequests, MerkleDispatches)

	// initialize runtime.NumCPU() many processors
	for i := 0; i < runtime.NumCPU(); i++ {
		go merkle.Processor(Key, MerkleDispatches)
	}
}
