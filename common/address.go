package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"golang.org/x/crypto/sha3"
	"log"
	"math"
	"math/rand"
	"os"
)

/**
Key:
[
	0: Public Key
	1: Private Key
]
*/

type KeyConfig struct {
	PublicKey  string
	PrivateKey string
	Mode       string
}

type ClientNodeConfig struct {
	Key      KeyConfig
	Instance Instance
}

type ServerNodeConfig struct {
	Key       KeyConfig
	Instances []Instance
}

type TopologyConfig struct {
	BLSMasterKey KeyConfig
	Servers      []ServerNodeConfig
	Clients      []ClientNodeConfig
}

func read(conf KeyConfig, index int) *Key {
	var mode int
	if conf.Mode == "eddsa" {
		mode = ModeEdDSA
	} else if conf.Mode == "bls" {
		mode = ModeBLS
	} else {
		log.Panicf("encountered invalid mode %v (only \"eddsa\" and \"bls\" are valid)\n", conf.Mode)
	}

	res, err := DecodeKey(mode, index, conf.PublicKey, conf.PrivateKey)
	if err != nil {
		panic(err)
	}
	return res
}

const IndexLength = 48 // max(EdDSAPublicKeyLength, BLSPublicKeyLength)
var clients map[[IndexLength]byte]ClientNode
var ClientAddresses []Address
var servers map[[IndexLength]byte]ServerNode
var ServerAddresses []Address
var masterKey *Key

func getIndex(addr Address) [IndexLength]byte {
	index := [IndexLength]byte{}
	copy(index[:], addr[:])
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

	clients = map[[IndexLength]byte]ClientNode{}
	ClientAddresses = []Address{}
	servers = map[[IndexLength]byte]ServerNode{}
	ServerAddresses = []Address{}

	var topology TopologyConfig
	dec := json.NewDecoder(file)
	dec.Decode(&topology)

	for i, client := range topology.Clients {
		key := read(client.Key, i)
		index := getIndex(key.GetAddress())
		clients[index] = ClientNode{
			Instance: client.Instance,
			Key:      key,
		}
		ClientAddresses = append(ClientAddresses, key.GetAddress())
	}

	for i, server := range topology.Servers {
		key := read(server.Key, i)
		index := getIndex(key.GetAddress())
		servers[index] = ServerNode{
			Instances: server.Instances,
			Key:       key,
		}
		ServerAddresses = append(ServerAddresses, key.GetAddress())
	}

	masterKeyConfig := topology.BLSMasterKey
	if masterKeyConfig.Mode == "bls" {
		masterKey = read(masterKeyConfig, 0)
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

func GetKey(address Address) *Key {
	client, err := lookupClient(address)
	if err == nil {
		return client.Key
	}

	server, err := lookupServer(address)
	if err == nil {
		return server.Key
	}

	log.Panicf("could not find address %x\n", address)
	return nil
}

func GetBLSMasterKey() *Key {
	return masterKey
}

func GetBLSMasterPublicKey() bls.PublicKey {
	return masterKey.BLS.Address
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
	num |= uint32(hash[hashLen-1])
	num |= uint32(hash[hashLen-2]) << 8
	num |= uint32(hash[hashLen-3]) << 16
	num |= uint32(hash[hashLen-4]) << 24

	return int(num % numInstances)
}

func IsValidServer(serverAddr Address) bool {
	_, err := lookupServer(serverAddr)
	return err == nil
}

func GetClients() []Address {
	return ClientAddresses
}

func GetServers() []Address {
	return ServerAddresses
}

func QuorumSize() int {
	n := GetNumServers()
	return int(math.Ceil(2.0 / 3.0 * float64(n)))
}

// returns shuffled >2/3 server quorum
func ServerQuorum() []Address {
	numServers := GetNumServers()
	quorumSize := QuorumSize()
	keys := make([]Address, quorumSize)
	indexes := rand.Perm(numServers)[:quorumSize]

	for i, v := range indexes {
		keys[i] = ServerAddresses[v]
	}
	return keys
}

func GetNumServers() int {
	return len(GetServers())
}
