package main

import (
	"github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"github.com/urfave/cli"
	"log"
	"os"
)

const bufferLen int = 255


func runClient(c *cli.Context) error {
	addr, err := client.ReadAddress(c.String("address"))
	if err != nil {
		log.Fatal(err)
	}

	common.InitAddresses(c.String("topology"))

	adapter := "rest"
	if c.String("adapter") != "" {
		adapter = c.String("adapter")
	}
	client.SetAdapterMode(adapter)

	port := common.GetPort(addr)

	log.Printf("initialized client; addr = 0x%x port = %d adapter=%s\n", addr, port, adapter)

	incoming := make(chan common.Value, bufferLen)

	w := util.NewWallet(addr)

	go client.HandleIncoming(w, incoming)
	go client.LaunchClientConsole(w)

	client.Init(port, incoming)

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
