package util

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/wallet"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Show help
func help() {
	fmt.Println("send 0x00 100\t\tSend 100 money to 0x00")
	fmt.Println("utxo\t\t\tShow local UTXOs")
	fmt.Println("info\t\t\tShow client information")
	fmt.Println("set addr 0x03\t\tSet own address to 0x03 (for testing)")
	fmt.Println("clear\t\t\tClear console")
	fmt.Println("help\t\t\tShow this help section")
}

func ReadAddress(s string) (common.Address, error) {
	split := strings.Split(s, "0x")
	if len(split) != 2 {
		return common.Address{}, errors.New("hex should look like 0x04\n")
	}

	addrInput, err := hex.DecodeString(split[1])
	if err != nil {
		return common.Address{}, errors.New("could not decode hex\n")
	}

	var address common.Address
	copy(address[:], addrInput)
	return address, nil
}

func send(w *common.Wallet, s []string, outgoing chan common.Transaction) {
	if len(s) != 3 {
		fmt.Println("invalid format. Sample format: send X 100")
		return
	}

	address, err := ReadAddress(s[1])
	if err != nil {
		fmt.Print("Error: ", err.Error())
		return
	}

	value, err := strconv.Atoi(s[2])
	if err != nil {
		fmt.Printf("%s is not a valid number.", s[2])
		return
	}

	t, err := wallet.PrepareTransaction(w, address, value)
	if err != nil {
		fmt.Println("failed to prepare transaction")
		return
	}

	outgoing <- t
}

func info(w *common.Wallet) {
	fmt.Printf("Address:\t0x%x\t%d\n", w.Address, w.Address)
}

func utxo(w *common.Wallet) {
	for _, t := range w.UTXO {
		fmt.Println(t)
	}
}

func balance(w *common.Wallet) {
	balance := 0
	for _, t := range w.UTXO {
		balance += t.Value
	}
	fmt.Printf("Balance: %d\n", balance)
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func setAddr(w *common.Wallet, s []string) {
	if len(s) != 3 || s[1] != "addr" {
		fmt.Println("invalid format. Sample format: set addr 0x03")
		return
	}

	addr, err := ReadAddress(s[2])
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	w.Address = addr
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "send", Description: "Send money to X"},
		{Text: "utxo", Description: "Show local UTXOs"},
		{Text: "balance", Description: "Show balance"},
		{Text: "info", Description: "Show client information"},
		{Text: "set", Description: "Set property"},
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

func LaunchClientConsole(w *common.Wallet, outgoing chan common.Transaction) {
	fmt.Println("For help, enter \"help\"")
	for {
		t := prompt.Input("> ", completer)
		s := strings.Split(t, " ")
		cmd := s[0]
		switch cmd {
		case "help":
			help()
		case "send":
			send(w, s, outgoing)
		case "utxo":
			utxo(w)
		case "balance":
			balance(w)
		case "info":
			info(w)
		case "clear":
			clear()
		case "set":
			setAddr(w, s)
		case "":
		default:
			fmt.Printf("Unknown command \"%s\"\n", cmd)
		}
	}
}
