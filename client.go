package main

import (
	"github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
)

const bufferLen int = 255

func main() {
	incoming := make(chan common.Transaction, bufferLen)
	outgoing := make(chan common.Transaction, bufferLen)

	// empty address for testing
	w := util.NewWallet(common.Address{0})

	go client.HandleIncomingTransactions(w, incoming)
	go client.HandleOutgoingTransactions(w, outgoing)

	util.LaunchClientConsole(w, outgoing)
}
