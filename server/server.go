package main

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/mmathys/acfts/common"
	serverAdapter "github.com/mmathys/acfts/server/adapter"
	"github.com/mmathys/acfts/server/store"
	"github.com/mmathys/acfts/server/util"
	"github.com/urfave/cli/v2"
	"log"
	//_ "net/http/pprof"
	"os"
	"runtime"
)

var TxCounter = new(int32)
var UTXOMap store.UTXOMap

type serverOpt struct {
	address common.Address
	instanceIndex int
	benchmark bool
	topology string
	pprof bool
	disableBatch bool
	mapType int
	scheme int
}

func runServer(opt serverOpt) error {
	bls.Init(bls.BLS12_381)
	bls.SetETHmode(bls.EthModeDraft07)

	common.InitAddresses(opt.topology)

	port := common.GetServerPort(opt.address, opt.instanceIndex)

	mapTypeReadable := "unrecognized"
	if opt.mapType == store.TypeSyncMap {
		mapTypeReadable = "syncMap"
	} else if opt.mapType == store.TypeInsertOnly {
		mapTypeReadable = "insertOnly"
	}

	UTXOMap.SetType(opt.mapType)

	fmt.Println("initializing server with:")
	fmt.Printf("addr=%x, instance=%d, port=%d, benchmark=%t, pprof=%t, batchVerification=%t, mapType=%s\n",
		opt.address, opt.instanceIndex, port, opt.benchmark, opt.pprof, !opt.disableBatch, mapTypeReadable)

	if opt.benchmark {
		go util.Ticker(TxCounter)
	}

	if opt.pprof {
		runtime.SetBlockProfileRate(1)
	}

	key := common.GetKey(opt.address)
	serverAdapter.Init(serverAdapter.AdapterOpt{
		Port:              port,
		Key:               key,
		NoSigning:         false,
		Benchmark:         opt.benchmark,
		TxCounter:         TxCounter,
		UTXOMap:           &UTXOMap,
		BatchVerification: !opt.disableBatch,
	})

	return nil
}

func main() {
	app := &cli.App{
		Name:  "ACFTS server",
		Usage: "Asynchronous Consensus-Free Transaction System server",
		Action: func(c *cli.Context) error {
			addr, err := common.ReadAddress(c.String("address"))
			if err != nil {
				log.Fatal(err)
			}

			if addr == nil {
				log.Fatalf("must define address")
			}

			mapType := c.String("map-type")
			mapTypeInt := -1

			if mapType == "syncMap" {
				mapTypeInt = store.TypeSyncMap
			} else if mapType == "insertOnly" {
				mapTypeInt = store.TypeInsertOnly
			} else {
				log.Panicf("--map-type must be either 'syncMap' or 'insertOnly' (got %s)", mapType)
			}

			runServer(serverOpt{
				address:       addr,
				instanceIndex: c.Int("instance"),
				benchmark:     c.Bool("benchmark"),
				topology:      c.String("topology"),
				pprof:         c.Bool("pprof"),
				disableBatch:  c.Bool("disable-batch"),
				mapType:       mapTypeInt,
			})
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Aliases:  []string{"a"},
				Usage:    "Set own address to `ADDRESS`. Format: e.g. 0x04",
				Required: true,
			},
			&cli.IntFlag{
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
			&cli.BoolFlag{
				Name:     "disable-batch",
				Usage:    "Disable EdDSA batch signature verification",
				Required: false,
			},
			&cli.StringFlag{
				Name:		"map-type",
				Value:		store.DefaultMapTypeString,
				Usage:    	"Sets the internal map type. 'syncMap' or 'insertOnly'.",
				Required:	false,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
