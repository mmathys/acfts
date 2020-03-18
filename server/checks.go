package server

import (
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/crypto"
)

/**
Functions, which check the validity of incoming UTXOs before signing. Used by server only.
*/

func CheckValidity(id *common.Identity, tx *common.Transaction) error {
	err := common.CheckFormat(tx)
	if err != nil {
		return err
	}

	err = common.CheckConstraints(tx)
	if err != nil {
		return err
	}

	err = checkInputSignatures(id, tx)
	if err != nil {
		return err
	}

	err = checkUnspent(id, tx)
	if err != nil {
		return err
	}

	return nil
}

/**
checks if inputs signatures are valid.
*/
func checkInputSignatures(id *common.Identity, tx *common.Transaction) error {
	for _, input := range tx.Inputs {
		err := crypto.VerifyValue(id.Key, &input)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
checks whether transactions is unspent TODO
*/
func checkUnspent(id *common.Identity, tx *common.Transaction) error {
	return nil
}
