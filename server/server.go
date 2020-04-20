package main

import (
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/server/rpc"
	"github.com/mmathys/acfts/util"
	"github.com/urfave/cli/v2"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"
)


var TxCounter = new(int32)
type Adapter interface {
	Init(port int, id *common.Identity, debug bool, benchmark bool, TxCounter *int32, SignedUTXO *sync.Map)
}

var rpcAdapter = &rpc.RPCAdapter{}
var currentAdapter Adapter = rpcAdapter

var SignedUTXO sync.Map

func SetAdapterMode(mode string) {
	if mode == "rest" {
		panic("rest is not supported anymore")
	} else if mode == "rpc" {
		currentAdapter = rpcAdapter
	} else {
		log.Fatalf("unrecognized mode %s", mode)
	}
}

func Init(port int, id *common.Identity, debug bool, benchmark bool, TxCounter *int32) {
	currentAdapter.Init(port, id, debug, benchmark, TxCounter, &SignedUTXO)
}

func runServer(address common.Address, benchmark bool, adapter string, topology string, pprof bool) error {
	common.InitAddresses(topology)

	port := common.GetPort(address)
	SetAdapterMode(adapter)

	log.Printf("initialized server; port = %d; benchmark = %t; adapter=%s; pprof=%t;\n", port, benchmark, adapter, pprof)

	if benchmark {
		go util.Ticker(TxCounter)
	}

	if pprof {
		runtime.SetBlockProfileRate(1)
	}

	id := util.GetIdentity(address)
	Init(port, id, true, benchmark, TxCounter)

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

			addr, err := util.ReadAddress(c.String("address"))
			if err != nil {
				log.Fatal(err)
			}

			if addr == nil {
				log.Fatal("must define address")
			}

			benchmark := c.Bool("benchmark")
			pprof := c.Bool("pprof")

			runServer(addr, benchmark, adapter, c.String("topology"), pprof)
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
