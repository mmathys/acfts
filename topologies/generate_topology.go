package main

import (
	"encoding/json"
	"io/ioutil"
)

var generatedKeysServer = [][]string{
	{"a77000e6e7e77df500e39144e34e31f1715c98a5ecaaf34b8863f2d907c35ab4", "26665d594658641d981cb13e74e130c12bb5a0924c217e10b2a0e20437cabf97a77000e6e7e77df500e39144e34e31f1715c98a5ecaaf34b8863f2d907c35ab4"},
	{"c8766183072ad3a0aeb17e672152bff7d5382a62f14a818f2a7dbe3c123d1fc3", "19d771ac0e1f2b76082820840dd8b7df7c7e2614f8d7346186460946b5f6173dc8766183072ad3a0aeb17e672152bff7d5382a62f14a818f2a7dbe3c123d1fc3"},
	{"564c67c3cd08863a570e632c60409d57b02f576a0582bcaf297aa4610f98aead", "eb1a60b82e72545395f8ce4334703f13597985aed204107ea5bdaf6289c38a68564c67c3cd08863a570e632c60409d57b02f576a0582bcaf297aa4610f98aead"},
	{"4a155d1c1d86854aae1008373c3e9a92dd02cbbef52e5372aee13c2efe22398d", "e76b8dee9a5eca0000da56aadaa39e8cdb564ee6f82008052caad39741e72b134a155d1c1d86854aae1008373c3e9a92dd02cbbef52e5372aee13c2efe22398d"},
}

var generatedKeysClient = [][]string{
	{"5e7da91d8083e79d742f39a1e218c92fdfb8e5684b60c3d2d82bce1fcc2c8084", "135f014fc76ca27c3be0227ac8c356ab7c6cb38dda1f350bb8fda7487f8aa2dd5e7da91d8083e79d742f39a1e218c92fdfb8e5684b60c3d2d82bce1fcc2c8084"},
	{"13a619eba72c54e669d67b910d3e927f3dc1bb02764b69b4163345d55c2eb0aa", "43ccf68f7190d1fc0f6fa6b2dc35e1244290c816cd0eb4fc12d73a5c1914a86e13a619eba72c54e669d67b910d3e927f3dc1bb02764b69b4163345d55c2eb0aa"},
	{"826f9850a0685f79877b8049f3216a35b4a5b9894d4f70b50946b254ea878e22", "5f4939311ce1bc5cd889e9d35a7ea7c2db64db0d5e22eb8d8914b044ed388141826f9850a0685f79877b8049f3216a35b4a5b9894d4f70b50946b254ea878e22"},
	{"d46a71d600ad623d21ced5c5398984def57875fe4ab6d9b5d4ed6883d06e6884", "4e6eb64343c2a3c9f5241771484545fdb1bcc2749225b87d9589034c1bb81414d46a71d600ad623d21ced5c5398984def57875fe4ab6d9b5d4ed6883d06e6884"},
	{"6c8c77537d8640d7b530a4695310151bbd8b42d45f9234be29ee83b3e9419474", "ac8a1e63bb6440fd4937849cbe9296a9a92925ff19bd7e0f1f29500887a682046c8c77537d8640d7b530a4695310151bbd8b42d45f9234be29ee83b3e9419474"},
	{"cf6d3ac66bec3408012f736463590c62739798036681cd73be39071ed00a2b67", "7f8f40262ab33bfd5a973622e046fab042123b88a2114fcb97c3e6858c4b1099cf6d3ac66bec3408012f736463590c62739798036681cd73be39071ed00a2b67"},
	{"00bb932e1c935c102cf1d53452b0f10ff9ece1264ec96a6ace30bb1037ea5507", "827d7b276e4ff3efbc3d3957f156a2a958b3deae890277ab7e1a847a0415dccf00bb932e1c935c102cf1d53452b0f10ff9ece1264ec96a6ace30bb1037ea5507"},
	{"62a0061c9ec6b2a1f186dfc0f490a7eb158b072ba2d7ef049248b1337917a817", "163f4cbafb36ee4a925e22e373153e01d807a5f7de4ede92e80b93a5a4fa18c362a0061c9ec6b2a1f186dfc0f490a7eb158b072ba2d7ef049248b1337917a817"},
	{"3f4d44a67deccaa314dd203105ff9231f382a84561013cf253352059b90b692d", "e650af1698bc8299268788a9ae97935f21e970180b4f7f1a1d94e3c550e736fb3f4d44a67deccaa314dd203105ff9231f382a84561013cf253352059b90b692d"},
	{"4d886b2926a41ef39fd990c90b298cafbddd811cca2c8e678a0066c44bf31523", "faf4c5c35216b0209957b098b779cc6da9a8dbe0106556260b3035cce8b175124d886b2926a41ef39fd990c90b298cafbddd811cca2c8e678a0066c44bf31523"},
	{"bc4e0b9464a805d958d1e09e4949b451e191cf25924c52a6bed010df5d857282", "679bde789d45515b076a68517bdfbfba042f55dad50df3043b14a8384f45006fbc4e0b9464a805d958d1e09e4949b451e191cf25924c52a6bed010df5d857282"},
	{"be89288412fd3d03c0b6f512d1e5bbbc91957bf86040e0959b200933825a9e1e", "22b59214fa4fcaaadb4ffb01bd3d413e6fdfb93faa21b27f241e64c4001137efbe89288412fd3d03c0b6f512d1e5bbbc91957bf86040e0959b200933825a9e1e"},
	{"278b81feb12c387a49fadc1649e66690dd7d92e8b0419f16c1eca80b81e05e8a", "64e4ccdd76cd951aa179d32741450c925d5ee1eb70241070798badabff86dfed278b81feb12c387a49fadc1649e66690dd7d92e8b0419f16c1eca80b81e05e8a"},
	{"45b90f249fef00684a0c1ca2a395f15c8c3aedde9691b52e9ee673b547265bf5", "51d929f193bedcf1edffa72825afe4cb7757a18c81cadc7b53b9d26b2759e39445b90f249fef00684a0c1ca2a395f15c8c3aedde9691b52e9ee673b547265bf5"},
	{"26ae03b17e85cc8bee2274310c6b36ab916709ee5c93f8ae3c67b4507b48569f", "0dd90e6463ec8cab750b414f9c1d79bf7eef9cefb65304457f1abaac9f2ab64126ae03b17e85cc8bee2274310c6b36ab916709ee5c93f8ae3c67b4507b48569f"},
	{"a023dd2de6a360810033baf18ad1ea8a513bca533a67bc9b572eb3e7021dbc84", "0de3ea14676fa097d219598d062224eb374b04675e68a5dd5b778261691545fca023dd2de6a360810033baf18ad1ea8a513bca533a67bc9b572eb3e7021dbc84"},
}

