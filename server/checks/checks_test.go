package checks

import (
	client "github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	"testing"
)

func TestValidSig(t *testing.T) {
	w := common.NewWallet(common.Address{0})
	tx, _ := cli.PrepareTransaction(w, common.Address{1}, 1)
	_, err := client.SignTransaction(w, tx)

	if err != nil {
		t.Error(err)
	}
}

func TestMinZero(t *testing.T) {
	w := common.NewWallet(common.Address{0})

	tx, _ := cli.PrepareTransaction(w, common.Address{1}, 0)
	_, err := client.SignTransaction(w, tx)

	if err == nil {
		t.Error("should throw an error")
	}

	tx, _ = cli.PrepareTransaction(w, common.Address{1}, -1)
	_, err = client.SignTransaction(w, tx)

	if err == nil {
		t.Error("should throw an error")
	}
}
