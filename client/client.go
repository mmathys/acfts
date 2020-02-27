package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"github.com/mmathys/acfts/wallet"
	"net/http"
	"reflect"
	"sync"
)

const bufferLen int = 255

func HandleIncomingTransactions(w *common.Wallet, incoming chan common.Transaction) {
	for {
		t := <-incoming
		fmt.Printf("got transaction %v\n", t)
	}
}

func HandleOutgoingTransactions(w *common.Wallet, outgoing chan common.Transaction) {
	for {
		t := <-outgoing
		fmt.Printf("handle outgoing %v\n", t)
		go doTransaction(w, t)
	}
}

func doTransaction(w *common.Wallet, t common.Transaction) {
	// wait for all servers
	var wg sync.WaitGroup

	// TODO Only wait for Math.ceil(2/3 * n) of n servers!
	n := len(core.GetServers())
	sigs := make(chan common.TransactionSignRes, n)

	for _, server := range core.GetServers() {
		wg.Add(1)
		go requestSignature(server, t, &wg, &sigs)
	}

	wg.Wait()

	fmt.Println("got sigs from all servers")

	// TODO validate and store sigs
	sig := <-sigs

	// remove own UTXOs, (is spent at this point)
	wallet.RemoveUTXO(w, &t.Inputs)

	// add own outputs
	var ownOutputs []common.Tuple
	for _, t := range sig.Outputs {
		if reflect.DeepEqual(t.Address, w.Address) {
			ownOutputs = append(ownOutputs, t)
		} else {
			go forwardSignature(t)
		}
	}

	wallet.AddUTXO(w, &ownOutputs)
}

func requestSignature(serverAddr common.Address, t common.Transaction, wg *sync.WaitGroup, sigs *chan common.TransactionSignRes) {
	net, err := core.LookupNetworkFromAddress(serverAddr)
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
		fmt.Printf("could not fetch sig at %s\n", net)
		return
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

func forwardSignature(t common.Tuple) {
	net, err := core.LookupNetworkFromAddress(t.Address)
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
		fmt.Printf("failed forwarding tx to %s\n", net)
	} else {
		fmt.Println("tx forwarded successfully")
	}
}
