package client

import (
	"fmt"
	"github.com/mmathys/acfts/common"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

var incoming chan common.Value

type Client struct {}
type Empty struct {}

func (c *Client) ForwardSignature(req common.Value, res *Empty) error {
	incoming <- req
	*res = Empty{}
	return nil
}

func initRPC(port int, _incoming chan common.Value) {
	incoming = _incoming

	client := new(Client)
	rpc.Register(client)
	rpc.HandleHTTP()
	addr := fmt.Sprintf(":%d", port)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

var once sync.Once
var client *rpc.Client

func requestSignatureRPC(serverAddr common.Address, id *common.Identity, t common.Transaction, wg *sync.WaitGroup, sigs *chan common.TransactionSignRes) {
	net, err := common.GetNetworkAddress(serverAddr)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	req := common.TransactionSigReq{Transaction: t}
	err = common.SignTransactionSigRequest(id.Key, &req)
	if err != nil{
		log.Panic(err)
	}

	once.Do(func() {
		client, err = rpc.DialHTTP("tcp", net)
	})
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var res common.TransactionSignRes
	err = client.Call("Server.Sign", req, &res)
	if err != nil {
		msg := fmt.Sprintf("could not fetch sig at %s\n", net)
		fmt.Println(err)
		log.Panic(msg)
	}

	*sigs <- res
	wg.Done()
}

func forwardValueRPC(t common.Value) {
	net, err := common.GetNetworkAddress(t.Address)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	client, err := rpc.DialHTTP("tcp", net)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var res Empty
	err = client.Call("Client.ForwardSignature", t, &res)
	if err != nil {
		msg := fmt.Sprintf("failed forwarding tx to %s.\n", net)
		fmt.Println(err)
		log.Panic(msg)
	} else {
		//fmt.Println("tx forwarded successfully")
	}
}