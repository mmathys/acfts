package main

import (
	"encoding/json"
	"fmt"
	"github.com/mmathys/acfts/common"
	"io/ioutil"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	// Declare a request struct
	var body common.Tuple

	/*
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	 */

	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}

	err = json.Unmarshal(payload, &body)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	fmt.Fprintf(w, "thanks, %v\n", body)
}

func main() {
	a := []byte{0,0}
	out, _ := json.Marshal(a)
	err := json.Unmarshal(out, &a)
	if err != nil {
		fmt.Println("failed")
	}
	fmt.Println(string(out))

	http.HandleFunc("/sign", hello)
	http.ListenAndServe(":6666", nil)
}
