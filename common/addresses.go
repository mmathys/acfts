package common

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
)

type NodeConfig struct {
	Net      string
	Port     int
	Key      []string
}

type Topology struct {
	Servers []NodeConfig
	Clients []NodeConfig
}

func readKey(keypair []string) *ecdsa.PrivateKey {
	res, err := crypto2.HexToECDSA(keypair[0])
	if err != nil {
		log.Fatal(err)
	}
	return res
}

var clients = map[[AddressLength]byte]Node{}
var clientKeys []Address
var servers = map[[AddressLength]byte]Node{}
var serverKeys []Address

func getIndex(addr Address) [AddressLength]byte {
	index := [AddressLength]byte{}
	copy(index[:], addr[:AddressLength])
	return index
}

func InitAddresses(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Panicf("%s does not exist", path)
	}

	file, err := os.Open(path)
	if err != nil {
		log.Panicf("error opening %s", path)
	}

	defer file.Close()


	var topo Topology
	dec := json.NewDecoder(file)
	dec.Decode(&topo)

	for _, client := range topo.Clients {
		key := readKey(client.Key)
		addr := MarshalPubkey(&key.PublicKey)
		index := getIndex(addr)
		clients[index] = Node{
			Net:      client.Net,
			Port:     client.Port,
			Key:      key,
		}
		clientKeys = append(clientKeys, addr)
	}

	for _, server := range topo.Servers {
		key := readKey(server.Key)
		addr := MarshalPubkey(&key.PublicKey)
		index := getIndex(addr)
		servers[index] = Node{
			Net:      server.Net,
			Port:     server.Port,
			Key:      key,
		}
		serverKeys = append(serverKeys, addr)
	}
}

func lookup(address Address) (Node, error) {
	index := getIndex(address)
	client, ok := clients[index]
	if ok {
		return client, nil
	}

	server, ok := servers[index]
	if ok {
		return server, nil
	}

	msg := fmt.Sprintf("could not find address 0x%x\n", address)
	return Node{}, errors.New(msg)
}

func GetNetworkAddress(address Address) (string, error) {
	res, err := lookup(address)
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s:%d", res.Net, res.Port), nil
	}
}

func GetKey(address Address) *ecdsa.PrivateKey {
	res, err := lookup(address)
	if err != nil {
		log.Panicf("could not find address 0x%x\n", address)
		return nil
	} else {
		return res.Key
	}
}

func GetPort(address Address) int {
	res, err := lookup(address)
	if err != nil {
		log.Fatalf("could not find address 0x%x", address)
		return -1
	} else {
		return res.Port
	}
}

func GetClients() []Address {
	return clientKeys
}

func GetServers() []Address {
	return serverKeys
}

func GetNumServers() int {
	return len(GetServers())
}
