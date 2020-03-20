package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"github.com/mmathys/acfts/wallet"
	"github.com/olekukonko/tablewriter"
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

func ReadAlias(s string) (common.Alias, error) {
	split := strings.Split(s, "0x")
	if len(split) != 2 {
		return common.Alias{}, errors.New("hex should look like 0x04\n")
	}

	addrInput, err := hex.DecodeString(split[1])
	if err != nil {
		return common.Alias{}, errors.New("could not decode hex\n")
	}

	var alias common.Alias
	copy(alias[:], addrInput)
	return alias, nil
}

func send(w *common.Wallet, s []string, outgoing chan common.Transaction) {
	if len(s) != 3 {
		fmt.Println("invalid format. Sample format: send X 100")
		return
	}

	address, err := ReadAlias(s[1])
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
	fmt.Printf("Alias:\t\t0x%x\n", w.Alias)

	net, _ := core.GetNetworkAddress(w.Alias)
	fmt.Printf("Network:\t%s\n", net)

	fmt.Printf("Private Key:\t0x%x\n", crypto.FromECDSA(w.Key))
	fmt.Printf("Public Key:\t0x%x\n", crypto.FromECDSAPub(&w.Key.PublicKey))
}

func utxo(w *common.Wallet) {
	table := tablewriter.NewWriter(os.Stdout)
	sum := 0
	table.SetHeader([]string{"address", "amount", "id", "sig count"})
	w.UTXO.Range(func(_ interface{}, value interface{}) bool {
		v := value.(common.Value)
		trimmedAddr := fmt.Sprintf("0x%x", v.Address)[:10]
		table.Append([]string{
			fmt.Sprintf("%s...", trimmedAddr),
			fmt.Sprintf("%d", v.Amount),
			fmt.Sprintf("%d", v.Id),
			fmt.Sprintf("%d", len(v.Signatures)),
		})
		sum += v.Amount
		return true
	})
	table.SetBorder(false)
	table.Render()
}

func balance(w *common.Wallet) {
	balance := 0
	w.UTXO.Range(func(_ interface{}, value interface{}) bool {
		v := value.(common.Value)
		balance += v.Amount
		return true
	})
	fmt.Printf("Balance: %d\n", balance)
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func setAlias(w *common.Wallet, s []string) {
	if len(s) != 3 || s[1] != "alias" {
		fmt.Println("invalid format. Sample format: set alias 0x03")
		return
	}

	alias, err := ReadAlias(s[2])
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	w.Alias = alias
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
			setAlias(w, s)
		case "":
		default:
			fmt.Printf("Unknown command \"%s\"\n", cmd)
		}
	}
}