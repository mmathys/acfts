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
	{"b799f2efc21900d0985f98112905f98804b99eedaf25716562618e8aec07aea7", "04bd0bb937b9623e1df5329262a8bda685a5049f5a6c44b90030744737e5916d627ab6e7a85a57007d46eebb32291c7760101d33615cde137accb8127ab0f593c3"},
	{"a45ab9fc855b3354f0069083192c907763648c18f0abc699f40a1bb256b4709b", "047bad17f5b243978dbb99b03400e057e95296b4d9612eef212e53b3f5792e5d5c07a815da98313c1ced143712ba0beb8647ed882ccd16599a81504cdc166f4293"},
	{"4c24865b51e00655e6cc383d7e22e8d00d47728a1e5d5f71dca6b9a300c7ea2f", "04a5527f944046ced1526264f2a7484d73402c494a1d9cc50938242d36cc98bda0c9121aefc210df3964bf870bf550b5d2befc235f7a7cd4e52cad26c5d049d8dd"},
	{"55f5776bcda3a3e1d8f2810a11e3c3c4ab7ef30bfd3c834625dc7dd9d59f6cd1", "0479ee713a7bcca5a1d83b0addfff9cff44673fddecce88b6d299f6ba6f9d6372e1c34ed1c65c4cb50f86341b3a4d3a1d4d652b3f869218f8b09fda7072007c186"},
	{"9234da203c4108c9484ee8db64fdc4c0bd59bf3686495768d2d88ed2878c72d2", "04737a47bf6e71b41186be97b2edc6a4d348d16cb395b154e6a3e90dd01b8ba9207d4dbc4c6460f27059d833c5a21aaa6122e8033057e658d1ad4f64c480d711ef"},
	{"e67785183231a26a17b729e6683506aeccf5497c7d6735a729955539b49f0e57", "04a32653c7439dc87c9c3f639adf00f39643f094b01095f66f7dac658dc05db4f586e7136afce16ca0ec9809c80177446c6aaad9fe7e4048f550a258be66370835"},
	{"45734a533232304fce8f98e76aac7f3db9ea1bbaa5f28391416f0537a80ee2e4", "04fa4de39285b1a7aff5b226792f55cbb3dcb1a69b1a58f07e75e4a2335cf2b486ca2500f09d82ee78dc983b57f5a9b39593a9012078b32874d78562f6379c84fb"},
	{"d14cf893a767cf9354ee594afa2eee111391c78f88a34d2074a1ffad2cc01f35", "04a775a6e44e897f96232ca2c752070b46681dd12eecd1c494315e46634f499ed83a1018854c00780c5ab7f40d4d23f41d81d1ff5a8724f0fea0a9c7a6ed90e97b"},
	{"cab38d74a8311482489c99209604dbba09296be182df069cdebc37f8800962a5", "043b1e0a072dcb1f69d847ef74473665ebc20450d3594d708d13d7dc61744ab443f735848675ffc8036a71c57d7a844a42e745f3fea4482e08d8f4cb5032a6ecc9"},
	{"d5830ac0b696931ee3f2e88fbb1b859198aa01d7b49f1fedc239355f9a5752a7", "04d34e04d720691c6d392cfc49d59d501e813ab6d879a1aba563a26b76c6f3109791c7aa8297fca7d399b068ff1baae0c9587f371ca04ccadded614a49e474cc25"},
	{"18496d30042a8c5eb58772c4c23c5557ade72f3edbb24259fd154613a9250229","0401eb21077a36b36c09f590dcd52ef3feb62b36177df7314f6818485b78296998ce8fedb19e62f7df72f52dc9f9ebf7f0f1782ab6bbcb033ccf0dd75e60c9419d"},
	{"fc3ff9856860914dd5cbd8504b9a61ea85414cc8d634ffa5a384e57fec6d5d84","04ee67afaa950a4e9a3358f02ed4e56716c1c9b5b9f17d6c1d94ac9a31c3f32f669830cfdacf0b408164d3237ec5e6b8995b8e62d05cbfafce149051f582dcd7be"},
	{"424cf13a70d3b5706d7b8fd1be3a957bae5163b5e7fccf3ad35901c076f68768","043c9b75ee68c84752e6da58a3ec366ba82e67815e15100604eea17897a76d8b02e49f3bdd6a08cc657c6508fdbc6e72d247a63eb5c7148aa62f2a788ab85ba367"},
}

