package main

import (
	"encoding/json"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"sync"
)

//var SignedUTXO = map[int]common.Tuple{}
var SignedUTXO sync.Map

func handleSign(w http.ResponseWriter, req *http.Request) {
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

	fmt.Println(SignedUTXO)

	// Sign the request
	outputs, err := core.Sign(&tx.Inputs, &tx.Outputs)
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

func runServer(c *cli.Context) error {
	port := c.Int("port")
	if port < 1025 || port > 65535 {
		log.Fatal("port must be between 1025 and 65535")
	}

	log.Printf("initialized server; port = %d\n", port)

	http.HandleFunc("/sign", handleSign)
	localAddr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(localAddr, nil)
	return nil
}

func main() {
	app := &cli.App{
		Name:   "ACFTS server",
		Usage:  "Asynchronous Consensus-Free Transaction System server",
		Action: runServer,
		Flags: []cli.Flag {
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Set server port to `PORT` for signing endpoint",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
