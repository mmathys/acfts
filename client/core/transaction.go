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
		fmt.Println("failed to sign transaction:")
		fmt.Println(err)
		return
	}

	// own UTXOs, (is spent at this point)
	RemoveUTXOMultiple(w, &t.Inputs)

	sig := combineSignatures(res)

	// add own outputs
	var ownOutputs []common.Value
	for _, v := range sig.Outputs {
		if bytes.Equal(v.Address, w.GetAddress()) {
			ownOutputs = append(ownOutputs, v)
		} else if forward {
			go adapter.ForwardValue(v)
		}
	}

	AddUTXOMultiple(w, &ownOutputs)
}

// TODO Only request Math.ceil(2/3 * n) of n (randomly chosen) servers!
func SignTransaction(w *common.Wallet, t common.Transaction) (*[]common.TransactionSignRes, error) {
	n := len(common.GetServers())

	sigs := make(chan common.TransactionSignRes, n)
	errs := make(chan error, n)

	var wg sync.WaitGroup
	wg.Add(common.QuorumSize())
	for _, server := range common.ServerQuorum() {
		go adapter.RequestSignature(server, w.Key, w.Mode, t, &wg, sigs, errs)
	}
	wg.Wait()

	close(sigs)
	close(errs)

	var errors []error
	for err := range errs {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return nil, errors[0]
	}

	var res []common.TransactionSignRes
	for sig := range sigs {
		res = append(res, sig)
	}

	return &res, nil
}
