package common

import (
	"encoding/gob"
	"fmt"
	"golang.org/x/crypto/sha3"
)

func HashValueSprintf(value Value) []byte {
	d := sha3.New256()
	value.Signatures = nil                    // zero out signatures before hash
	d.Write([]byte(fmt.Sprintf("%v", value))) // this may be slow!
	return d.Sum(nil)
}

func HashValue(value Value) []byte {
	value.Signatures = nil                    		// zero out signatures before hash

	d := sha3.New256()
	enc := gob.NewEncoder(d)
	enc.Encode(value)
	return d.Sum(nil)
}

func HashTransactionSigRequestSprintf(req TransactionSigReq) []byte {
	d := sha3.New256()
	req.Signature = ECDSASig{}              // zero out signatures before hash
	d.Write([]byte(fmt.Sprintf("%v", req))) // this may be slow!
	return d.Sum(nil)
}

func HashTransactionSigRequest(req TransactionSigReq) []byte {
	d := sha3.New256()
	enc := gob.NewEncoder(d)
	req.Signature = ECDSASig{}              // zero out signatures before hash
	enc.Encode(req)
	return d.Sum(nil)
}