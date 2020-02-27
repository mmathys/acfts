package benchmark

import (
	"fmt"
	"github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"os"
	"testing"
)

/**
This is an easy (synchronous!) benchmark
A -> B and B -> A back.
Do not expect high numbers from this
*/

var addrA = common.Address{0}
var addrB = common.Address{1}
var A = util.NewWallet(addrA)
var B = util.NewWallet(addrB)

var bufferLen = 1 // important for benchmark
var A_in = make(chan common.Transaction, bufferLen)
var B_in = make(chan common.Transaction, bufferLen)
var A_out = make(chan common.Transaction, bufferLen)
var B_out = make(chan common.Transaction, bufferLen)

func TestMain(m *testing.M) {
	fmt.Println("before!")

	go client.HandleIncomingTransactions(A, A_in)
	go client.HandleIncomingTransactions(B, B_in)
	go client.HandleOutgoingTransactions(A, A_out)
	go client.HandleOutgoingTransactions(B, B_out)

	code := m.Run()

	fmt.Println("after!")
	os.Exit(code)
}

var mode bool = true // mode <=> A -> B
func BenchmarkBasic(b *testing.B) {
	var origin *common.Wallet
	var target common.Address
	var outCh *chan common.Transaction

	if mode {
		origin = A
		target = addrB
		outCh = &A_out
	} else {
		origin = B
		target = addrA
		outCh = &B_out
	}

	tx, err := wallet.PrepareTransaction(origin, target, 100)
	if err != nil {
		b.Error("failed to send transaction")
	}

	*outCh <- tx

	mode = !mode
}
