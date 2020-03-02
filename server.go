package main

import (
	"encoding/json"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"net/http"
	"sync"
)

//var SignedUTXO = map[int]common.Tuple{}
var SignedUTXO sync.Map

func handleSign(w http.ResponseWriter, req *http.Request) {
	// Parse the request
	var tx common.Transaction
	err := json.NewDecoder(req.Body).Decode(&tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, input := range tx.Inputs {
		SignedUTXO.Store(input.Id, input)
	}

	//fmt.Println(SignedUTXO)

	// Sign the request
	outputs, err := core.Sign(&tx.Inputs, &tx.Outputs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Form the response
	response := common.TransactionSignRes{Outputs: outputs}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func main() {
	http.HandleFunc("/sign", handleSign)
	http.ListenAndServe(":6666", nil)
}
