package cli

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/mmathys/acfts/client/core"
	"github.com/mmathys/acfts/common"
	"github.com/olekukonko/tablewriter"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Show help
func help() {
	fmt.Println("send <address> 100\tSend 100 credits to <address>. Must be encoded in hexadecimal.")
	fmt.Println("utxo\t\t\tShow local UTXOs")
	fmt.Println("balance\t\t\tShow balance")
	fmt.Println("info\t\t\tShow client information")
	fmt.Println("clear\t\t\tClear console")
	fmt.Println("help\t\t\tShow this help section")
}

func send(w *common.Wallet, s []string) {
	if len(s) != 3 {
		fmt.Println("invalid format. Sample format: send X 100")
		return
	}

	address, err := common.ReadAddress(s[1])
	if err != nil {
		fmt.Print("Error: ", err.Error())
		return
	}

	value, err := strconv.Atoi(s[2])
	if err != nil {
		fmt.Printf("%s is not a valid number.", s[2])
		return
	}

	t, err := core.PrepareTransaction(w, address, value)
	if err != nil {
		fmt.Println("failed to prepare transaction")
		return
	}

	core.DoTransaction(w, t, true)
}

func info(w *common.Wallet) {
	addr := w.GetAddress()
	key := common.GetKey(addr)
	keySerialized := w.SerializePrivateKey()
	net, _ := common.GetClientNetworkAddress(addr)
	fmt.Printf("Address (public key):\t%x\n", addr)
	fmt.Printf("Private Key:\t\t%x\n", keySerialized)
	modeReadable := "unrecognized"
	if key.Mode == common.ModeEdDSA {
		modeReadable = "EdDSA"
	} else if key.Mode == common.ModeBLS {
		modeReadable = "BLS"
	} else if key.Mode == common.ModeMerkle {
		modeReadable = "Merkle"
	}
	fmt.Printf("Key Mode:\t\t%v\n", modeReadable)
	fmt.Printf("Network:\t\t%s\n", net)
}

func utxo(w *common.Wallet) {
	table := tablewriter.NewWriter(os.Stdout)
	sum := 0
	table.SetHeader([]string{"address", "amount", "id", "sig count", "sig mode"})
	w.UTXO.Range(func(_ interface{}, value interface{}) bool {
		v := value.(common.Value)
		trimmedAddr := fmt.Sprintf("%x", v.Address)[:10]
		trimmedId := fmt.Sprintf("%x", v.Id)[:10]
		sigMode := "unknown"
		if len(v.Signatures) > 0 {
			if v.Signatures[0].Mode == common.ModeEdDSA {
				sigMode = "multisig (EdDSA)"
			} else if v.Signatures[0].Mode == common.ModeBLS {
				sigMode = "threshold (BLS)"
			} else if v.Signatures[0].Mode == common.ModeMerkle {
				sigMode = "merkle (EdDSA)"
			}
		}

		table.Append([]string{
			fmt.Sprintf("%s...", trimmedAddr),
			fmt.Sprintf("%d", v.Amount),
			fmt.Sprintf("%s...", trimmedId),
			fmt.Sprintf("%d", len(v.Signatures)),
			sigMode,
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

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "send", Description: "Send credits to X"},
		{Text: "utxo", Description: "Show local UTXOs"},
		{Text: "balance", Description: "Show balance"},
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

func LaunchClientConsole(w *common.Wallet) {
	fmt.Println("For help, enter \"help\"")
	for {
		t := prompt.Input("> ", completer)
		s := strings.Split(t, " ")
		cmd := s[0]
		switch cmd {
		case "help":
			help()
		case "send":
			send(w, s)
		case "utxo":
			utxo(w)
		case "balance":
			balance(w)
		case "info":
			info(w)
		case "clear":
			clear()
		case "":
		default:
			fmt.Printf("Unknown command \"%s\"\n", cmd)
		}
	}
}
