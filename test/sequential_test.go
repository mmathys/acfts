package test

import (
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

// in this benchmark, in each iteration, a new wallet gets created. then, the wallet spends all of its credits.
func BenchmarkSequentialNewWallet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		walletA := common.NewWallet(A)

		tx, err := core.PrepareTransaction(walletA, B, 100)
		if err != nil {
			b.Fatal("failed to prepare transaction")
		}

		core.DoTransaction(walletA, tx, false)
	}
}

// in this benchmark, a wallet gets created once. Then, the wallet spends all of its credits, 1 credit per iteration.
// the server
func BenchmarkSequentialSpendSingle(b *testing.B) {
	walletA := common.NewWalletWithAmount(A, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx, err := core.PrepareTransaction(walletA, B, 1)
		if err != nil {
			b.Fatal("failed to prepare transaction")
		}

		core.DoTransaction(walletA, tx, false)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestSequentialSpendSingle(t *testing.T) {
	N := 10000
	walletA := common.NewWalletWithAmount(A, N)

	for i := 0; i < N; i++ {
		tx, err := core.PrepareTransaction(walletA, B, 1)
		if err != nil {
			panic("failed to prepare transaction")
		}

		core.DoTransaction(walletA, tx, false)
	}
}
