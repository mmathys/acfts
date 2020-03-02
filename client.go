package main

import (
	"github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"net/http"
)

const bufferLen int = 255

func main() {
	incoming := make(chan common.Tuple, bufferLen)
	outgoing := make(chan common.Transaction, bufferLen)

	// create new wallet.
	w := util.NewWallet(common.Address{0})

	go client.HandleIncoming(w, incoming)
	go client.HandleOutgoing(w, outgoing)
	go util.LaunchClientConsole(w, outgoing)

	http.HandleFunc("/transaction", client.ReceiveSignature(incoming))
	http.ListenAndServe(":5555", nil)
}
