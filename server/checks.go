package server

import (
	"github.com/mmathys/acfts/common"
)

/**
Functions, which check the validity of incoming UTXOs before signing. Used by Server only.
*/

func CheckValidity(id *common.Identity, req *common.TransactionSigReq) error {
	tx := req.Transaction
	err := common.CheckFormat(&tx)
	if err != nil {
		return err
	}

	err = common.CheckConstraints(&tx)
	if err != nil {
		return err
	}

	err = checkRequestSignature(id, req)
	if err != nil {
		return err
	}

	err = checkInputSignatures(&tx)
	if err != nil {
		return err
	}

	err = checkUnspent(id, &tx)
	if err != nil {
		return err
	}

	return nil
}

/**
Checks whether the request signature is valid; i.e. the public key of the signature matches with the public key of the inputs
*/
func checkRequestSignature(id *common.Identity, req *common.TransactionSigReq) error {
	return nil
}

/**
checks if inputs signatures are valid.
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
checks whether transactions is unspent TODO
*/
func checkUnspent(id *common.Identity, tx *common.Transaction) error {
	return nil
}
