package sign

import (
	"fmt"
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/adapter"
	"github.com/mmathys/acfts/server/store"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"testing"
)

var numMultisig = 0
var numWorkers = 0
var topology = "none"
var batchVerification = true
var mapType = -1

// This sets up the environment and the profiler.
func TestMain(m *testing.M) {
	_numMultisig, err := strconv.Atoi(os.Getenv("NUM_MULTISIG"))
	if err != nil {
		panic("Error getting environment variable NUM_MULTISIG")
	}
	numMultisig = _numMultisig

	_numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))
	if err != nil {
		panic("Error getting environment variable NUM_WORKERS")
	}
	numWorkers = _numWorkers

	_topology := os.Getenv("TOPOLOGY")
	if _topology == "" {
		panic("Environment variable TOPOLOGY is not set")
	}
	topology = _topology

	_batchVerification := os.Getenv("BATCH_VERIFICATION")

	if _batchVerification == "true" {
		batchVerification = true
	} else if _batchVerification == "false" {
		batchVerification = false
	} else if _batchVerification == "" {
		panic("Environment variable BATCH_VERIFICATION is not set")
	} else {
		panic("Environment variable BATCH_VERIFICATION must be either \"true\" or \"false\"")
	}

	_mapType := os.Getenv("MAP_TYPE")
	if _mapType == "insertOnly" {
		mapType = store.TypeInsertOnly
	} else if _mapType == "syncMap" {
		mapType = store.TypeSyncMap
	} else if _mapType == "" {
		mapType = store.DefaultMapType
	} else {
		panic("Environment variable MAP_TYPE must be \"insertOnly\", \"syncMap\" or empty.")
	}

	go func() {
		//runtime.SetBlockProfileRate(1)
		log.Println(http.ListenAndServe(":6666", nil))
	}()

	os.Exit(m.Run())
}

// Benchmarks the speed of the whole server (without network) for a variable number of workers. The topolpgy and the
// number of workers is passed as the last argument in the command line.
// Hint: this easiest way to run this test is with docker-compose.
func BenchmarkSignNoNetwork(b *testing.B) {
	err := worker(b.N, b)
	if err != nil {
		b.Error(err)
		b.Fail()
	}
}

// Benchmarks the speed of the whole server (without network) for 50k transactions and for a variable number of workers.
// The topolpogy and the number of workers are passed as the last argument in the command line.
// Hint: this easiest way to run this test is with docker-compose.
func TestSignNoNetwork(t *testing.T) {
	err := worker(1000000, nil)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

// This function is used by the test and benchmarks.
func worker(N int, b *testing.B) error {
	mapTypeString := "unrecognized"
	if mapType == store.TypeSyncMap {
		mapTypeString = "syncMap"
	} else if mapType == store.TypeInsertOnly {
		mapTypeString = "insertOnly"
	}

	fmt.Printf("topology=%s, numWorkers=%d, numMultisig=%d, batchVerification=%t, mapType=%s\n",
		topology, numWorkers, numMultisig, batchVerification, mapTypeString)

	common.InitAddresses(topology)
	// set the number of server according to numMultisig. numMultisig must be >= num servers.
	if len(common.ServerAddresses) < numMultisig {
		log.Panicf("not enough servers. numMultisig=%d, but we only have %d servers.", numMultisig, len(common.ServerAddresses))
		panic("not enough servers")
	}
	common.ServerAddresses = common.ServerAddresses[:numMultisig]

	// initialize adapter
	utxoMap := new(store.UTXOMap)
	utxoMap.SetType(mapType)
	utxoMap.Init()
	adapter.UTXOMap = utxoMap
	adapter.TxCounter = new(int32)
	adapter.CheckTransactions = true
	adapter.Benchmark = false
	adapter.Key = common.GetIdentity(common.GetServers()[0])
	adapter.AllowDoublespend = false
	adapter.UseUTXOMap = true
	adapter.CheckTransactions = true
	adapter.BatchVerification = batchVerification

	// get clients from topology
	client := common.GetClients()[0]
	clientId := common.GetIdentity(client)
	target := common.GetClients()[1]

	// generate requests
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
			err = common.SignTransactionSigRequest(clientId, &req)
			if err != nil {
				return err
			}

			requests[i] = append(requests[i], req)
		}
	}

	server := new(adapter.Server)
	res := common.TransactionSignRes{}

	//startDelay := 1 * time.Millisecond / time.Duration(numWorkers) // distribute start over 1ms
	var wg sync.WaitGroup

	if b != nil {
		b.ResetTimer()
	}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		//time.Sleep(startDelay)
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
