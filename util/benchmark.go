package util

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"sync/atomic"
	"time"
)

const (
	header = "time,tx,cpu_util"
)

func CountTx(txCounter *int32) {
	atomic.AddInt32(txCounter, 1)
}

func Ticker(txCounter *int32) {
	fmt.Println(header)
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		printTx(txCounter)
	}
}

func printTx(txCounter *int32) {
	count := atomic.SwapInt32(txCounter, 0)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	cpuUtil, err := cpu.Percent(time.Second, false)
	if err != nil {
		cpuUtil = []float64{0}
	}


	if count > 0 {
		fmt.Printf("%d,%d,%f\n", timestamp, count, cpuUtil[0])
	}
}
