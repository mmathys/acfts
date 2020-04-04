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

func BenchmarkSignNoNetwork(t *testing.B) {
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
		t.Error(err)
	}

	req := common.TransactionSigReq{Transaction: tx}
	err = common.SignTransactionSigRequest(clientId.Key, &req)
	if err != nil {
		t.Error(err)
	}

	server := new(rpc.Server)

	res := common.TransactionSignRes{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		err := server.Sign(req, &res)
		if err != nil {
			t.Error(err)
		}
	}
}
