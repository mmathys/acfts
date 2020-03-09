# Client-side

## Setup

### mysql (e.g. cluster number 0)

```
$ mysql.server start
$ mysql -u root -p
mysql> create database acfts_client_0
```

## Run (e.g. run with port 3000)

```
$ go run client.go
Input cluster number: 0
...
[GIN-debug] Listening and serving HTTP on :3000/
```

## Usage

You can create a request of a new transaction with  `struct generalTx`.

### generalTx

`generalTx` is a golang struct.

```go
type generalTx struct {
	From    model.Address
	To      []model.Address
	Amounts []int
}
```

You can designate the sender (`From`), the receivers (`To`) and amounts for each receiver (`Amounts`).

The orders of `To` and `Amounts` must correspond to each other.

### insideTx

If you want to make a transaction inside a cluster, you can use `insideTx` instead of `generalTx`.

```go
type insideTx struct {
	From    int
	To      []int
	Amounts []int
}
```

`insedeTx` has `From`, `To` and `Amounts` as well as `generalTx`, but the types of the first two are not `model.Address`, but `int` for simplicity.

It means you can designate the sender and the receivers by the client indexes in one cluster.

You can convert `insideTx` to `generalTx` with `convertInsideTxs()`.

### Execution

You generate clients in each cluster with `	generateClients(numClients, myurl)`.

Then, you get all addresses of clients including other clusters with `getAllAddrs()`.

You need to call `collectOtherAddrs()` before calling `getAllAddrs()` to get addresses of different clusters.

```go
collectOtherAddrs(otherClients)
addrs := getAllAddrs()
```

Before you start creating general transactions, you should make genesis with `createGenesis(serverURLs, owner, amount)`. This function must be called only once among all clusters.

Now, you can designate the sender and the receivers with `addrs`.

```go
tx := generalTx{From: addrs[0], To: []int{addrs[1], addrs[2]}, Amounts: []int{10, 20}}
```

In this case, transactions (addrs[0] → addrs[1], 10) and (addrs[0] → addrs[2], 20) will be created.

To send a request of the transactions to servers to get their signatures, call `executeTxs()`.

```go
txs := []{tx}
executeTxs(serverURLs, txs)
```

## Benchmarking

You can benchmark the programs with `go test` command.

You need to prepare benchmark programs first.

```go
func BenchmarkScenario1(b *testing.B) {
	serverURLs := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	addrs := benchmarkSetup(serverURLs, 1, 4, 0)
	defer db.Close()

	atxs1 := []generalTx{
		{From: addrs[0], To: []model.Address{addrs[1]}, Amounts: []int{1}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executeTxs(serverURLs, atxs1)
	}
	manners.Close()
}
```

By default, 5 scenarios are prepared in `client_test.go`.

In those secarios, client {0, 1, 2, 3} are in cluster0 and client {4, 5, 6, 7} are in cluster1.

The genesis is bound to client 0.

Transaction: (from → to, amount)

- case 1: (0 → 1, 1)
- case 2: (0 → 1, 1), (1 → 0, 1)
- case 3: (0 → 1, 500000), (0 → 1, 500000), (2 → 3, 500000), (3 → 2, 500000)
- case 4: (random → random, 1)
- case 5: (0 → 4, 1)

In order to run the benchmark, execute the following command.

```
$ go test -bench Scenario1 -benchtime 10000x -timeout 24h
```

If you want to benchmark transactions between different clusters, you need to run the other client programs before executing the benchmark command.

#### Options

`-bench regexp`: Choose scenario(s) that you want to run benchmark matching a regular expression.

`-benchtime t`: Set the number of iterations (`b.N`). If you put `Nx` as t,  each benchmark will be run N times.

`-timeout d`: Set the limit time for benchmarks. The default is 10 mins.

You can see the details of the options [here](https://golang.org/cmd/go/#hdr-Testing_flags).

### pprof

pprof is a tool for profiling and visualization of programs. pprof provides various information such as flame graphs.

As setup, the server program needs to listen and serve the specified port for pprof.

```go
go func() {
	log.Println(http.ListenAndServe("localhost:7000", nil))
}()
```

When you want to profile localhost:7000, run the following command during the execution.

```
$ go tool pprof -http=":7777" -seconds 60 localhost:7000
```

When finishing the profiling, your browser will automatically open and show the result.

#### Options

`-http`: Specify host:post on which you can get an interactive web interface

`-seconds`: Set duration for time-based profile collection

#### Result

After finishing the profiling, the result is saved in your computer.

You can review the result on localhost:7777 whenever you want by the following command.

```
$ go tool pprof -http="localhost:7777" /path/to/pprof/pprof.samples.cpu.001.pb.gz
```

