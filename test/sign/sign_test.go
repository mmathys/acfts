package sign

import (
	"fmt"
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/adapter"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

var numWorkers = 8

// This sets up the environment and the profiler.
func TestMain(m *testing.M) {
	numWorkers, err := strconv.Atoi(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Printf("numWorkers must be supplied")
		panic(err)
	}
	fmt.Printf("numWorkers = %d\n", numWorkers)

	go func() {
		runtime.SetBlockProfileRate(1)
		log.Println(http.ListenAndServe(":6666", nil))
	}()

	os.Exit(m.Run())
}

// Benchmarks the speed of the whole server (without network) for a variable number of workers. The topolpgy and the
// number of workers is passed as the last argument in the command line.
// Hint: this easiest way to run this test is with docker-compose.
func BenchmarkSignNoNetwork(b *testing.B) {
	err := worker(b.N, numWorkers, b)
	if err != nil {
		b.Error(err)
		b.Fail()
	}
}

// Benchmarks the speed of the whole server (without network) for 50k transactions and for a variable number of workers.
// The topolpogy and the number of workers are passed as the last argument in the command line.
// Hint: this easiest way to run this test is with docker-compose.
func TestSignNoNetwork(t *testing.T) {
	err := worker(50000, numWorkers, nil)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

// This function is used by the test and benchmarks. It contains some tests about whether a delay/distribution has an
// effect on the profile
func worker(N int, numWorkers int, b *testing.B) error {
	args := os.Args
	topo := args[len(args)-2]
	common.InitAddresses(topo)
	adapter.TxCounter = new(int32)
	adapter.SignedUTXO = new(sync.Map)
	adapter.CheckTransactions = true
	adapter.Benchmark = false
	adapter.Id = common.GetIdentity(common.GetServers()[0])
	adapter.AllowDoublespend = false
	adapter.UseUTXOMap = true
	adapter.CheckTransactions = true

	client := common.GetClients()[0]
	clientId := common.GetIdentity(client)
	target := common.GetClients()[1]

	var requests [][]common.TransactionSigReq
	for i := 0; i < numWorkers; i++ {
		requests = append(requests, []common.TransactionSigReq{})
		for j := 0; j < N/numWorkers; j++ {
			w := common.NewWalletWithAmount(client, 1)
			tx := common.Transaction{Inputs: nil, Outputs: nil}
			tx, err := core.PrepareTransaction(w, target, 1)
			if err != nil {
				return err
			}
			req := common.TransactionSigReq{Transaction: tx}
			err = common.SignTransactionSigRequest(clientId.Key, &req)
			if err != nil {
				return err
			}

			requests[i] = append(requests[i], req)
		}
	}

	server := new(adapter.Server)
	res := common.TransactionSignRes{}

	startDelay := 1 * time.Millisecond / time.Duration(numWorkers) // distribute start over 1ms
	var wg sync.WaitGroup

	if b != nil {
		b.ResetTimer()
	}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		time.Sleep(startDelay)
		go func(work []common.TransactionSigReq) {
			for j := 0; j < N/numWorkers; j++ {
				err := server.Sign(work[j], &res)
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}(requests[i])
	}

	wg.Wait()
	return nil
}
