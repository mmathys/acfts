package server

import (
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)


var id *common.Identity
var debug bool
var benchmarkMode bool
var SignedUTXO *sync.Map
var TxCounter *int32

type Server struct {}
func (s *Server) Sign(req common.TransactionSigReq, res *common.TransactionSignRes) error {
	//log.Printf("got sign request: %v", req)

	if benchmarkMode {
		defer util.CountTx(TxCounter)
	}

	if !debug {
		err := CheckValidity(id, &req)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	tx := req.Transaction
	if !debug {
		for _, input := range tx.Inputs {
			index := [common.IdentifierLength]byte{}
			copy(index[:], input.Id[:common.IdentifierLength])
			SignedUTXO.Store(index, input)
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

func initRPC(port int, _id *common.Identity, _debug bool, _benchmark bool, _SignedUTXO *sync.Map, _TxCounter *int32) {
	id = _id
	debug = _debug
	benchmarkMode  = _benchmark
	SignedUTXO = _SignedUTXO
	TxCounter = _TxCounter

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
