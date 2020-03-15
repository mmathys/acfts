package util

import (
	"fmt"
	"sync/atomic"
	"time"
)

func CountTx(txCounter *int32) {
	atomic.AddInt32(txCounter, 1)
}

func Ticker(txCounter *int32) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		printTx(txCounter)
	}
}

func printTx(txCounter *int32) {
	count := atomic.SwapInt32(txCounter, 0)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	if count > 0 {
		fmt.Printf("%d,%d\n", timestamp, count)
	}
}
