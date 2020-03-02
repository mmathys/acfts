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
Do not expect high numbers from this
*/

func TestMain(m *testing.M) {
	fmt.Println("before!")

	code := m.Run()

	fmt.Println("after!")
	os.Exit(code)
}


// in this benchmark, in each iteration, a new wallet gets created. then, the wallet spends all of its cash.
func BenchmarkSequentialNewWallet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var addrA = common.Address{0}
		var addrB = common.Address{1}

		A := util.NewWallet(addrA)

		tx, err := wallet.PrepareTransaction(A, addrB, 100)
		if err != nil {
			b.Error("failed to prepare transaction")
		}

		_, err = client.SignTransaction(A, tx)
		if err != nil {
			fmt.Println("failed to sign transaction")
			return
		}
	}
}


// in this benchmark, a wallet gets created once. Then, the wallet spends all of its cash, 1 money per iteration.
func BenchmarkSequentialSpendSingle(b *testing.B) {
	var addrA = common.Address{0}
	var addrB = common.Address{1}
	A := util.NewWalletWithAmount(addrA, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx, err := wallet.PrepareTransaction(A, addrB, 1)
		if err != nil {
			b.Error("failed to prepare transaction")
		}

		_, err = client.SignTransaction(A, tx)
		if err != nil {
			fmt.Println("failed to sign transaction")
			return
		}
	}
}
