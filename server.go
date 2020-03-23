package main

import (
	"github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server"
	"github.com/mmathys/acfts/util"
	"github.com/urfave/cli"
	"log"
	_ "net/http/pprof"
	"os"
	"sync"
)

var SignedUTXO sync.Map
var TxCounter = new(int32)

func runServer(address common.Address, benchmark bool) error {
	port := common.GetPort(address)

	if !benchmark {
		log.Printf("initialized server; port = %d; benchmark = %t\n", port, benchmark)
	} else {
		go util.Ticker(TxCounter)
	}

	id := util.GetIdentity(address)
	server.InitREST(port, id, false, benchmark, &SignedUTXO, TxCounter)

	return nil
}

func main() {
	app := &cli.App{
		Name:  "ACFTS server",
		Usage: "Asynchronous Consensus-Free Transaction System server",
		Action: func(c *cli.Context) error {
			addr, err := client.ReadAddress(c.String("address"))
			if err != nil {
				log.Fatal(err)
			}

			if addr == nil {
				log.Fatal("must define address")
			}

			benchmark := c.Bool("benchmark")

			runServer(addr, benchmark)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Aliases:  []string{"a"},
				Usage:    "Set own address to `ADDRESS`. Format: e.g. 0x04",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "benchmark",
				Aliases:  []string{"b"},
				Usage:    "Enables benchmark mode. If set, then outputs number of tx/s to stdout, separated by a newline.",
				Required: false,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
