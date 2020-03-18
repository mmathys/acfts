package benchmark

import (
	"fmt"
	"github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"testing"
)

/**
This is parallel benchmark
*/

var targetAddr = common.Address{0}

// in this benchmark, a wallet gets created once. Then, the wallet spends all of its cash, 1 money per iteration.
func BenchmarkParallelSpendSingle(b *testing.B) {
	var numWorkers uint8 = 1
	N := 100000

	for ; numWorkers <= 16; numWorkers++ {
		s := fmt.Sprintf("ParallelSpendSingle, workers=%d", numWorkers)
		b.Run(s, func(b *testing.B) {

			jobs := make(chan bool, N)
			done := make(chan bool, N)

			var i uint8 = 0
			for ; i < numWorkers; i++ {
				addr := common.Address{i + 1}
				w := util.NewWalletWithAmount(addr, N)
				go worker(w, b, jobs, done)
			}

			for j := 0; j < N; j++ {
				jobs <- true
			}
			close(jobs)

			for k := 0; k < N; k++ {
				<-done
			}
		})
	}
}

func worker(w *common.Wallet, b *testing.B, jobs <-chan bool, done chan<- bool) {
	for _ = range jobs {
		tx, err := wallet.PrepareTransaction(w, targetAddr, 1)
		if err != nil {
			b.Error("failed to prepare transaction")
		}

		_, err = client.SignTransaction(w, tx)
		if err != nil {
			fmt.Println("failed to sign transaction")
			return
		}

		done <- true
	}
}
