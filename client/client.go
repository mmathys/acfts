package client

import (
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"github.com/mmathys/acfts/wallet"
	"reflect"
	"sync"
)

const bufferLen int = 255

func HandleIncoming(w *common.Wallet, incoming chan common.Value) {
	for {
		t := <-incoming
		//fmt.Printf("got tuple %v\n", t)
		wallet.AddUTXO(w, t)
	}
}

func HandleOutgoing(w *common.Wallet, outgoing chan common.Transaction) {
	for {
		t := <-outgoing
		//fmt.Printf("handle outgoing %v\n", t)
		doTransaction(w, t)
	}
}

func doTransaction(w *common.Wallet, t common.Transaction) {
	res, err := SignTransaction(w, t)
	if err != nil {
		fmt.Println("failed to sign transaction")
		return
	}

	// own UTXOs, (is spent at this point)
	wallet.RemoveUTXOMultiple(w, &t.Inputs)

	// TODO combine signatures.
	sig := (*res)[0]

	// add own outputs
	var ownOutputs []common.Value
	for _, t := range sig.Outputs {
		if reflect.DeepEqual(t.Address, w.Key.PublicKey) {
			ownOutputs = append(ownOutputs, t)
		} else {
			go ForwardSignature(t)
		}
	}

	wallet.AddUTXOMultiple(w, &ownOutputs)
}

// TODO Only wait for Math.ceil(2/3 * n) of n servers!
func SignTransaction(w *common.Wallet, t common.Transaction) (*[]common.TransactionSignRes, error) {
	n := len(core.GetServers())
	sigs := make(chan common.TransactionSignRes, n)

	var wg sync.WaitGroup

	for _, server := range core.GetServers() {
		wg.Add(1)
		go RequestSignature(server, t, &wg, &sigs)
	}

	wg.Wait()

	//fmt.Println("got sigs from all servers")

	// TODO validate and store sigs
	var res []common.TransactionSignRes
	for i := 0; i < n; i++ {
		sig := <-sigs
		res = append(res, sig)
	}

	return &res, nil
}