type Instance struct {
	Net  string // network address, with http
	Port int    // port
}

type ClientNode struct {
	Instance Instance
	Key      []string
	Balance  int
}

type ServerNode struct {
	Instances []Instance
	Key       []string
}

type Topology struct {
	Servers []ServerNode
	Clients []ClientNode
}


func check(e error) {
	if e != nil {
		panic(e)
	}
}

// prints topology to stdout
func main() {
	err := ioutil.WriteFile("topologies/localSimple.json", localSimple(), 0644)
	check(err)
	err = ioutil.WriteFile("topologies/localSimpleExtended.json", localSimpleExtended(), 0644)
	check(err)
	err = ioutil.WriteFile("topologies/localFull.json", localFull(), 0644)
	check(err)
	err = ioutil.WriteFile("topologies/dockerSimple.json", dockerSimple(), 0644)
	check(err)
	err = ioutil.WriteFile("topologies/aws.json", aws(), 0644)
	check(err)
}

func getClientNetwork(dockerNetwork bool) string {
	if dockerNetwork {
		return "client"
	} else {
		return "localhost"
	}
}

func getServerNetwork(dockerNetwork bool) string {
	if dockerNetwork {
		return "server"
	} else {
		return "localhost"
	}
}

func config(numClients int, numServers int, numServerInstances int, dockerNetwork bool) []byte {
	topo := Topology{}

	counter := 0
	for i := 0; i < numServers; i++ {
		var instances []Instance
		for j := 0; j < numServerInstances; j++ {
			instances = append(instances, Instance{
				Net:  getServerNetwork(dockerNetwork),
				Port: 6666 + counter,
			})
			counter++
		}

		topo.Servers = append(topo.Servers, ServerNode{
			Instances: instances,
			Key:       generatedKeysServer[i],
		})
	}

	for i := 0; i < numClients; i++ {
		topo.Clients = append(topo.Clients, ClientNode{
			Instance: Instance{
				Net:  getClientNetwork(dockerNetwork),
				Port: 5555 + i,
			},
			Balance: 100,
			Key:     generatedKeysClient[i],
		})
	}

	out, _ := json.Marshal(topo)
	return out
}

// topology optimized for full local testing (includes shards)
func localFull() []byte {
	numClients := 16
	numServers := 1
	numServerInstances := 3

	return config(numClients, numServers, numServerInstances, false)
}

// topology optimized for local testing
func localSimple() []byte {
	numClients := 3
	numServers := 1
	numServerInstances := 1

	return config(numClients, numServers, numServerInstances, false)
}

// topology optimized for local testing, extended
func localSimpleExtended() []byte {
	numClients := 16
	numServers := 1
	numServerInstances := 1

	return config(numClients, numServers, numServerInstances, false)
}

// topology optimized for aws
func aws() []byte {
	numClients := 16
	numServers := 1
	numInstances := 5

	return config(numClients, numServers, numInstances, false)
}

// topology optimized for docker
func dockerSimple() []byte {
	numClients := 16
	numServers := 1
	numInstances := 1

	return config(numClients, numServers, numInstances, true)
}
