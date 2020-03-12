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

func handleSign(id *common.Identity) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// parse the request
		var tx common.Transaction
		err := json.NewDecoder(req.Body).Decode(&tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, input := range tx.Inputs {
			SignedUTXO.Store(input.Id, input)
		}

		// Sign the request
		outputs, err := core.Sign(id.Key, tx.Outputs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Form the response
		response := common.TransactionSignRes{Outputs: outputs}
		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func runServer(addr common.Address) error {
	port := core.GetPort(addr)

	log.Printf("initialized server; port = %d\n", port)
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
			runServer(addr)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Aliases:  []string{"a"},
				Usage:    "Set own address to `ADDRESS`. Format: e.g. 0x04",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
