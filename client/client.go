package main

import (
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"github.com/mmathys/acfts/wallet"
	"github.com/urfave/cli"
	"log"
	"os"
)

const bufferLen int = 255

func handleIncoming(w *common.Wallet, incoming chan common.Value) {
	for {
		t := <-incoming
		//fmt.Printf("got tuple %v\n", t)
		wallet.AddUTXO(w, t)
	}
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
	core.SetAdapterMode(adapter)

	port := common.GetPort(addr)

	log.Printf("initialized client; addr = 0x%x port = %d adapter=%s\n", addr, port, adapter)

	incoming := make(chan common.Value, bufferLen)

	w := util.NewWallet(addr)

	go handleIncoming(w, incoming)
	go core.LaunchClientConsole(w)

	core.Init(port, incoming)

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
