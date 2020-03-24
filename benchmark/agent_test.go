package benchmark

import (
	"bytes"
	"fmt"
	"github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func getRandomAddress(a common.Agent) common.Address {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	i := r.Intn(len(a.Topology) - 1)
	m := a.Topology[i]
	if bytes.Equal(a.Address, m) {
		i++
	}
	return a.Topology[i]
}

/**
Clients send 1 money to random other clients
*/
func simpleAgent(a common.Agent, wg *sync.WaitGroup) {
	defer wg.Done()

	w := util.NewWalletWithAmount(a.Address, a.NumTransactions)

	time.Sleep(a.StartDelay) // wait before starting tx

	for i := 0; i < a.NumTransactions; i++ {
		to := getRandomAddress(a)

		t, err := wallet.PrepareTransaction(w, to, 1)
		if err != nil {
			fmt.Println(err)
			panic("failed to prepare transaction")
		}

		client.DoTransaction(w, t, false)
	}
}

// there are 16 clients
func testAgents(t *testing.T) {
	maxClients := 9
	numTx := 100000/maxClients
	delay := 500 * time.Millisecond
	clients := common.GetClients()

	var wg sync.WaitGroup
	topology := clients[:maxClients]

	for _, addr := range topology {
		a := common.Agent{NumTransactions: numTx, StartDelay: delay, Address: addr, Topology: topology}
		wg.Add(1)
		go simpleAgent(a, &wg)
	}

	wg.Wait()
}

func TestAgentsREST(t *testing.T) {
	testAgents(t)
}

func TestAgentsRPC(t *testing.T) {
	client.SetAdapterMode("rpc")
	testAgents(t)
}