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
	wg.Add(len(core.GetServers()))
	sigs := make(chan common.TransactionSignRes)
	
	for _, server := range core.GetServers() {
		go requestSignature(server, t, &wg, &sigs)
	}

	wg.Wait()

	// TODO validate and store sigs
	sig := <- sigs

	fmt.Println("got sigs from all servers")

	// remove and add UTXOs
	wallet.RemoveUTXO(&t.Inputs)
	wallet.AddUTXO(&sig.Outputs)

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
	wg.Done()
}

func main() {
	incoming := make(chan common.Transaction, bufferLen)
	outgoing := make(chan common.Transaction, bufferLen)

	go handleIncomingTransactions(incoming)
	go handleOutgoingTransactions(outgoing)

	util.LaunchClientConsole(outgoing)
}