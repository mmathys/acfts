package wallet

import (
	"errors"
	"github.com/mmathys/acfts/common"
	"math/rand"
)

func PrepareTransaction(w *common.Wallet, addr common.Address, val int) (common.Transaction, error) {
	// Linear Scan through UTXOs
	var inputs []common.Tuple
	current := 0
	for _, tx := range w.UTXO {
		if current < val {
			inputs = append(inputs, tx)
			current += tx.Value
		}
	}

	if current < val {
		return common.Transaction{}, errors.New("not enough funds")
	}

	var outputs []common.Tuple

	// add remaining funds to inputs
	if current > val {
		remaining := current - val
		outputs = append(outputs, common.Tuple{Address: w.Address, Value: remaining, Id: rand.Int()})
	}

	// add counterpart
	outputs = append(outputs, common.Tuple{Address: addr, Value: val, Id: rand.Int()})

	t := common.Transaction{Inputs: inputs, Outputs: outputs}
	return t, nil
}

func RemoveUTXO(wallet *common.Wallet, ts *[]common.Tuple) {
	for _, t := range *ts {
		delete(wallet.UTXO, t.Id)
	}
}

func AddUTXO(wallet *common.Wallet, ts *[]common.Tuple) {
	for _, t := range *ts {
		wallet.UTXO[t.Id] = t
	}
}
