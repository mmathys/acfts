package common

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"golang.org/x/crypto/sha3"
	"hash"
)

func HashValueSprintf(value Value) []byte {
	d := sha3.New256()
	value.Signatures = nil                    // zero out signatures before hash
	d.Write([]byte(fmt.Sprintf("%v", value))) // this may be slow!
	return d.Sum(nil)
}

func HashValue(value Value) []byte {
	value.Signatures = nil // zero out signatures before hash

	d := sha3.New256()
	enc := gob.NewEncoder(d)
	enc.Encode(value)
	return d.Sum(nil)
}

func writeValue(d *hash.Hash, value *Value) {
	(*d).Write(value.Address)
	binary.Write(*d, binary.LittleEndian, value.Amount)
	(*d).Write(value.Id[:])
	if value.Signatures != nil {
		for _, signature := range value.Signatures {
			(*d).Write(signature)
		}
	}
}

func HashTransactionSigRequest(req TransactionSigReq) []byte {
	req.Signature = []byte{} // zero out signatures before hash
	d := sha3.New256()

	for _, input := range req.Transaction.Inputs {
		writeValue(&d, &input)
	}

	for _, output := range req.Transaction.Outputs {
		writeValue(&d, &output)
	}

	hash := d.Sum(nil)
	return hash
}
