package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"net/http"
	"reflect"
	"sync"
)

const bufferLen int = 255

func handleIncomingTransactions(incoming chan common.Transaction) {
	for {
		t := <- incoming
		fmt.Printf("got transaction %v\n", t)
	}
}

func handleOutgoingTransactions(outgoing chan common.Transaction) {
	for {
		t := <- outgoing
		fmt.Printf("handle outgoing %v\n", t)
		go doTransaction(t)
	}
}

func doTransaction(t common.Transaction) {
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
	wallet.RemoveUTXO(&t.Inputs)

	// add own outputs
	var ownOutputs []common.Tuple
	for _, t := range sig.Outputs {
		if reflect.DeepEqual(t.Address, core.GetOwnAddress()) {
			ownOutputs = append(ownOutputs, t)
		}
	}
	wallet.AddUTXO(&ownOutputs)

	// TODO: send UTXO to counterpart
}

func requestSignature(serverAddr common.Address, t common.Transaction, wg *sync.WaitGroup, sigs *chan common.TransactionSignRes) {
	net, err := core.LookupNetworkFromAddress(serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&t)
	if err != nil {
		fmt.Println("could not encode transaction")
		return
	}

	res, err := http.Post(net + "/sign", "raw", &buf)
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

func main() {
	incoming := make(chan common.Transaction, bufferLen)
	outgoing := make(chan common.Transaction, bufferLen)

	go handleIncomingTransactions(incoming)
	go handleOutgoingTransactions(outgoing)

	util.LaunchClientConsole(outgoing)
}