package wallet

import (
	"errors"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/core"
	"math/rand"
	"time"
)

// for testing, each user will have 100 money
// TODO: this is pseudo-deterministic. Used for testing only.
var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)
var key = r1.Int()
var UTXO = map[int]common.Tuple {
	key: {[]byte{0},100, key},
}

func PrepareTransaction(addr common.Address, val int) (common.Transaction, error){

	// Linear Scan through UTXOs
	var inputs []common.Tuple
	current := 0
	for _, tx := range UTXO {
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
		outputs = append(outputs, common.Tuple{core.GetOwnAddress(), remaining, rand.Int()})
	}

	// add counterpart
	outputs = append(outputs, common.Tuple{addr, val, 0})

	t := common.Transaction{inputs, outputs}
	return t, nil
}

func RemoveUTXO(ts *[]common.Tuple) {
	for _, t := range *ts {
		delete(UTXO, t.Id)
	}
}

func AddUTXO(ts *[]common.Tuple) {
	for _, t := range *ts {
		UTXO[t.Id] = t
	}
}