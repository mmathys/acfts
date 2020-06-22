package bls

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/mmathys/acfts/common"
	"testing"
)

func TestPrintGeneratedBLSKey(t *testing.T) {
	bls.Init(bls.BLS12_381)
	bls.SetETHmode(bls.EthModeDraft07)
	for i := 0; i < 64; i++ {
		//var key bls.SecretKey
		//key.SetByCSPRNG()
		key := common.GenerateKey(common.ModeBLS)
		pub := key.GetAddress()
		fmt.Printf("{\"%x\",\"%x\"},\n", pub.Serialize(), key.Serialize())
	}
}
