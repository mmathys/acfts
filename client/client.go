package main

import (
	clientAdapter "github.com/mmathys/acfts/client/adapter"
	"github.com/mmathys/acfts/client/cli"
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	urfaveCli "github.com/urfave/cli/v2"
	"log"
	"os"
)

const bufferLen int = 255

func handleIncoming(w *common.Wallet, incoming chan common.Value) {
	for {
		t := <-incoming
		// verify
		err := common.VerifyValue(&t)
		if err != nil {
			panic(err)
		}
		core.AddUTXO(w, t)
	}
}

func runClient(c *urfaveCli.Context) error {
	addr, err := common.ReadAddress(c.String("address"))
	if err != nil {
		log.Fatal(err)
	}

	common.InitAddresses(c.String("topology"))

	port := common.GetClientPort(addr)

	log.Printf("initialized client; addr = %x port = %d", addr, port)

	incoming := make(chan common.Value, bufferLen)

	w := common.NewWallet(addr)

	go handleIncoming(w, incoming)
	go cli.LaunchClientConsole(w)

	clientAdapter.Init(port, incoming)

	return nil
}

func main() {
	app := &urfaveCli.App{
		Name:   "ACFTS client",
		Usage:  "Asynchronous Consensus-Free Transaction System client",
		Action: runClient,
		Flags: []urfaveCli.Flag{
			&urfaveCli.StringFlag{
				Name:     "address",
				Aliases:  []string{"a"},
				Usage:    "Set own address to `ADDRESS`. Format: e.g. 0x04",
				Required: true,
			},
			&urfaveCli.StringFlag{
				Name:     "topology",
				Aliases:  []string{"t"},
				Usage:    "Path to the topology json file",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
