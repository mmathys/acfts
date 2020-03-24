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

	port := common.GetPort(addr)

	log.Printf("initialized client; addr = 0x%x, port = %d\n", addr, port)

	incoming := make(chan common.Value, bufferLen)

	w := util.NewWallet(addr)

	go client.HandleIncoming(w, incoming)
	go client.LaunchClientConsole(w)

	client.InitREST(port, incoming)

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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
