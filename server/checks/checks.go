package checks

import (
	"github.com/mmathys/acfts/common"
)

/**
Functions, which check the validity of incoming UTXOs before signing. Used by Server only.
*/

func CheckValidity(req *common.TransactionSigReq) error {
	tx := req.Transaction
	err := common.CheckFormat(&tx)
	if err != nil {
		return err
	}

	err = common.CheckConstraints(&tx)
	if err != nil {
		return err
	}

	err = checkInputSignatures(&tx)
	if err != nil {
		return err
	}

	err = checkRequestSignature(req)
	if err != nil {
		return err
	}

	return nil
}

/**
checks if input signatures are valid.
*/
func checkInputSignatures(tx *common.Transaction) error {
	for _, input := range tx.Inputs {
		err := common.VerifyValue(&input)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
Checks whether the request signature is valid
*/
func checkRequestSignature(req *common.TransactionSigReq) error {
	return common.VerifyTransactionSigRequest(req)
}
