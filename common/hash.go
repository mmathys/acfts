package common

import (
	"crypto"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"hash"
)

func HashValueSprintf(value Value) []byte {
	d := crypto.SHA512.New()
	value.Signatures = nil                    // zero out signatures before hash
	d.Write([]byte(fmt.Sprintf("%v", value))) // this may be slow!
	return d.Sum(nil)
}

func HashValue(value Value) []byte {
	value.Signatures = nil // zero out signatures before hash

	d := crypto.SHA512.New()
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
			(*d).Write(signature.Address)
			(*d).Write(signature.Signature)
		}
	}
}

func HashTransactionSigRequest(req TransactionSigReq) []byte {
	req.Signature = EdDSASig{} // zero out signature before hash
	d := crypto.SHA512.New()

	for _, input := range req.Transaction.Inputs {
		writeValue(&d, &input)
	}

	for _, output := range req.Transaction.Outputs {
		writeValue(&d, &output)
	}

	hash := d.Sum(nil)
	return hash
}
