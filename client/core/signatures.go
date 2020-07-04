package core

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/mmathys/acfts/common"
)

func combineSignatures(res *[]common.TransactionSignRes) common.TransactionSignRes {
	baseRes := (*res)[0]
	if len(baseRes.Outputs) == 0 {
		panic("encountered response with zero outputs")
	}
	if len(baseRes.Outputs[0].Signatures) == 0 {
		panic("encountered response with zero signatures")
	}

	mode := baseRes.Outputs[0].Signatures[0].Mode

	for _, r := range (*res)[1:] {
		for _, rOutput := range r.Outputs {
			for i, baseOutput := range baseRes.Outputs {
				if baseOutput.Id == rOutput.Id {
					baseRes.Outputs[i].Signatures = append(baseRes.Outputs[i].Signatures, rOutput.Signatures[0])
				}
			}
		}
	}

	if mode == common.ModeEdDSA || mode == common.ModeMerkle {
		// return response unmerged
		return baseRes
	} else if mode == common.ModeBLS {
		// For each output...
		for i := range baseRes.Outputs {
			// ... apply BLS threshold
			var sigs []bls.Sign
			var ids []bls.ID
			for _, sig := range baseRes.Outputs[i].Signatures {
				var partialSig bls.Sign
				partialSig.Deserialize(sig.Signature)
				sigs = append(sigs, partialSig)
				var id bls.ID
				id.Deserialize(sig.BLSID)
				ids = append(ids, id)
			}
			var masterSig bls.Sign
			err := masterSig.Recover(sigs, ids)
			if err != nil {
				panic(err)
			}
			masterPub := common.GetBLSMasterPublicKey()
			blsSig := common.Signature{
				Address:   masterPub.Serialize(),
				Signature: masterSig.Serialize(),
				Mode:      mode,
			}
			baseRes.Outputs[i].Signatures = []common.Signature{blsSig}
		}
		return baseRes
	} else {
		panic("unsupported mode")
	}
}