func readKey(keypair []string) *ecdsa.PrivateKey {
	res, err := crypto2.HexToECDSA(keypair[0])
	if err != nil {
		log.Fatal(err)
	}
	return res
}

var m = map[common.Address]common.Node{
	common.Address{0}:  {"client", common.Address{0}, "http://localhost", 5555, readKey(generatedKeyPairs[0])}, // 0x00 (client)
	common.Address{1}:  {"client", common.Address{1}, "http://localhost", 5556, readKey(generatedKeyPairs[1])}, // 0x01 (client)
	common.Address{2}:  {"client", common.Address{2}, "http://localhost", 5557, readKey(generatedKeyPairs[2])}, // 0x02 (client)
	common.Address{3}:  {"client", common.Address{3}, "http://localhost", 5558, readKey(generatedKeyPairs[3])},
	common.Address{4}:  {"client", common.Address{4}, "http://localhost", 5559, readKey(generatedKeyPairs[4])},
	common.Address{5}:  {"client", common.Address{5}, "http://localhost", 5560, readKey(generatedKeyPairs[5])},
	common.Address{6}:  {"client", common.Address{6}, "http://localhost", 5561, readKey(generatedKeyPairs[6])},
	common.Address{7}:  {"client", common.Address{7}, "http://localhost", 5562, readKey(generatedKeyPairs[7])},
	common.Address{8}:  {"client", common.Address{8}, "http://localhost", 5563, readKey(generatedKeyPairs[8])},
	common.Address{9}:  {"client", common.Address{9}, "http://localhost", 5564, readKey(generatedKeyPairs[9])},
	common.Address{10}: {"client", common.Address{10}, "http://localhost", 5565, readKey(generatedKeyPairs[10])},
	common.Address{11}: {"client", common.Address{11}, "http://localhost", 5566, readKey(generatedKeyPairs[11])},
	common.Address{12}: {"client", common.Address{12}, "http://localhost", 5567, readKey(generatedKeyPairs[12])},
	common.Address{13}: {"client", common.Address{13}, "http://localhost", 5568, readKey(generatedKeyPairs[13])},
	common.Address{14}: {"client", common.Address{14}, "http://localhost", 5569, readKey(generatedKeyPairs[14])},
	common.Address{15}: {"client", common.Address{15}, "http://localhost", 5570, readKey(generatedKeyPairs[15])},
	common.Address{16}: {"server", common.Address{16}, "http://localhost", 6666, readKey(generatedKeyPairs[16])},
	common.Address{17}: {"server", common.Address{17}, "http://localhost", 6667, readKey(generatedKeyPairs[17])},
	common.Address{18}: {"server", common.Address{18}, "http://localhost", 6668, readKey(generatedKeyPairs[18])},
	common.Address{19}: {"server", common.Address{19}, "http://localhost", 6669, readKey(generatedKeyPairs[19])},
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
		log.Panicf("could not find address 0x%x\n", address)
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
		common.Address{3},
		common.Address{4},
		common.Address{5},
		common.Address{6},
		common.Address{7},
		common.Address{8},
		common.Address{9},
		common.Address{10},
		common.Address{11},
		common.Address{12},
		common.Address{13},
		common.Address{14},
		common.Address{15},
	}
}

func GetServers() []common.Address {
	return []common.Address{
		common.Address{16},
		//common.Address{17},
		//common.Address{18},
		//common.Address{19},
	}
}

func GetNumServers() int {
	return len(GetServers())
}
