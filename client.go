package main

import (
	"fmt"
	"github.com/mmathys/acfts/client"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/util"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
)

const bufferLen int = 255


func runClient(c *cli.Context) error {
	addr, err := util.ReadAddress(c.String("address"))
	if err != nil {
		log.Fatal(err)
	}

	port := c.Int("port")
	if port < 1025 || port > 65535 {
		log.Fatal("port must be between 1025 and 65535")
	}

	log.Printf("initialized client; addr = 0x%x, port = %d\n", addr, port)

	incoming := make(chan common.Tuple, bufferLen)
	outgoing := make(chan common.Transaction, bufferLen)

	// create new wallet.
	w := util.NewWallet(addr)

	go client.HandleIncoming(w, incoming)
	go client.HandleOutgoing(w, outgoing)
	go util.LaunchClientConsole(w, outgoing)

	http.HandleFunc("/transaction", client.ReceiveSignature(incoming))
	localAddr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(localAddr, nil)

	return nil
}

func main() {
	app := &cli.App{
		Name:   "ACFTS client",
		Usage:  "Asynchronous Consensus-Free Transaction System client",
		Action: runClient,
		Flags: []cli.Flag {
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Set client port to `PORT` for accepting forwarded transactions",
				Required: true,
			},
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
