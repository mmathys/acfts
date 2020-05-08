package core

import (
	"bytes"
	"fmt"
	"github.com/mmathys/acfts/client/adapter"
	"github.com/mmathys/acfts/common"
	"sync"
)

func DoTransaction(w *common.Wallet, t common.Transaction, forward bool) {
	res, err := SignTransaction(w, t)
	if err != nil {
		fmt.Println("failed to sign transaction")
		return
	}

	// own UTXOs, (is spent at this point)
	RemoveUTXOMultiple(w, &t.Inputs)

	sig := combineSignatures(res)

	// add own outputs
	var ownOutputs []common.Value
	for _, v := range sig.Outputs {
		if bytes.Equal(v.Address, w.Identity.Address) {
			ownOutputs = append(ownOutputs, v)
		} else if forward {
			go adapter.ForwardValue(v)
		}
	}

	AddUTXOMultiple(w, &ownOutputs)
}

// TODO Only wait for Math.ceil(2/3 * n) of n servers!
func SignTransaction(w *common.Wallet, t common.Transaction) (*[]common.TransactionSignRes, error) {
	n := len(common.GetServers())

	sigs := make(chan common.TransactionSignRes, n)

	var wg sync.WaitGroup

	for _, server := range common.GetServers() {
		wg.Add(1)
		go adapter.RequestSignature(server, w.Identity, t, &wg, &sigs)
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
