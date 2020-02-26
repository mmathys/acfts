package util

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"github.com/mmathys/acfts/wallet"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Show help
func help() {
	fmt.Println("send X 100\t\tSend 100 money to X")
	fmt.Println("utxo\t\t\tShow local UTXOs")
	fmt.Println("info\t\t\tShow client information")
	fmt.Println("clear\t\t\tClear console")
	fmt.Println("help\t\t\tShow this help section")
}


func send(s []string, outgoing chan common.Transaction) {
	if len(s) != 3 {
		fmt.Println("invalid format. Sample format: send X 100")
		return
	}

	alias := s[1]
	address := core.LookupAddress(alias)
	value, err := strconv.Atoi(s[2])
	if err != nil {
		fmt.Printf("%s is not a valid number.", s[2])
		return
	}

	t, err := wallet.PrepareTransaction(address, value)
	if err != nil {
		fmt.Println("failed to prepare transaction")
		return
	}

	outgoing <- t
}

func info() {
	fmt.Println("WIP")
}

func utxo() {
	for _, t := range wallet.UTXO {
		fmt.Println(t)
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "send", Description: "Send money to X"},
		{Text: "utxo", Description: "Show local UTXOs"},
		{Text: "info", Description: "Show client information"},
		{Text: "help", Description: "Show help"},
		{Text: "clear", Description: "Clear screen"},
	}
	x := d.FindStartOfPreviousWord()
	isFirstWord := x == 0
	if isFirstWord {
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	} else {
		return nil
	}
}

func LaunchClientConsole(outgoing chan common.Transaction) {
	fmt.Println("For help, enter \"help\"")
	for {
		t := prompt.Input("> ", completer)
		s := strings.Split(t, " ")
		cmd := s[0]
		switch cmd {
		case "help": help()
		case "send": send(s, outgoing)
		case "utxo": utxo()
		case "info": info()
		case "clear": clear()
		case "":
		default: fmt.Printf("Unknown command \"%s\"\n", cmd)
		}
	}
}