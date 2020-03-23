package client

import (
	"bytes"
	"github.com/mmathys/acfts/common"
)

func combineSignatures(res *[]common.TransactionSignRes) common.TransactionSignRes {
	baseRes := (*res)[0]

	for _, r := range (*res)[1:] {
		for _, rOutput := range r.Outputs {
			for i, baseOutput := range baseRes.Outputs {
				if bytes.Equal(baseOutput.Id, rOutput.Id) {
					baseRes.Outputs[i].Signatures = append(baseRes.Outputs[i].Signatures, rOutput.Signatures[0])
				}
			}
		}
	}

	return baseRes
}