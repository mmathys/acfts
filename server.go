package main

import (
	"encoding/json"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"github.com/mmathys/acfts/util"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"sync"
)

var SignedUTXO sync.Map
var TxCounter *int32 = new(int32)
var BenchmarkMode bool = false

func handleSign(id *common.Identity) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		if BenchmarkMode {
			defer util.CountTx(TxCounter)
		}

		// parse the request
		var tx common.Transaction
		err := json.NewDecoder(req.Body).Decode(&tx)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, input := range tx.Inputs {
			SignedUTXO.Store(input.Id, input)
		}

		// Sign the request
		outputs, err := core.Sign(id.Key, tx.Outputs)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Form the response
		response := common.TransactionSignRes{Outputs: outputs}
		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func runServer(addr common.Address, benchmark bool) error {
	port := core.GetPort(addr)

	if !benchmark {
		log.Printf("initialized server; port = %d; benchmark = %t\n", port, benchmark)
	} else {
		BenchmarkMode = true
		go util.Ticker(TxCounter)
	}
	id := util.GetIdentity(addr)

	http.HandleFunc("/sign", handleSign(id))
	localAddr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(localAddr, nil)
	return nil
}

func main() {
	app := &cli.App{
		Name:  "ACFTS server",
		Usage: "Asynchronous Consensus-Free Transaction System server",
		Action: func(c *cli.Context) error {
			addr, err := util.ReadAddress(c.String("address"))
			if err != nil {
				log.Fatal(err)
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
