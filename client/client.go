package main

import (
	"fmt"
	clientAdapter "github.com/mmathys/acfts/client/adapter"
	"github.com/mmathys/acfts/client/cli"
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	urfaveCli "github.com/urfave/cli/v2"
	"log"
	"os"
)

const bufferLen int = 255
var disableBatchVerification = false

func handleIncoming(w *common.Wallet, incoming chan common.Value) {
	for {
		t := <-incoming
		// verify
		err := common.VerifyValue(&t, !disableBatchVerification)
		if err != nil {
			panic(err)
		}
		core.AddUTXO(w, t)
	}
}

func runClient(c *urfaveCli.Context) error {
	addr, err := common.ReadAddress(c.String("address"))
	if err != nil {
		log.Fatal(err)
	}

	common.InitAddresses(c.String("topology"))

	port := common.GetClientPort(addr)
	balance := common.GetClientBalance(addr)
	if balance == 0 {
		newBalance := 100
		fmt.Printf("Warning: no balance set for client. Initializing balance to %d.\n", newBalance)
		balance = newBalance
	}
	
	disableBatchVerification = c.Bool("disable-batch")

	log.Printf("initialized client; addr = %x port = %d balance = %d", addr, port, balance)

	incoming := make(chan common.Value, bufferLen)

	w := common.NewWalletWithAmount(addr, balance)

	go handleIncoming(w, incoming)
	go cli.LaunchClientConsole(w)

	clientAdapter.Init(port, incoming)

	return nil
}

func main() {
	app := &urfaveCli.App{
		Name:   "ACFTS client",
		Usage:  "Asynchronous Consensus-Free Transaction System client",
		Action: runClient,
		Flags: []urfaveCli.Flag{
			&urfaveCli.StringFlag{
				Name:     "address",
				Aliases:  []string{"a"},
				Usage:    "Set own address to `ADDRESS`. Format: e.g. 0x04",
				Required: true,
			},
			&urfaveCli.BoolFlag{
				Name:     "disable-batch",
				Usage:    "Disable EdDSA batch signature verification",
				Required: false,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
