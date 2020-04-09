package sign

import (
	"fmt"
	"github.com/cornelk/hashmap"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/rpc"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	_ "net/http/pprof"
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	/*
	go func() {
		runtime.SetBlockProfileRate(1)
		log.Println(http.ListenAndServe("localhost:6666", nil))
	}()
	 */
	os.Exit(m.Run())
}

func BenchmarkSignNoNetwork(b *testing.B) {
	for numWorkers := 0; numWorkers < 100; numWorkers++ {
		name := fmt.Sprintf("num workers = %d", numWorkers)
		b.Run(name, func(b *testing.B) {
			err := worker(b.N, numWorkers, b)
			if err != nil {
				b.Error(err)
				b.Fail()
			}
		})
	}


}

/*
func TestSignNoNetwork(t *testing.T) {
	numWorkers, err := strconv.Atoi(os.Args[len(os.Args)-1])
	if err != nil {
		panic(err)
	}

	err = worker(1000000, numWorkers, nil)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
 */

func worker(N int, numWorkers int, b *testing.B) error {
	args := os.Args
	topo := args[len(args)-1]
	common.InitAddresses(topo)
	rpc.TxCounter = new(int32)
	rpc.SignedUTXO = new(hashmap.HashMap)
	rpc.Debug = false
	rpc.BenchmarkMode = false
	rpc.Id = util.GetIdentity(common.GetServers()[0])
	rpc.AllowDoublespend = true

	client := common.GetClients()[0]
	clientId := util.GetIdentity(client)
	target := common.GetClients()[1]
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

	server := new(rpc.Server)

	res := common.TransactionSignRes{}
	if b != nil {
		b.ResetTimer()
	}
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < N/numWorkers; j++ {
				err := server.Sign(req, &res)
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}
