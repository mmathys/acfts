package crypto

import (
	"fmt"
	"github.com/mmathys/acfts/common"
	"golang.org/x/crypto/sha3"
)

func HashValue(value common.Value) []byte {
	d := sha3.New256()
	value.Signatures = nil                    // zero out signatures before hash
	d.Write([]byte(fmt.Sprintf("%v", value))) // this may be slow!
	return d.Sum(nil)
}

func HashTransactionSigRequest(req common.TransactionSigReq) []byte {
	d := sha3.New256()
	req.Signature = common.ECDSASig{}                // zero out signatures before hash
	d.Write([]byte(fmt.Sprintf("%v", req))) // this may be slow!
	return d.Sum(nil)
}