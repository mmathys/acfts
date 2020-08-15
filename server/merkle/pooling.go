package merkle

import (
	"fmt"
	"github.com/mmathys/acfts/common"
	"sync"
)

type PoolMsg struct {
	Req       *common.TransactionSigReq
	Res       *common.TransactionSignRes
	WaitGroup *sync.WaitGroup
}

func CollectAndDispatch(threshold int, requests chan *PoolMsg, dispatches chan []*PoolMsg) {
	fmt.Println("collect and dispatch:", threshold)
	for {
		res := make([]*PoolMsg, threshold)
		for i := 0; i < threshold; i++ {
			res[i] = <-requests
		}
		dispatches <- res
	}
}
