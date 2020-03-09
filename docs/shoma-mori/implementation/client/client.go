package main

import (
	"acfts-client/api"
	"acfts-client/model"
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	// mrand "math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "net/http/pprof"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func initDB(num int) *gorm.DB {
	setting := fmt.Sprintf("root:@tcp(127.0.0.1:3306)/acfts_client_%d?charset=utf8&parseTime=True&loc=Local", num)
	db, err := gorm.Open("mysql", setting)
	if err != nil {
		log.Println("failed to connect database")
		log.Panicln(err)
	}

	db.AutoMigrate(&model.Output{})
	db.AutoMigrate(&model.Signature{})
	db.AutoMigrate(&model.Cluster{})
	db.AutoMigrate(&model.Client{})

	// Add a index to improve throughput of queries
	db.Model(&model.Output{}).
		AddIndex("idx_address_hash", "address1", "address2", "previous_hash")
	db.Model(&model.Signature{}).AddIndex("idx_signature", "output_id")

	return db
}

// delete all data in DB
func deleteAll(db *gorm.DB) {
	output := model.Output{}
	signature := model.Signature{}
	client := model.Client{}
	cluster := model.Cluster{}

	db.DropTable(&output)
	db.DropTable(&signature)
	db.DropTable(&client)
	db.DropTable(&cluster)

	db.AutoMigrate(&output)
	db.AutoMigrate(&signature)
	db.AutoMigrate(&client)
	db.AutoMigrate(&cluster)

	db.Model(&model.Output{}).
		AddIndex("idx_address_hash", "address1", "address2", "previous_hash")
	db.Model(&model.Signature{}).AddIndex("idx_signature", "output_id")
}

func initRoute(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// APIs
	r.GET("/address", api.GetAddrs(db))
	r.POST("/output", api.ReceiveUTXO(db))
	r.DELETE("/output", api.ClearOutputs(db))

	return r
}

