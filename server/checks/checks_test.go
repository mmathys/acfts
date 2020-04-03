package checks

import (
	client "github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"testing"
)

func TestValidSig(t *testing.T) {
	w := util.NewWallet(common.Address{0})
	tx, _ := wallet.PrepareTransaction(w, common.Address{1}, 1)
	_, err := client.SignTransaction(w, tx)

	if err != nil {
		t.Error(err)
	}
}

func TestMinZero(t *testing.T) {
	w := util.NewWallet(common.Address{0})

	tx, _ := wallet.PrepareTransaction(w, common.Address{1}, 0)
	_, err := client.SignTransaction(w, tx)

	if err == nil {
		t.Error("should throw an error")
	}

	tx, _ = wallet.PrepareTransaction(w, common.Address{1}, -1)
	_, err = client.SignTransaction(w, tx)

	if err == nil {
		t.Error("should throw an error")
	}
}