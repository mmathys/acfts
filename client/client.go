package client

import (
	"bytes"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/wallet"
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

func DoTransaction(w *common.Wallet, t common.Transaction, forward bool) {
	res, err := SignTransaction(w, t)
	if err != nil {
		fmt.Println("failed to sign transaction")
		return
	}

	// own UTXOs, (is spent at this point)
	wallet.RemoveUTXOMultiple(w, &t.Inputs)

	sig := combineSignatures(res)

	// add own outputs
	var ownOutputs []common.Value
	for _, v := range sig.Outputs {
		if bytes.Equal(v.Address, w.Identity.Address) {
			ownOutputs = append(ownOutputs, v)
		} else if forward {
			go ForwardValue(v)
		}
	}

	wallet.AddUTXOMultiple(w, &ownOutputs)
}

// TODO Only wait for Math.ceil(2/3 * n) of n servers!
func SignTransaction(w *common.Wallet, t common.Transaction) (*[]common.TransactionSignRes, error) {
	n := len(common.GetServers())

	sigs := make(chan common.TransactionSignRes, n)

	var wg sync.WaitGroup

	for _, server := range common.GetServers() {
		wg.Add(1)
		go RequestSignature(server, w.Identity, t, &wg, &sigs)
	}

	wg.Wait()

	// TODO validate and store sigs
	var res []common.TransactionSignRes
	for i := 0; i < n; i++ {
		sig := <-sigs
		res = append(res, sig)
	}

	return &res, nil
}