func get(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicln(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	return body
}

// Do post request asynchronously
func post(url, jsonStr string, b chan []byte) {
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		log.Panicln(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	b <- body
}

// Delete all data in database
func _delete(url string) []byte {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Panicln(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	return body
}

func getClientSig(utxo model.Output) (string, string) {
	// Convert addresses to bytes to get its hash
	buf := []byte(fmt.Sprintf("%v%v%v", utxo.Address1, utxo.Address2, utxo.PreviousHash))

	// Get hash using SHA256
	h := crypto.Hash.New(crypto.SHA256)
	h.Write(buf)
	hashed := h.Sum(nil)

	// Get signature using ellipse curve cryptography
	pub := utxo.Address1 + utxo.Address2
	pri, ok := pub2Pri[pub]
	if !ok {
		log.Panicf("A public key %s is not in this cluster.\n", pub)
	}
	r, s, err := ecdsa.Sign(rand.Reader, pri, hashed)
	if err != nil {
		log.Panicln(err)
	}

	return r.String(), s.String()
}

func getServerSigs(utxo model.Output) string {
	var signatures []model.Signature
	db.Table("signatures").Where("output_id = ?", utxo.ID).Find(&signatures)

	str := ""
	for _, signature := range signatures {
		s := `
					{
						"address1": "` + signature.Address1 + `",
						"address2": "` + signature.Address2 + `",
						"signature1": "` + signature.Signature1 + `",
						"signature2": "` + signature.Signature2 + `"
					},`
		str += s
	}
	// Remove the last ','
	if len(str) > 0 {
		str = "[" + str[:len(str)-1] + "]"
	} else {
		str = "[]"
	}

	return str
}

func getSiblings(utxo model.Output) string {
	var siblings []model.Output
	db.Where("id <> ? AND previous_hash = ?", utxo.ID, utxo.PreviousHash).Find(&siblings)

	str := ""
	for _, sibling := range siblings {
		// used and signatures are initial values because it is when they are signed.
		s := `
				{
					"amount": ` + strconv.Itoa(sibling.Amount) + `,
					"address1": "` + sibling.Address1 + `",
					"address2": "` + sibling.Address2 + `",
					"previous_hash": "` + sibling.PreviousHash + `",
					"index": ` + strconv.Itoa(int(sibling.Index)) + `,
					"used": false,
					"server_signatures": []
				},`
		str += s
	}
	// Remove the last ','
	if len(str) > 0 {
		str = "[" + str[:len(str)-1] + "]"
	} else {
		str = "[]"
	}

	return str
}

func createInputStr(utxos []model.Output) string {
	inputs := ""
	for _, utxo := range utxos {
		sig1, sig2 := getClientSig(utxo)
		inputStr := `
		{
			"utxo": {
				"address1": "` + utxo.Address1 + `",
				"address2": "` + utxo.Address2 + `",
				"previous_hash": "` + utxo.PreviousHash + `",
				"index": ` + strconv.Itoa(int(utxo.Index)) + `,
				"server_signatures": ` + getServerSigs(utxo) + `
			},
			"siblings": ` + getSiblings(utxo) + `,
			"signature1": "` + sig1 + `",
			"signature2": "` + sig2 + `"
		},`

		inputs += inputStr
	}
	// Remove the last ','
	if len(inputs) > 0 {
		inputs = "[" + inputs[:len(inputs)-1] + "]"
	} else {
		inputs = "[]"
	}

	return inputs
}

func createSigStr(sigs []model.Signature) string {
	str := ""
	for _, sig := range sigs {
		s := `
					{
						"address1": "` + sig.Address1 + `",
						"address2": "` + sig.Address2 + `",
						"signature1": "` + sig.Signature1 + `",
						"signature2": "` + sig.Signature2 + `"
					},`
		str += s
	}
	// Remove the last ','
	if len(str) > 0 {
		str = "[" + str[:len(str)-1] + "]"
	} else {
		str = "[]"
	}

	return str
}

func createOutputStr(ops []model.Output, hash string, sigs []model.Signature) string {
	outputs := ""
	for _, op := range ops {
		outputStr := `
		{
			"amount": ` + strconv.Itoa(op.Amount) + `,
			"address1": "` + op.Address1 + `",
			"address2": "` + op.Address2 + `",
			"previous_hash": "` + hash + `",
			"index": ` + strconv.Itoa(int(op.Index)) + `,
			"server_signatures": ` + createSigStr(sigs) + `
		},`

		outputs += outputStr
	}
	// Remove the last ','
	if len(outputs) > 0 {
		outputs = "[" + outputs[:len(outputs)-1] + "]"
	} else {
		outputs = "[]"
	}

	return outputs
}

func getPreviousHash(previous string) string {
	bytes := sha256.Sum256([]byte(previous))
	num := fmt.Sprintf("%x", bytes)
	return string(num)
}

// From, To: Number of clients
type insideTx struct {
	From    int
	To      []int
	Amounts []int
}

type generalTx struct {
	From    model.Address
	To      []model.Address
	Amounts []int
}

func convertInsideTxs(its []insideTx) []generalTx {
	ats := make([]generalTx, len(its))
	for i, it := range its {
		fromPri := keys[it.From]
		fromPub := &fromPri.PublicKey
		from1 := fromPub.X.String()
		from2 := fromPub.Y.String()
		from := model.Address{from1, from2}

		to := make([]model.Address, len(it.To))
		for j, v := range it.To {
			toPri := keys[v]
			toPub := &toPri.PublicKey
			t1 := toPub.X.String()
			t2 := toPub.Y.String()
			t := model.Address{t1, t2}
			to[j] = t
		}

		at := generalTx{
			From:    from,
			To:      to,
			Amounts: it.Amounts,
		}
		ats[i] = at
	}

	return ats
}

func findUTXOs(addr model.Address, amount int) ([]model.Output, int) {
	var candidates []model.Output
	address1 := addr.Address1
	address2 := addr.Address2

	// TODO: signatures > 2/3のUTXOをSQLで見つける
	db.Where("address1 = ? AND address2 = ? AND used = false", address1, address2).
		Find(&candidates)

	// Select outputs which have enough signatures
	utxos := make([]model.Output, 0)
	for _, candidate := range candidates {
		count := 0
		db.Table("signatures").Where("output_id = ?", candidate.ID).Count(&count)
		if float64(count) >= 2.0*float64(n)/3.0 {
			utxos = append(utxos, candidate)
		}
	}
	if len(utxos) == 0 {
		return nil, 0
	}

	// Collect transactions so that the amount of them becomes equal to the arguement
	// FIXME: utxosをsortした方がいいかもしれない
	sum := 0
	for i, utxo := range utxos {
		sum += utxo.Amount
		if sum >= amount {
			utxos = utxos[:i+1]
			for _, utxo := range utxos {
				// FIXME: serverにこのUTXOが拒否されたら整合性が失われる
				db.Model(&utxo).Update("used", true)
			}
			return utxos, sum
		}
	}

	return nil, 0
}

func createJSONStr(tx generalTx) (string, error) {
	pubKey := tx.From
	want := 0
	for _, val := range tx.Amounts {
		want += val
	}
	utxos, sum := findUTXOs(pubKey, want)
	if sum == 0 {
		// TODO: このclient番号に紐づいている有効なUTXOがなかった場合
		// トランザクションの作成を保留にすべき？
		return "", fmt.Errorf("Error: there are no enough utxos of client %v", tx.From)
	}
	inputs := createInputStr(utxos)
	hash := getPreviousHash(inputs)

	ops := make([]model.Output, len(tx.To))
	for i, addr := range tx.To {
		ops[i] = model.Output{
			Amount: tx.Amounts[i],
			Address: model.Address{
				Address1: addr.Address1,
				Address2: addr.Address2,
			},
			PreviousHash: hash,
			Index:        uint(i),
		}
	}
	// Create a change transaction if the sum of utxos exceeds want
	if sum > want {
		change := model.Output{
			Amount: sum - want,
			Address: model.Address{
				Address1: tx.From.Address1,
				Address2: tx.From.Address2,
			},
			PreviousHash: hash,
			Index:        uint(len(ops)),
		}
		ops = append(ops, change)
	}
	outputs := createOutputStr(ops, hash, nil)

	jsonStr := `
{
	"inputs": ` + inputs + `,
	"outputs": ` + outputs + `
}`

	return jsonStr, nil
}

// FIXME: この構造体いらない？
type genesisJSON struct {
	Message string       `json:"message"`
	Genesis model.Output `json:"genesis"`
}

func createGenesis(urls []string, owner model.Address, amount int) {
	// Dummy signatures for genesis
	sig := model.Signature{
		Address: model.Address{
			Address1: "gene",
			Address2: "sis",
		},
		Signature1: "dum",
		Signature2: "my",
		OutputID:   1,
	}
	sigs := make([]model.Signature, n)
	for i := 0; i < n; i++ {
		sigs[i] = sig
	}

	genesis := model.Output{
		Amount: amount,
		Address: model.Address{
			Address1: owner.Address1,
			Address2: owner.Address2,
		},
		PreviousHash: "genesis",
		Index:        0,
		Used:         false,
		Signatures:   sigs,
	}
	db.Create(&genesis)

	jsonStr := `
	{
		"amount": ` + strconv.Itoa(amount) + `,
		"address1": "` + owner.Address1 + `",
		"address2": "` + owner.Address2 + `"
	}`

	// Create genesis records in servers
	b := make(chan []byte, n)
	for _, surl := range urls {
		url := surl + "/genesis"
		go post(url, jsonStr, b)
	}

	// Wait for responses from all servers
	for i := 0; i < n; i++ {
		body := <-b
		var g genesisJSON
		err := json.Unmarshal(body, &g)
		if err != nil {
			log.Panicln(err)
		}
		log.Println(g.Message)
	}
}

func generateClients(num int, url string) {
	keys = make([]*ecdsa.PrivateKey, num)
	pub2Pri = make(map[string]*ecdsa.PrivateKey)

	cluster := model.Cluster{URL: url}
	db.Create(&cluster)

	for i := 0; i < num; i++ {
		var err error
		keys[i], err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			log.Panicln(err)
		}
		pub := &keys[i].PublicKey
		pub2Pri[pub.X.String()+pub.Y.String()] = keys[i]

		client := model.Client{
			Address: model.Address{
				Address1: pub.X.String(),
				Address2: pub.Y.String(),
			},
		}
		db.Model(&cluster).Association("Clients").Append(client)
	}
}

func getOtherCURLs(base string, numClusters, num int) []string {
	otherClients := make([]string, 0)
	for i := 0; i < numClusters; i++ {
		if i != num {
			cport := 3000 + i
			otherClients = append(otherClients, base+":"+strconv.Itoa(cport))
		}
	}

	return otherClients
}

// FIXME: サーバー側もstruct Addressを作ったらaddr1, addr2をまとめる
func updateOutputs(outputs []model.Output, sigs []model.Signature) {
	for _, output := range outputs {
		// Make a record if the output is not created yet
		cnt := 0
		db.Where("address1 = ? AND address2 = ? AND previous_hash = ? AND output_index = ?",
			output.Address1, output.Address2, output.PreviousHash, output.Index).
			First(&output).Count(&cnt)
		if cnt == 0 {
			db.Create(&output)
		}
		db.Model(&output).Association("Signatures").Append(sigs)
	}
}

func executeTxs(serverURLs []string, txs []generalTx) {
	for i := 0; i < len(txs); i++ {
		jsonStr, err := createJSONStr(txs[i])
		if err != nil {
			log.Println(err)
			return
		}
		// fmt.Println(jsonStr)

		// Send requests to all servers
		b := make(chan []byte, n)
		for _, serverURL := range serverURLs {
			url := serverURL + "/transaction"
			go post(url, jsonStr, b)
		}

		var outputs []model.Output
		sigs := make([]model.Signature, n)
		// Wait for responses from all servers
		for j := 0; j < n; j++ {
			body := <-b
			response := model.Response{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				log.Panicln(err)
			}
			log.Println("Message: " + response.Message)

			outputs = response.Transaction.Outputs
			sig := model.Signature{
				Address: model.Address{
					Address1: response.Address1,
					Address2: response.Address2,
				},
				Signature1: response.Signature1.String(),
				Signature2: response.Signature2.String(),
			}
			sigs[j] = sig
		}

		if len(outputs) == 0 {
			log.Panicln("Error: there are no outputs.")
		}

		related := make(map[uint]bool)
		for _, output := range outputs {
			// Update signatures
			db.Where("address1 = ? AND address2 = ? AND previous_hash = ?",
				output.Address1, output.Address2, output.PreviousHash).First(&output)

			// Get a related client
			receiver := model.Client{}
			raddr1 := output.Address1
			raddr2 := output.Address2
			db.Where("address1 = ? AND address2 = ?", raddr1, raddr2).First(&receiver)
			related[receiver.ClusterID] = true
		}

		// log.Printf("related %v\n", related)

		r := make(chan []byte, len(related))
		cnt := 0
		for cid := range related {
			// If myself is related, add outputs to db
			// If else, send outputs to different clusters
			if cid == 1 {
				updateOutputs(outputs, sigs)
			} else {
				cluster := model.Cluster{}
				db.Where("id = ?", cid).First(&cluster)
				url := cluster.URL + "/output"

				hash := outputs[0].PreviousHash
				outputs := createOutputStr(outputs, hash, sigs)
				jsonStr := `
{
	"outputs": ` + outputs + `
}`

				go post(url, jsonStr, r)
				cnt++
			}
		}

		// Wait for responses from all servers
		for j := 0; j < cnt; j++ {
			body := <-r
			response := struct {
				Message string         `json:"message"`
				UTXOs   []model.Output `json:"utxos"`
			}{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				log.Panicln(err)
			}
			log.Printf("Message: %s\n", response.Message)
			// log.Println("received related utxos")
			// jsonBytes, _ := json.MarshalIndent(response.UTXOs, "", "    ")
			// log.Println(string(jsonBytes))
		}
	}
}

func collectOtherAddrs(others []string) {
	for _, ourl := range others {
		url := ourl + "/address"
		body := get(url)

		response := struct {
			Message string          `json:"message"`
			Addrs   []model.Address `json:"addresses"`
		}{}
		err := json.Unmarshal(body, &response)
		if err != nil {
			log.Panicln(err)
		}
		log.Printf("Message: %s\n", response.Message)
		log.Printf("Addrs: %s\n", response.Addrs)

		cluster := model.Cluster{URL: ourl}
		db.Create(&cluster)

		addrs := response.Addrs
		for _, addr := range addrs {
			client := model.Client{Address: addr}
			db.Model(&cluster).Association("Clients").Append(client)
		}
	}
}

func getAllAddrs() []model.Address {
	clients := []model.Client{}
	db.Find(&clients)
	addrs := make([]model.Address, len(clients))
	for i, client := range clients {
		addrs[i] = client.Address
	}

	return addrs
}

var db *gorm.DB

// FIXME: DBに秘密鍵と公開鍵を保存する
var keys []*ecdsa.PrivateKey

// Get a private key from a public key (PublicKey.X + PublicKey.Y)
var pub2Pri map[string]*ecdsa.PrivateKey

// n is the number of servers
const n = 4

func main() {
	serverURLs := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	fmt.Printf("Input cluster number: ")
	var num int
	fmt.Scan(&num)
	db = initDB(num)
	defer db.Close()

	deleteAll(db)

	const basePort = 3000
	port := basePort + num
	cBase := "http://localhost"

	// Generate private keys
	const numClients = 4
	myurl := cBase + ":" + strconv.Itoa(port)
	generateClients(numClients, myurl)

	/* For pprof */
	go func() {
		log.Println(http.ListenAndServe("localhost:7000", nil))
	}()

	r := initRoute(db)
	go r.Run(":" + strconv.Itoa(port))
	// r.Run(":" + strconv.Itoa(port))

	const numClusters = 1
	otherClients := getOtherCURLs(cBase, numClusters, num)

	fmt.Println("Input something when all clusters have been registered by all servers.")
	var dummy string
	fmt.Scan(&dummy)

	collectOtherAddrs(otherClients)

	addrs := getAllAddrs()
	log.Println(addrs)

	/* Make the genesis */
L_FOR:
	for {
		fmt.Println("Have a genesis?")
		fmt.Println("1. yes")
		fmt.Println("2. no")
		var g int
		fmt.Scan(&g)
		switch g {
		case 1:
			owner := addrs[0]
			createGenesis(serverURLs, owner, 200)
			break L_FOR
		case 2:
			break L_FOR
		default:
			fmt.Println("Please input valid number.")
		}
	}

	/* Make sample transactions */

	/* Use generalTx */

	// Inside one cluster
	atxs := []generalTx{
		{From: addrs[0], To: []model.Address{addrs[1]}, Amounts: []int{200}},
	}
	executeTxs(serverURLs, atxs)

	// atxs := make([]generalTx, 200)
	// tx := generalTx{From: addrs[0], To: []model.Address{addrs[1]}, Amounts: []int{1}}
	// for i := 0; i < 200; i++ {
	// 	atxs[i] = tx
	// }
	// executeTxs(serverURLs, atxs)

	// Between two clusters
	// atxs := []generalTx{
	// 	{From: addrs[0], To: []model.Address{addrs[4]}, Amounts: []int{200}},
	// }
	// executeTxs(serverURLs, atxs)

	// atxs := make([]generalTx, 200)
	// tx := generalTx{From: addrs[0], To: []model.Address{addrs[4]}, Amounts: []int{1}}
	// for i := 0; i < 200; i++ {
	// 	atxs[i] = tx
	// }
	// executeTxs(serverURLs, atxs)

	/* Use insideTx */

	// case 1: 0 -> 1
	// itxs := []insideTx{
	// 	{From: 0, To: []int{1}, Amounts: []int{200}},
	// }
	// atxs := convertInsideTxs(itxs)
	// executeTxs(serverURLs, atxs)

	// case 2: 0 <-> 1
	// tx1 := insideTx{From: 0, To: []int{1}, Amounts: []int{200}}
	// tx2 := insideTx{From: 1, To: []int{0}, Amounts: []int{200}}
	// itxs := []insideTx{}
	// for i := 0; i < 10; i++ {
	// 	itxs = append(itxs, tx1)
	// 	itxs = append(itxs, tx2)
	// }
	// atxs := convertInsideTxs(itxs)
	// executeTxs(serverURLs, atxs)

	// case 3: 0 <-> 1 & 2 <-> 3
	// tx0 := insideTx{From: 0, To: []int{0, 2}, Amounts: []int{50, 150}}
	// itxs0 := []insideTx{}
	// itxs0 = append(itxs0, tx0)
	// atxs0 := convertInsideTxs(itxs0)
	// executeTxs(serverURLs, atxs0)
	//
	// tx1 := insideTx{From: 0, To: []int{1}, Amounts: []int{50}}
	// tx2 := insideTx{From: 1, To: []int{0}, Amounts: []int{50}}
	// tx3 := insideTx{From: 2, To: []int{3}, Amounts: []int{150}}
	// tx4 := insideTx{From: 3, To: []int{2}, Amounts: []int{150}}
	// itxs1 := []insideTx{}
	// for i := 0; i < 10; i++ {
	// 	itxs1 = append(itxs1, tx1)
	// 	itxs1 = append(itxs1, tx2)
	// 	itxs1 = append(itxs1, tx3)
	// 	itxs1 = append(itxs1, tx4)
	// }
	// atxs1 := convertInsideTxs(itxs1)
	// executeTxs(serverURLs, atxs1)

	// case 4: random -> random
	// mrand.Seed(time.Now().UnixNano())
	// tx0 := insideTx{From: 0, To: []int{0, 1, 2, 3}, Amounts: []int{50, 50, 50, 50}}
	// itxs0 := []insideTx{}
	// itxs0 = append(itxs0, tx0)
	// atxs0 := convertInsideTxs(itxs0)
	// executeTxs(serverURLs, atxs0)
	//
	// itxs1 := []insideTx{}
	// for i := 0; i < 25; i++ {
	// 	from := mrand.Intn(4)
	// 	to := mrand.Intn(4)
	// 	tx1 := insideTx{From: from, To: []int{to}, Amounts: []int{1}}
	// 	itxs1 = append(itxs1, tx1)
	// }
	// atxs1 := convertInsideTxs(itxs1)
	// executeTxs(serverURLs, atxs1)

	// Wait for receiving outputs of this cluster
	// for {
	// }
}
