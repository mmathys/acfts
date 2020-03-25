package rest

/**
Includes functions, which clients use to communicate to each other. Handles transport, serialization and deserialization.
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmathys/acfts/common"
	"log"
	"net/http"
	"sync"
)

type Adapter struct {}

func (a *Adapter) Init(port int, incoming chan common.Value) {
	http.HandleFunc("/transaction", receiveSignatureREST(incoming))
	localAddr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(localAddr, nil)
}

func (a *Adapter) RequestSignature(serverAddr common.Address, id *common.Identity, t common.Transaction, wg *sync.WaitGroup, sigs *chan common.TransactionSignRes) {
	net, err := common.GetNetworkAddress(serverAddr)
	net = "http://" + net
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	request := common.TransactionSigReq{Transaction: t}
	err = common.SignTransactionSigRequest(id.Key, &request)
	if err != nil{
		log.Panic(err)
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&request)
	if err != nil {
		fmt.Println("could not encode transaction sign request")
		return
	}

	res, err := http.Post(net+"/sign", "raw", &buf)
	if err != nil {
		msg := fmt.Sprintf("could not fetch sig at %s\n", net)
		fmt.Println(err)
		panic(msg)
	}

	var sig common.TransactionSignRes
	err = json.NewDecoder(res.Body).Decode(&sig)
	if err != nil {
		fmt.Println("could not decode transaction sig response")
		return
	}

	*sigs <- sig
	wg.Done()
}

func (a *Adapter) ForwardValue(t common.Value) {
	net, err := common.GetNetworkAddress(t.Address)
	net = "http://" + net
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&t)
	if err != nil {
		fmt.Println("could not encode transaction")
		return
	}

	res, err := http.Post(net+"/transaction", "raw", &buf)
	if err != nil || res.StatusCode != 200 {
		msg := fmt.Sprintf("failed forwarding tx to %s.\n", net)
		fmt.Println(err)
		panic(msg)
	} else {
		//fmt.Println("tx forwarded successfully")
	}
}

func receiveSignatureREST(c chan common.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Parse the request
		var t common.Value
		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// put into channel
		c <- t

		w.WriteHeader(200)
	}
}
