package benchmark

import (
	"fmt"
	"github.com/mmathys/acfts/client/core"
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

		core.DoTransaction(w, t, false)
	}
}

// there are 16 clients
func testAgentsMultipleParallel(t *testing.T) {
	clients := common.GetClients()
	maxClients := int(.3 * float64(len(clients)))
	//maxClients = 1000

	for numClients := maxClients; numClients <= maxClients; numClients++ {
		testAgents(t, numClients)
	}
}

func testAgents(t *testing.T, numClients int) {
	delay := 500 * time.Millisecond
	totalTx := 1000000 // 1 million
	clients := common.GetClients()
	msg := fmt.Sprintf("numClients: %d", numClients)
	t.Run(msg, func(t *testing.T) {
		numTx := totalTx / numClients
		var wg sync.WaitGroup
		topology := clients[:numClients]

		for _, addr := range topology {
			wg.Add(1)
			go simpleAgent(common.Agent{NumTransactions: numTx, StartDelay: delay, Address: addr, Topology: topology}, &wg)
		}

		wg.Wait()
	})
}

func TestAgentsREST(t *testing.T) {
	common.InitAddresses("../topologies/localSimple.json")
	testAgentsMultipleParallel(t)
}

func TestAgentsRPC(t *testing.T) {
	core.SetAdapterMode("rpc")
	common.InitAddresses("../topologies/localSimple.json")
	testAgentsMultipleParallel(t)
}

func TestAgentsAWS(t *testing.T) {
	core.SetAdapterMode("rpc")
	common.InitAddresses("../topologies/aws.json")
	testAgentsMultipleParallel(t)
}

func TestAgentsVSOS(t *testing.T) {
	core.SetAdapterMode("rpc")
	common.InitAddresses("../topologies/vsos.json")
	testAgentsMultipleParallel(t)
}

func TestAgentsRPCFixed(t *testing.T) {
	args := os.Args

	numClients, err := strconv.Atoi(args[len(args)-1])
	if err != nil {
		panic(err)
	}

	core.SetAdapterMode("rpc")

	topo := args[len(args)-2]
	common.InitAddresses(topo)

	testAgents(t, numClients)
}
