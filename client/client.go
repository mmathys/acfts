package main

import (
	"bytes"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"github.com/urfave/cli"
	"log"
	"os"
	"sync"
)

const bufferLen int = 255

func handleIncoming(w *common.Wallet, incoming chan common.Value) {
	for {
		t := <-incoming
		//fmt.Printf("got tuple %v\n", t)
		wallet.AddUTXO(w, t)
	}
}

func DoTransaction(w *common.Wallet, t common.Transaction, forward bool) {
	res, err := SignTransaction(w, t)
	if err != nil {
		fmt.Println("failed to sign transaction")
		return
	}

	// own UTXOs, (is spent at this point)
	wallet.RemoveUTXOMultiple(w, &t.Inputs)

	sig := combineSignatures(res)

	// add own outputs
	var ownOutputs []common.Value
	for _, v := range sig.Outputs {
		if bytes.Equal(v.Address, w.Identity.Address) {
			ownOutputs = append(ownOutputs, v)
		} else if forward {
			go ForwardValue(v)
		}
	}

	wallet.AddUTXOMultiple(w, &ownOutputs)
}

// TODO Only wait for Math.ceil(2/3 * n) of n servers!
func SignTransaction(w *common.Wallet, t common.Transaction) (*[]common.TransactionSignRes, error) {
	n := len(common.GetServers())

	sigs := make(chan common.TransactionSignRes, n)

	var wg sync.WaitGroup

	for _, server := range common.GetServers() {
		wg.Add(1)
		go RequestSignature(server, w.Identity, t, &wg, &sigs)
	}

	wg.Wait()

	// TODO validate and store sigs
	var res []common.TransactionSignRes
	for i := 0; i < n; i++ {
		sig := <-sigs
		res = append(res, sig)
	}

	return &res, nil
}

func runClient(c *cli.Context) error {
	addr, err := util.ReadAddress(c.String("address"))
	if err != nil {
		log.Fatal(err)
	}

	common.InitAddresses(c.String("topology"))

	adapter := "rest"
	if c.String("adapter") != "" {
		adapter = c.String("adapter")
	}
	SetAdapterMode(adapter)

	port := common.GetPort(addr)

	log.Printf("initialized client; addr = 0x%x port = %d adapter=%s\n", addr, port, adapter)

	incoming := make(chan common.Value, bufferLen)

	w := util.NewWallet(addr)

	go handleIncoming(w, incoming)
	go launchClientConsole(w)

	Init(port, incoming)

	return nil
}

func main() {
	app := &cli.App{
		Name:   "ACFTS client",
		Usage:  "Asynchronous Consensus-Free Transaction System client",
		Action: runClient,
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Usage:   "Set own address to `ADDRESS`. Format: e.g. 0x04",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "topology",
				Aliases:  []string{"t"},
				Usage:    "Path to the topology json file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "adapter",
				Usage:    "Set the adapter. Either \"rest\" or \"rpc\"",
				Required: false,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
