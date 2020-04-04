package benchmark

import (
	"fmt"
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"testing"
)

/**
This is parallel benchmark
*/

// in this benchmark, a wallet gets created once. Then, the wallet spends all of its cash, 1 money per iteration.
func TestParallelSpendSingle(t *testing.T) {
	var numWorkers uint8 = 3
	N := 10000

	jobs := make(chan bool, N)
	done := make(chan bool, N)

	var i uint8 = 0
	for ; i < numWorkers; i++ {
		addr := common.GetClients()[i]
		w := util.NewWalletWithAmount(addr, N)
		go worker(w, t, jobs, done)
	}

	for j := 0; j < N; j++ {
		jobs <- true
	}
	close(jobs)

	for k := 0; k < N; k++ {
		<-done
	}
}

func worker(w *common.Wallet, t *testing.T, jobs <-chan bool, done chan<- bool) {
	targetAddr := common.GetClients()[0]
	for _ = range jobs {
		tx, err := wallet.PrepareTransaction(w, targetAddr, 1)
		if err != nil {
			t.Error("failed to prepare transaction")
		}

		_, err = core.SignTransaction(w, tx)
		if err != nil {
			fmt.Println("failed to sign transaction")
			return
		}

		done <- true
	}
}
