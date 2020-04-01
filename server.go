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
	"runtime"
)

var TxCounter = new(int32)

func runServer(address common.Address, benchmark bool, adapter string, topology string) error {
	common.InitAddresses(topology)

	port := common.GetPort(address)
	server.SetAdapterMode(adapter)

	if !benchmark {
		log.Printf("initialized server; port = %d; benchmark = %t; adapter=%s\n", port, benchmark, adapter)
	} else {
		go util.Ticker(TxCounter)
	}

	if benchmark {
		runtime.SetBlockProfileRate(.1)
	}

	id := util.GetIdentity(address)
	server.Init(port, id, false, benchmark, TxCounter)

	return nil
}

func main() {
	app := &cli.App{
		Name:  "ACFTS server",
		Usage: "Asynchronous Consensus-Free Transaction System server",
		Action: func(c *cli.Context) error {
			adapter := "rest"
			if c.String("adapter") != "" {
				adapter = c.String("adapter")
			}

			addr, err := client.ReadAddress(c.String("address"))
			if err != nil {
				log.Fatal(err)
			}

			if addr == nil {
				log.Fatal("must define address")
			}

			benchmark := c.Bool("benchmark")

			runServer(addr, benchmark, adapter, c.String("topology"))
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Aliases:  []string{"a"},
				Usage:    "Set own address to `ADDRESS`. Format: e.g. 0x04",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "adapter",
				Usage:    "Set the adapter. Either \"rest\" or \"rpc\"",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "topology",
				Aliases:  []string{"t"},
				Usage:    "Path to the topology json file",
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
