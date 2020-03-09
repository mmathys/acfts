package main

import (
	"acfts-client/model"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/braintree/manners"

	"net/http"
	_ "net/http/pprof"
)

func benchmarkSetup(srvURLs []string, numClusters, numClients, gOwner int) []model.Address {
	// Cluster number
	num := 0
	db = initDB(num)

	// Delete all data in a client
	deleteAll(db)

	// Delete all data in servers
	for _, serverURL := range srvURLs {
		url := serverURL + "/all"
		body := _delete(url)
		log.Printf("Delete: %v\n", string(body))
	}

	const basePort = 3000
	port := basePort + num
	cBase := "http://localhost"

	// Generate private keys
	myurl := cBase + ":" + strconv.Itoa(port)
	generateClients(numClients, myurl)

	// Use manners to stop the server every scenaio
	r := initRoute(db)
	go manners.ListenAndServe(":"+strconv.Itoa(port), r)

	// For pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:7000", nil))
	}()

	otherClients := getOtherCURLs(cBase, numClusters, num)
	fmt.Println("other clients")
	fmt.Println(otherClients)

	for _, other := range otherClients {
		url := other + "/output"
		body := _delete(url)
		log.Printf("Delete: %v\n", string(body))
	}

	collectOtherAddrs(otherClients)

	addrs := getAllAddrs()
	log.Println(addrs)

	// Make a genesis transaction
	owner := addrs[gOwner]
	createGenesis(srvURLs, owner, 1000000)

	return addrs
}

// Scenario1: 0 -> 1
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

// Scenario2: 0 <-> 1
func BenchmarkScenario2(b *testing.B) {
	serverURLs := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	addrs := benchmarkSetup(serverURLs, 1, 4, 0)
	defer db.Close()

	atxs1 := []generalTx{
		// {From: addrs[0], To: []model.Address{addrs[1]}, Amounts: []int{1000000}},
		// {From: addrs[1], To: []model.Address{addrs[0]}, Amounts: []int{1000000}},
		{From: addrs[0], To: []model.Address{addrs[1]}, Amounts: []int{1}},
		{From: addrs[1], To: []model.Address{addrs[0]}, Amounts: []int{1}},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		executeTxs(serverURLs, atxs1)
	}

	manners.Close()
}

// Scenario3: 0 <-> 1, 2 <-> 3
func BenchmarkScenario3(b *testing.B) {
	serverURLs := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	addrs := benchmarkSetup(serverURLs, 1, 4, 0)
	defer db.Close()

	atx := generalTx{
		From:    addrs[0],
		To:      []model.Address{addrs[0], addrs[2]},
		Amounts: []int{500000, 500000},
	}
	executeTxs(serverURLs, []generalTx{atx})

	atxs1 := []generalTx{
		{From: addrs[0], To: []model.Address{addrs[1]}, Amounts: []int{500000}},
		{From: addrs[1], To: []model.Address{addrs[0]}, Amounts: []int{500000}},
		{From: addrs[2], To: []model.Address{addrs[3]}, Amounts: []int{500000}},
		{From: addrs[3], To: []model.Address{addrs[2]}, Amounts: []int{500000}},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		executeTxs(serverURLs, atxs1)
	}

	manners.Close()
}

// Scenario4: random -> random
func BenchmarkScenario4(b *testing.B) {
	serverURLs := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	addrs := benchmarkSetup(serverURLs, 1, 4, 0)
	defer db.Close()

	atx := generalTx{
		From:    addrs[0],
		To:      []model.Address{addrs[0], addrs[1], addrs[2], addrs[3]},
		Amounts: []int{250000, 250000, 250000, 250000},
	}
	executeTxs(serverURLs, []generalTx{atx})

	rand.Seed(time.Now().UnixNano())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		from := addrs[rand.Intn(4)]
		to := addrs[rand.Intn(4)]
		atxs1 := []generalTx{
			{From: from, To: []model.Address{to}, Amounts: []int{1}},
		}
		b.StartTimer()

		executeTxs(serverURLs, atxs1)
	}

	manners.Close()
}

// Scenario5: 0 -> 4
func BenchmarkScenario5(b *testing.B) {
	serverURLs := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	addrs := benchmarkSetup(serverURLs, 2, 4, 0)
	defer db.Close()

	atxs1 := []generalTx{
		{From: addrs[0], To: []model.Address{addrs[4]}, Amounts: []int{1}},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		executeTxs(serverURLs, atxs1)
	}

	manners.Close()
}
