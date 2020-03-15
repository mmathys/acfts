package benchmark

import (
	"fmt"
	"github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"math/rand"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"
)

func getRandomAddress(a common.Agent) common.Address {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	i := r.Intn(len(a.Topology) - 2)
	m := a.Topology[i]
	if reflect.DeepEqual(a.Address, m) {
		i++
	}
	return a.Topology[i]
}

/**
Clients send 1 money to random other clients
*/
func simpleAgent(a common.Agent, wg *sync.WaitGroup) {
	defer wg.Done()

	sendBufferLen := 1      // sync
	receiveBufferLen := 256 // async

	w := util.NewWalletWithAmount(a.Address, a.NumTransactions)

	incoming := make(chan common.Value, sendBufferLen)
	outgoing := make(chan common.Transaction, receiveBufferLen)

	go client.HandleIncoming(w, incoming)
	go client.HandleOutgoing(w, outgoing)

	mux := http.NewServeMux()

	mux.HandleFunc("/transaction", client.ReceiveSignature(incoming))
	localAddr := fmt.Sprintf(":%d", core.GetPort(a.Address))
	go http.ListenAndServe(localAddr, mux)

	time.Sleep(a.StartDelay) // wait before starting tx

	for i := 0; i < a.NumTransactions; i++ {
		to := getRandomAddress(a)
		t, err := wallet.PrepareTransaction(w, to, 1)
		if err != nil {
			fmt.Println(err)
			panic("failed to prepare transaction")
		}

		outgoing <- t
	}
	time.Sleep(a.EndDelay) // wait for others?
}

// there are 16 clients
func TestAgents(t *testing.T) {
	numTx := 100
	delay := 500 * time.Millisecond
	endDelay := 1 * time.Second
	clients := core.GetClients()
	var wg sync.WaitGroup

	for _, addr := range clients {
		a := common.Agent{NumTransactions: numTx, StartDelay: delay, EndDelay: endDelay, Address: addr, Topology: clients}
		wg.Add(1)
		go simpleAgent(a, &wg)
	}

	wg.Wait()
	fmt.Println("ended.")
}
