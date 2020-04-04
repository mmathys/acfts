package sign

import (
	"github.com/cornelk/hashmap"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/rpc"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"os"
	"testing"
)

func run(N int, b *testing.B) error {
	args := os.Args
	topo := args[len(args)-1]
	common.InitAddresses(topo)
	rpc.TxCounter = new(int32)
	rpc.SignedUTXO = new(hashmap.HashMap)
	rpc.Debug = false
	rpc.BenchmarkMode = false
	rpc.Id = util.GetIdentity(common.GetServers()[0])
	rpc.AllowDoublespend = true

	client := common.GetClients()[0]
	clientId := util.GetIdentity(client)
	target := common.GetClients()[1]
	w := util.NewWalletWithAmount(client, 1)
	tx := common.Transaction{Inputs:nil, Outputs:nil}

	tx, err := wallet.PrepareTransaction(w, target, 1)

	if err != nil {
		return err
	}

	req := common.TransactionSigReq{Transaction: tx}
	err = common.SignTransactionSigRequest(clientId.Key, &req)
	if err != nil {
		return err
	}

	server := new(rpc.Server)

	res := common.TransactionSignRes{}
	if b != nil {
		b.ResetTimer()
	}
	for i := 0; i < N; i++ {
		err := server.Sign(req, &res)
		if err != nil {
			return err
		}
	}
	return nil
}

func BenchmarkSignNoNetwork(b *testing.B) {
	err := run(b.N, b)
	if err != nil {
		b.Error(err)
		b.Fail()
	}
}

func TestSignNoNetwork(t *testing.T) {
	err := run(1000000, nil)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
