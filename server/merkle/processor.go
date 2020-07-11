package merkle

import "github.com/mmathys/acfts/common"

func Processor(key *common.Key, chDispatch chan []*PoolMsg) {
	for {
		dispatchGroup := <-chDispatch

		// calculate the hashes of all outputs, of each request.
		var hashes [][]byte
		for _, msg := range dispatchGroup {
			for _, output := range msg.Req.Transaction.Outputs {
				hash := common.HashValue(key.Mode, output)
				hashes = append(hashes, hash)
			}
		}

		// sign ALL hashes in one go.
		sigs := key.SignMultipleMerkle(hashes)

		// fill out the response struct for each request and notify that we're done.
		merkleI := 0
		for _, msg := range dispatchGroup {
			outputs := msg.Req.Transaction.Outputs
			for i, _ := range outputs {
				outputs[i].Signatures = []common.Signature{*sigs[merkleI]}
				merkleI++
			}
			// "respond"
			*msg.Res = common.TransactionSignRes{Outputs: outputs}
			msg.WaitGroup.Done()
		}

		if merkleI != len(hashes) {
			panic("something went wrong in the loop")
		}
	}
}
