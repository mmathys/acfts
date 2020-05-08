package checks

import (
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/common/test_util"
	"os"
	"testing"
)

var A common.Address
var B common.Address

func TestMain(m *testing.M) {
	test_util.TestEnvironment()
	A = test_util.TestClient(0)
	B = test_util.TestClient(1)
	os.Exit(m.Run())
}

func TestValidSig(t *testing.T) {
	w := common.NewWallet(A)
	tx, _ := core.PrepareTransaction(w, B, 1)
	_, err := core.SignTransaction(w, tx)

	if err != nil {
		t.Error(err)
	}
}

func TestMinZero(t *testing.T) {
	w := common.NewWallet(A)

	tx, _ := core.PrepareTransaction(w, B, 0)
	_, err := core.SignTransaction(w, tx)

	if err == nil {
		t.Error("should throw an error")
	}

	tx, _ = core.PrepareTransaction(w, B, -1)
	_, err = core.SignTransaction(w, tx)

	if err == nil {
		t.Error("should throw an error")
	}
}
