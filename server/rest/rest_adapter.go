package rest

import (
	"encoding/json"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server"
	"github.com/mmathys/acfts/util"
	"net/http"
	"sync"
)

type Adapter struct {}

func (a *Adapter) Init(port int, id *common.Identity, debug bool, benchmark bool, SignedUTXO *sync.Map, TxCounter *int32) {
	http.HandleFunc("/sign", handleSign(id, debug, benchmark, SignedUTXO, TxCounter))
	localAddr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(localAddr, nil)
}

func handleSign(id *common.Identity, debug bool, benchmarkMode bool, SignedUTXO *sync.Map, TxCounter *int32) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		if benchmarkMode {
			defer util.CountTx(TxCounter)
		}

		// parse the request
		var sigReq common.TransactionSigReq

		err := json.NewDecoder(req.Body).Decode(&sigReq)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !debug {
			err = server.CheckValidity(id, &sigReq)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		tx := sigReq.Transaction
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
			outputs, err = common.SignValues(id.Key, tx.Outputs)
		}

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Form the response
		response := common.TransactionSignRes{Outputs: outputs}
		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
