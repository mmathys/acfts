package common

import (
	"crypto"
	"encoding/binary"
	"encoding/gob"
	"hash"
)

func HashValue(mode int, value Value) []byte {
	value.Signatures = nil // zero out signatures before hash

	var d hash.Hash
	if mode == ModeEdDSA {
		d = crypto.SHA512.New()
	} else if mode == ModeBLS {
		d = crypto.SHA3_256.New()
	} else {
		panic("unsupported mode")
	}
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
			binary.Write(*d, binary.LittleEndian, signature.Mode)
		}
	}
}

func HashTransactionSigRequest(mode int, req TransactionSigReq) []byte {
	req.Signature = Signature{} // zero out signature before hash

	var d hash.Hash
	if mode == ModeEdDSA {
		d = crypto.SHA512.New()
	} else if mode == ModeBLS {
		d = crypto.SHA3_256.New()
	} else {
		panic("unsupported mode")
	}

	for _, input := range req.Transaction.Inputs {
		writeValue(&d, &input)
	}

	for _, output := range req.Transaction.Outputs {
		writeValue(&d, &output)
	}

	hash := d.Sum(nil)
	return hash
}
