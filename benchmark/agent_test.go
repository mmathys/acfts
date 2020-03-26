package benchmark

import (
	"fmt"
	"github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

/**
Clients send 1 money to random other clients
*/
func simpleAgent(a common.Agent, wg *sync.WaitGroup) {
	defer wg.Done()

	w := util.NewWalletWithAmount(a.Address, a.NumTransactions)

	time.Sleep(a.StartDelay) // wait before starting tx

	for i := 0; i < a.NumTransactions; i++ {
		to := a.Topology[0]

		t, err := wallet.PrepareTransaction(w, to, 1)
		if err != nil {
			fmt.Println(err)
			panic("failed to prepare transaction")
		}

		client.DoTransaction(w, t, false)
	}
}

// there are 16 clients
func testAgentsMultipleParallel(t *testing.T) {
	clients := common.GetClients()
	maxClients := len(clients)

	for numClients := 9; numClients <= maxClients; numClients++ {
		testAgents(t, numClients)
	}
}

func testAgents(t *testing.T, numClients int) {
	delay := 500 * time.Millisecond
	clients := common.GetClients()
	msg := fmt.Sprintf("numClients: %d", numClients)
	t.Run(msg, func(t *testing.T) {
		numTx := 100000
		var wg sync.WaitGroup
		topology := clients[:numClients]

		for _, addr := range topology {
			a := common.Agent{NumTransactions: numTx, StartDelay: delay, Address: addr, Topology: topology}
			wg.Add(1)
			go simpleAgent(a, &wg)
		}

		wg.Wait()
	})
}

func TestAgentsREST(t *testing.T) {
	common.InitAddresses("../topologies/localSimple.json")
	testAgentsMultipleParallel(t)
}

func TestAgentsRPC(t *testing.T) {
	client.SetAdapterMode("rpc")
	common.InitAddresses("../topologies/localSimple.json")
	testAgentsMultipleParallel(t)
}

func TestAgentsRPCFixed(t *testing.T) {
	args := os.Args

	numClients, err := strconv.Atoi(args[len(args)-1])
	if err != nil {
		panic(err)
	}

	client.SetAdapterMode("rpc")

	topo := args[len(args)-2]
	common.InitAddresses(topo)

	testAgents(t, numClients)
}
