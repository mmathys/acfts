package main

import (
	"github.com/mmathys/acfts/common"
	serverAdapter "github.com/mmathys/acfts/server/adapter"
	"github.com/mmathys/acfts/server/util"
	"github.com/urfave/cli/v2"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"
)

var TxCounter = new(int32)
var SignedUTXO sync.Map

func runServer(address common.Address, instanceIndex int, benchmark bool, topology string, pprof bool) error {
	common.InitAddresses(topology)

	port := common.GetServerPort(address, instanceIndex)

	log.Println("initialized server")
	log.Printf("addr=%x, instance=%d, port=%d, benchmark = %t, pprof=%t\n", address, instanceIndex, port, benchmark, pprof)

	if benchmark {
		go util.Ticker(TxCounter)
	}

	if pprof {
		runtime.SetBlockProfileRate(1)
	}

	id := common.GetIdentity(address)
	serverAdapter.Init(port, id, false, benchmark, TxCounter, &SignedUTXO)

	return nil
}

func main() {
	app := &cli.App{
		Name:  "ACFTS server",
		Usage: "Asynchronous Consensus-Free Transaction System server",
		Action: func(c *cli.Context) error {
			instanceIndex := c.Int("instance") // if not set, value is 0

			addr, err := common.ReadAddress(c.String("address"))
			if err != nil {
				log.Fatal(err)
			}

			if addr == nil {
				log.Fatal("must define address")
			}

			benchmark := c.Bool("benchmark")
			pprof := c.Bool("pprof")

			runServer(addr, instanceIndex, benchmark, c.String("topology"), pprof)
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
				Name:     "instance",
				Aliases:  []string{"i"},
				Usage:    "Sets the zero-based instance index. This is used for load balancing/sharding. Default: 0",
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
			&cli.BoolFlag{
				Name:     "pprof",
				Usage:    "Enables pprof on default http server",
				Required: false,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
