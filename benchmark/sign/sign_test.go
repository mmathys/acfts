package sign

import (
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/rpc"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
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

func TestMain(m *testing.M) {
	numWorkers, err := strconv.Atoi(os.Args[len(os.Args)-1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("numWorkers = %d\n", numWorkers)

	go func() {
		runtime.SetBlockProfileRate(1)
		log.Println(http.ListenAndServe(":6666", nil))
	}()

	os.Exit(m.Run())
}

func BenchmarkSignNoNetwork(b *testing.B) {
	numWorkers, err := strconv.Atoi(os.Args[len(os.Args)-1])
	if err != nil {
		panic(err)
	}

	err = worker(b.N, numWorkers, b)
	if err != nil {
		b.Error(err)
		b.Fail()
	}
}

func TestSignNoNetwork(t *testing.T) {
	numWorkers, err := strconv.Atoi(os.Args[len(os.Args)-1])
	if err != nil {
		panic(err)
	}

	err = worker(50000, numWorkers, nil)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func worker(N int, numWorkers int, b *testing.B) error {
	fmt.Println("preparing...")
	args := os.Args
	topo := args[len(args)-2]
	common.InitAddresses(topo)
	rpc.TxCounter = new(int32)
	rpc.SignedUTXO = new(sync.Map)
	rpc.Debug = false
	rpc.BenchmarkMode = false
	rpc.Id = util.GetIdentity(common.GetServers()[0])
	rpc.AllowDoublespend = false
	rpc.UseUTXOMap = true
	rpc.CheckTransactions = true

	client := common.GetClients()[0]
	clientId := util.GetIdentity(client)
	target := common.GetClients()[1]

	var requests [][]common.TransactionSigReq
	for i := 0; i < numWorkers; i++ {
		requests = append(requests, []common.TransactionSigReq{})
		for j := 0; j < N/numWorkers; j++ {
			w := util.NewWalletWithAmount(client, 1)
			tx := common.Transaction{Inputs: nil, Outputs: nil}
			tx, err := wallet.PrepareTransaction(w, target, 1)
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

	server := new(rpc.Server)
	res := common.TransactionSignRes{}

	startDelay := 1 * time.Millisecond / time.Duration(numWorkers) // distribute start over 1ms
	var wg sync.WaitGroup

	fmt.Println("running tests... ")
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
