package test

import (
	"fmt"
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/test/environment"
	"os"
	"testing"
)

/**
Easy synchronous benchmark
Do not expect high numbers from this
*/

var A common.Address
var B common.Address

func TestMain(m *testing.M) {
	common.InitAddresses("../topologies/localSimple.json")
	A = environment.TestClient(0)
	B = environment.TestClient(1)
	os.Exit(m.Run())
}

// in this benchmark, in each iteration, a new wallet gets created. then, the wallet spends all of its cash.
func BenchmarkSequentialNewWallet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		walletA := common.NewWallet(A)

		tx, err := core.PrepareTransaction(walletA, B, 100)
		if err != nil {
			b.Error("failed to prepare transaction")
		}

		_, err = core.SignTransaction(walletA, tx)
		if err != nil {
			fmt.Println("failed to sign transaction")
			return
		}
	}
}

// in this benchmark, a wallet gets created once. Then, the wallet spends all of its cash, 1 money per iteration.
func BenchmarkSequentialSpendSingle(b *testing.B) {
	var addrA = common.GetClients()[0]
	var addrB = common.GetClients()[1]
	A := common.NewWalletWithAmount(addrA, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx, err := core.PrepareTransaction(A, addrB, 1)
		if err != nil {
			b.Error("failed to prepare transaction")
		}

		_, err = core.SignTransaction(A, tx)
		if err != nil {
			fmt.Println("failed to sign transaction")
			return
		}
	}
}

func TestSequentialSpendSingle(t *testing.T) {
	N := 10000
	var addrA = common.GetClients()[0]
	var addrB = common.GetClients()[1]
	A := common.NewWalletWithAmount(addrA, N)

	for i := 0; i < N; i++ {
		tx, err := core.PrepareTransaction(A, addrB, 1)
		if err != nil {
			panic("failed to prepare transaction")
		}

		_, err = core.SignTransaction(A, tx)
		if err != nil {
			fmt.Println("failed to sign transaction")
			return
		}
	}
}
