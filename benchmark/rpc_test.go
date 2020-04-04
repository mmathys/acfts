package benchmark

import (
	"errors"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/rpc"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func setupServer() {
	runtime.SetBlockProfileRate(1)
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

func runClient(wg *sync.WaitGroup, numReq int) {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &Args{7, 8}
	for i := 0; i < numReq; i++ {
		var reply int
		err = client.Call("Arith.Multiply", args, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
	}
	wg.Done()
}

func TestMain(m *testing.M) {
	go setupServer()
	time.Sleep(time.Second)
	code := m.Run()
	os.Exit(code)
}

func TestRPC(t *testing.T) {
	numClients := 100
	numReq := 1000 * 1000
	var wg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go runClient(&wg, numReq / numClients)
	}
	wg.Wait()
}
