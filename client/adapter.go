package client

/**
Includes functions, which clients use to communicate to each other. Handles transport, serialization and deserialization.
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"net/http"
	"sync"
)

func RequestSignature(serverAlias common.Alias, t common.Transaction, wg *sync.WaitGroup, sigs *chan common.TransactionSignRes) {
	net, err := core.GetNetworkAddress(serverAlias)
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
	defer wg.Done()
}

func ForwardSignature(t common.Value) {
	alias, err := core.GetAliasFromAddress(t.Address)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	net, err := core.GetNetworkAddress(alias)
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

func ReceiveSignature(c chan common.Value) http.HandlerFunc {
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
