package wallet

import (
	"errors"
	"github.com/mmathys/acfts/common"
	"math/rand"
)

func PrepareTransaction(w *common.Wallet, addr common.Address, val int) (common.Transaction, error) {
	// Linear Scan through UTXOs
	var inputs []common.Value
	current := 0
	w.UTXO.Range(func(_ interface{}, value interface{}) bool {
		v := value.(common.Value)
		if current < val {
			inputs = append(inputs, v)
			current += v.Amount
			return true
		} else {
			return false
		}
	})

	if current < val {
		return common.Transaction{}, errors.New("not enough funds")
	}

	var outputs []common.Value

	// add remaining funds to inputs
	if current > val {
		remaining := current - val
		outputs = append(outputs, common.Value{Address: w.Address, Amount: remaining, Id: rand.Int()})
	}

	// add counterpart
	outputs = append(outputs, common.Value{Address: addr, Amount: val, Id: rand.Int()})

	t := common.Transaction{Inputs: inputs, Outputs: outputs}
	return t, nil
}

func RemoveUTXO(wallet *common.Wallet, t common.Value) {
	wallet.UTXO.Delete(t.Id)
}

func RemoveUTXOMultiple(wallet *common.Wallet, ts *[]common.Value) {
	for _, t := range *ts {
		wallet.UTXO.Delete(t.Id)
	}
}

func AddUTXO(wallet *common.Wallet, t common.Value) {
	wallet.UTXO.Store(t.Id, t)
}

func AddUTXOMultiple(wallet *common.Wallet, ts *[]common.Value) {
	for _, t := range *ts {
		wallet.UTXO.Store(t.Id, t)
	}
}
