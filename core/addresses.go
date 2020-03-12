package core

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
	"github.com/mmathys/acfts/common"
	"log"
)

type Entry struct {
	network string
}

// note: the public key can be recovered from the private key. not necessarily needed in this array.
var generatedKeyPairs = [][]string{
	{"439ffd45a52551b1a5470b618a18b58e70c6f1201422760c02ffab17a1dc73e1", "04ae9c97b237b6dc65081c41bbbb45419933353a848867b4ea61bd92fa202669532ccac6cd35324ef47c5c402847f74f6e6d0873a81a9cee9290bf796ae0348ca2"},
	{"0401df1f7bdbff562bc87f21c4fcbe9bf09b2ffe8c2e4609b32f9b11bc79b9a0", "04799fb6d7685ea7fb7a0477db235dcebf4ec7b871a26287df9552a2ba035a46d00d538409e3d5ad4834d0cfa2f82c5426ebc78ee8d1bacd4d4783ad19c56c123c"},
	{"7675b8ef1a5ba5cfc037b0d778fc74374aef4b48e8f33d03b0c0236780941988", "048dadeefa70cfeb43f9e15a7df924cf1bb7878b081039ec2b9bce593252fe732ca043ec91ba4dcb639d7e70d778ecadd4ad2aadb11963c1206d3a6b62a8c0ee04"},
	{"6a26d73ecbd72ca5d4c97f3860656ac0bbe73b025f713e13e023bf9311e69c21", "04af34ebdcf83c3c122b20a0e0bbcefd8292400e829c2dcec6e63415a0b68659321b554f148be777aee885ddb226415bb08ead0fae84d7e997e13754ac518a46d5"},
	{"4fecbb37a9fee1857d72e46e9330f0fb012b0c2eab1c89cc7a2d5044b064570b", "0430ae31df0cd59c06648bf6183ddedb7f500f74e6fcf6b8088525a764e39b4665333c8e45285036d60ff7e9e963b2e98595706846c73f4f083a11d133e500db1a"},
	{"b8f360aeb3136db2f27ab44592f05b4cbc50ac8dd8f5f046dee4e64056ece91f", "042d67c7ba4deb532ce0212dcf71bb3feffe305785f3619dc428aa03d54c287e16386ec39b8117871668769db26e36405a66877d3aa5fc15039af0a798dc77f504"},
	{"76c41be05a08e389cb64d1e8cf9423964628e26aa8821703c7d55240b1cf435d", "04af86aeda3e395c01a6b9b00120f798c1c40ba6c238ea5d6a6e1d2de88d47c51081812a59114b894647e90a209073aa1448efab3595fb3c534d620892e4824f3e"},
}

func readKey(keypair []string) *ecdsa.PrivateKey {
	res, err := crypto2.HexToECDSA(keypair[0])
	if err != nil {
		log.Fatal(err)
	}
	return res
}

var m = map[common.Address]common.Node{
	common.Address{0}: {"client", common.Address{0}, "http://localhost", 5555, readKey(generatedKeyPairs[0])}, // 0x00 (client)
	common.Address{1}: {"client", common.Address{1}, "http://localhost", 5556, readKey(generatedKeyPairs[1])}, // 0x01 (client)
	common.Address{2}: {"client", common.Address{2}, "http://localhost", 5557, readKey(generatedKeyPairs[2])}, // 0x02 (client)
	common.Address{3}: {"server", common.Address{3}, "http://localhost", 6666, readKey(generatedKeyPairs[3])},             // 0x03 (server)
	common.Address{4}: {"server", common.Address{4}, "http://localhost", 6667, readKey(generatedKeyPairs[4])},             // 0x04 (server)
	common.Address{5}: {"server", common.Address{5}, "http://localhost", 6668, readKey(generatedKeyPairs[5])},             // 0x05 (server)
	common.Address{6}: {"server", common.Address{6}, "http://localhost", 6669, readKey(generatedKeyPairs[6])},             // 0x06 (server)
}

func GetNetworkAddress(address common.Address) (string, error) {
	res, ok := m[address]
	if ok {
		return fmt.Sprintf("%s:%d", res.Net, res.Port), nil
	} else {
		msg := fmt.Sprintf("could not find address 0x%x\n", address)
		return "", errors.New(msg)
	}
}

func GetKey(address common.Address) *ecdsa.PrivateKey {
	res, ok := m[address]
	if ok {
		return res.Key
	} else {
		log.Fatal("could not find address")
		return nil
	}
}

func GetPort(address common.Address) int {
	res, ok := m[address]
	if ok {
		return res.Port
	} else {
		log.Fatal("could not find address")
		return -1
	}
}

func GetClients() []common.Address {
	return []common.Address{
		common.Address{0},
		common.Address{1},
		common.Address{2},
	}
}

func GetServers() []common.Address {
	return []common.Address{
		common.Address{3},
		//common.Address{4},
		//common.Address{5},
		//common.Address{6},
	}
}
