package parallel

import (
	"fmt"
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/test/environment"
	"os"
	"testing"
)

/**
This is parallel benchmark
*/


var A common.Address
var B common.Address

func TestMain(m *testing.M) {
	common.InitAddresses("../../topologies/localSimple.json")
	A = environment.TestClient(0)
	B = environment.TestClient(1)
	os.Exit(m.Run())
}

// in this benchmark, a wallet gets created once. Then, the wallet spends all of its credits, 1 credit per iteration.
func TestParallelSpendSingle(t *testing.T) {
	var numWorkers = 3
	N := 10000

	jobs := make(chan bool, N)
	done := make(chan bool, N)

	var i = 0
	for ; i < numWorkers; i++ {
		addr := environment.TestClient(i)
		w := common.NewWalletWithAmount(addr, N)
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
		tx, err := core.PrepareTransaction(w, targetAddr, 1)
		if err != nil {
			t.Error("failed to prepare transaction")
		}

		core.DoTransaction(w, tx, false)
		if err != nil {
			fmt.Println("failed to sign transaction")
			return
		}

		done <- true
	}
}
