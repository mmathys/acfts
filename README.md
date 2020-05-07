<div align="center">
    <h1>
        ACFTS
    </h1>
    <p>
        Asynchronous Consensus-Free Transaction System
    </p>
</div>

## Table of Contents
1. [Install](#install)
2. [Introduction](#introduction)
3. [Executables](#executables)
4. [Code](#code)
5. [Mirroring](#mirroring)

## Install

**Prerequisites**: [Go](https://golang.org/doc/install) is required.

Clone the repo:

```bash
git clone git@github.com:mmathys/acfts.git

# or use the ETH Mirror:
git clone git@gitlab.ethz.ch:disco-students/fs20/mmathys-consensus-free-transaction-systems.git
```

Add `$GOPATH/bin` to your `$PATH`
```
# Add this line to your .bashrc, .zshrc or similar
export PATH=$PATH:$GOPATH/bin

# reload the .bashrc
source ~/.bashrc
```

Build source and install executables:

```bash
cd acfts
go build ./...
go install ./...
```

## Introduction

ACFTS is a asynchronous consensus-free transaction system. It consists of trusted servers and untrusted clients. Each server
and client have a unique ECDSA key. The address of each server and client is its public key.

The client has a CLI which allows easy transfer of currency. In this setup, each client automatically gets 100 valid credits.
The client CLI can be accessed after starting the client executable. 

### Topologies

Topologies are defined as JSON configuration files in `topologies/`.

Topologies encode all necessary information about server and clients, for example: address, keys, network address,
replication instances. In a system, every server and client have to use the same topology.

When launching a server or a client, its configuration can be given with a topology, an assigned address (which
must correspond to a node in the topology). When server shard replication is used, and replication instance index is
passed additionally.

## Executables

The CLIs are installed in `$GOPATH/bin` and are called `server` and `client` respectively.

### Server

```bash
server                      # CLI executable
    --address <address>     # Address assigned to server. Format: 0x...
    --topology <file>       # Topology configuration file
    --benchmark             # Outputs number of tx/s to stdout
    --pprof                 # Enables pprof profiler
    --adapter rpc           # Network adapter    
    --instance <number>     # Replication instance (used for sharding)
    --help                  # Prints help
```

### Client

Start the CLI:

```bash
client                      # CLI executable
    --address <address>     # Address assigned to client. Format: 0x...
    --topology <file>       # Topology configuration file
    --benchmark             # If set: outputs benchmark
    --adapter rpc           # Network adapter  
    --help                  # Prints help  
```

Run CLI commands:

```bash
> help                    # Show the help section
> send <address> 100      # Send 100 credits to <address>. Format: 0x...
> utxo                    # Show local UTXOs
> balance                 # Show balance
> info                    # Show client information
> clear                   # Clear console
```

## Code

| Folder | Description |
| :---: | :---: |
| `benchmark` | code for running benchmarks |
| `client` | Client specific code |
| `common` | Code which is used in both server and client, for example ECDSA-related code. |
| `docs` | Documentation |
| `server` | Server specific code |
| `tests` | Tests (not used in the executables) |
| `topologies` | Topology config files |
| `util` | Utility functions |
| `wallet` | Wallet specific code |


## Mirroring

The original repository is hosted on  [GitHub](https://github.com/mmathys/acfts) and mirrored to a [repository hosted
on gitlab.ethz.ch](https://gitlab.ethz.ch/disco-students/fs20/mmathys-consensus-free-transaction-systems).