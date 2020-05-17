package common

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"log"
	"os"
)

/**
Key:
[
	0: Public Key
	1: Private Key
]
 */

type ClientNodeConfig struct {
	Key      []string // 0: Address, 1: Private Key
	Instance Instance
}

type ServerNodeConfig struct {
	Key       []string // 0: Address, 1: Private Key
	Instances []Instance
}

type TopologyConfig struct {
	Servers []ServerNodeConfig
	Clients []ClientNodeConfig
}

func readKey(keypair []string) (Address, PrivateKey) {
	addr, err := hex.DecodeString(keypair[0])
	if err != nil {
		panic(err)
	}

	priv, err := hex.DecodeString(keypair[1])
	if err != nil {
		panic(err)
	}

	if len(addr) != AddressLength {
		panic("topology: encountered wrong address length")
	}

	if len(priv) != PrivateKeyLength {
		panic("topology: encountered wrong private key length")
	}

	return addr, priv
}

var clients map[[AddressLength]byte]ClientNode
var clientKeys []Address
var servers map[[AddressLength]byte]ServerNode
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

	clients = map[[AddressLength]byte]ClientNode{}
	clientKeys = []Address{}
	servers = map[[AddressLength]byte]ServerNode{}
	serverKeys = []Address{}

	var topology TopologyConfig
	dec := json.NewDecoder(file)
	dec.Decode(&topology)

	for _, client := range topology.Clients {
		addr, key := readKey(client.Key)
		index := getIndex(addr)
		clients[index] = ClientNode{
			Instance: client.Instance,
			Key:      key,
		}
		clientKeys = append(clientKeys, addr)
	}

	for _, server := range topology.Servers {
		addr, key := readKey(server.Key)
		index := getIndex(addr)
		servers[index] = ServerNode{
			Instances: server.Instances,
			Key:       key,
		}
		serverKeys = append(serverKeys, addr)
	}
}

func lookupClient(address Address) (ClientNode, error) {
	index := getIndex(address)
	client, ok := clients[index]
	if ok {
		return client, nil
	} else {
		msg := fmt.Sprintf("could not find client %x\n", address)
		return ClientNode{}, errors.New(msg)
	}
}

func lookupServer(address Address) (ServerNode, error) {
	index := getIndex(address)
	server, ok := servers[index]
	if ok {
		return server, nil
	} else {
		msg := fmt.Sprintf("could not find server %x\n", address)
		return ServerNode{}, errors.New(msg)
	}
}

func lookupServerInstance(address Address, instanceIndex int) (Instance, error) {
	server, err := lookupServer(address)
	if err != nil {
		return Instance{}, err
	}

	if instanceIndex < 0 || instanceIndex > len(server.Instances) {
		msg := fmt.Sprintf("instance index is out of bounds. got: %d, required: 0..%d", instanceIndex, len(server.Instances)-1)
		return Instance{}, errors.New(msg)
	}

	return server.Instances[instanceIndex], nil
}

func GetClientNetworkAddress(address Address) (string, error) {
	res, err := lookupClient(address)
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s:%d", res.Instance.Net, res.Instance.Port), nil
	}
}

func GetServerNetworkAddress(address Address, instanceIndex int) (string, error) {
	res, err := lookupServerInstance(address, instanceIndex)
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s:%d", res.Net, res.Port), nil
	}
}

func GetKey(address Address) *PrivateKey {
	client, err := lookupClient(address)
	if err == nil {
		return &client.Key
	}

	server, err := lookupServer(address)
	if err == nil {
		return &server.Key
	}

	log.Panicf("could not find address %x\n", address)
	return nil
}

func GetClientPort(address Address) int {
	res, err := lookupClient(address)
	if err != nil {
		log.Fatalf("could not find address %x", address)
		return -1
	} else {
		return res.Instance.Port
	}
}

func GetServerPort(address Address, instanceIndex int) int {
	res, err := lookupServerInstance(address, instanceIndex)
	if err != nil {
		log.Fatalf("could not find address %x", address)
		return -1
	} else {
		return res.Port
	}
}

func GetClientBalance(address Address) int {
	res, err := lookupClient(address)
	if err != nil {
		log.Fatalf("could not find address %x", address)
		return -1
	} else {
		return res.Balance
	}
}


func GetServerInstanceIndex(serverAddr Address, clientAddr Address) int {
	server, err := lookupServer(serverAddr)
	if err != nil {
		panic(err)
	}

	d := sha3.New256() // 256 bits / 32 bytes
	d.Write(clientAddr)
	hash := d.Sum(nil)

	// get the four most least significant bytes to the address, then mod it with the number of server instances
	numInstances := uint32(len(server.Instances))
	var num uint32
	hashLen := len(hash)
	num |= uint32(hash[hashLen - 1])
	num |= uint32(hash[hashLen - 2]) << 8
	num |= uint32(hash[hashLen - 3]) << 16
	num |= uint32(hash[hashLen - 4]) << 24

	return int(num % numInstances)
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
