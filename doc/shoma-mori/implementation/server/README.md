# Server-side

## Setup

### mysql (e.g. port 8080)

```
$ mysql.server start
$ mysql -u root -p
mysql> create database acfts_8080
```

## Run (e.g. run with port 8080)

```
$ go run main.go
Input port number: 8080
...
[GIN-debug] Listening and serving HTTP on :8080/
```

## Example

```shell
#!/bin/bash

# Create genesis
curl -H 'Content-Type:application/json' -d "$JSON" http://localhost:8080/genesis

# Request for the new transaction
JSON=$(cat <<EOS
{
        "inputs": [
                {
                        "utxo": {
                                "address1": "example01",
                                "address2": "example02",
                                "previous_hash": "genesis",
                                "index": 0,
                                "server_signatures": [
                                        {
                                                "address1": "gene",
                                                "address2": "sis",
                                                "signature1": "dum",
                                                "signature2": "my"
                                        },
                                        {
                                                "address1": "gene",
                                                "address2": "sis",
                                                "signature1": "dum",
                                                "signature2": "my"
                                        }]
                        },
                        "siblings": [],
                        "signature1": "sig01",
                        "signature2": "sig02"
                }
        ],
        "outputs": [
                {
                        "amount": 200,
                        "address1": "example11",
                        "address2": "example12",
                        "previous_hash": "a36d289b2ed18b196f07f59489710f45d92a97fdc2e2492dec67cf9cadb3a733",
                        "index": 0
                }
        ]
}
EOS
)

curl -H 'Content-Type:application/json' -d "$JSON" http://localhost:8080/transaction
```

Execute this shell script.
```
$ chmod 755 test.sh
$ ./test.sh
```

## Benchmarking with pprof

pprof is a tool for profiling and visualization of programs. pprof provides various information such as flame graphs.

As setup, the server program needs to listen and serve the specified port for pprof.

```go
go func() {
	log.Println(http.ListenAndServe("localhost:8000", nil))
}()
```

When you want to profile localhost:8000, run the following command during the execution.

```
$ go tool pprof -http=":8888" -seconds 60 localhost:8000
```

When finishing the profiling, your browser will automatically open and show the result.

#### Options

`-http`: Specify host:post on which you can get an interactive web interface

`-seconds`: Set duration for time-based profile collection

#### Result

After finishing the profiling, the result is saved in your computer.

You can review the result on localhost:8888 whenever you want by the following command.

```
$ go tool pprof -http="localhost:8888" /path/to/pprof/pprof.samples.cpu.001.pb.gz
```

