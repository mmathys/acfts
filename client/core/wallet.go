package core

import (
	"errors"
	"github.com/mmathys/acfts/common"
)

// prepare transaction mutex

func getIndex(id common.Identifier) [common.IdentifierLength]byte {
	if len(id) != common.IdentifierLength {
		panic("identifier length mismatch")
	}
	index := [common.IdentifierLength]byte{}
	copy(index[:], id[:common.IdentifierLength])
	return index
}

func PrepareTransaction(w *common.Wallet, target common.Address, val int) (common.Transaction, error) {
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

	addressOwn := w.Address
	counterpart := common.GetIdentity(target)
	addressCounterpart := counterpart.Address

	var outputs []common.Value

	// add remaining fund to output
	if current > val {
		remaining := current - val
		outputs = append(outputs, common.Value{Address: addressOwn, Amount: remaining, Id: common.RandomIdentifier()})
	}

	// add counterpart
	outputs = append(outputs, common.Value{Address: addressCounterpart, Amount: val, Id: common.RandomIdentifier()})

	t := common.Transaction{Inputs: inputs, Outputs: outputs}
	return t, nil
}

func RemoveUTXO(wallet *common.Wallet, t common.Value) {
	wallet.UTXO.Delete(getIndex(t.Id))
}

func RemoveUTXOMultiple(wallet *common.Wallet, ts *[]common.Value) {
	for _, t := range *ts {
		wallet.UTXO.Delete(getIndex(t.Id))
	}
}

func AddUTXO(wallet *common.Wallet, t common.Value) {
	wallet.UTXO.Store(getIndex(t.Id), t)
}

func AddUTXOMultiple(wallet *common.Wallet, ts *[]common.Value) {
	for _, t := range *ts {
		wallet.UTXO.Store(getIndex(t.Id), t)
	}
}
