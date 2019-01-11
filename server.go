package main

import (
	"github.com/przemekBielak/blockchain"
	"github.com/davecgh/go-spew/spew" // pretty printing, install: go get -u github.com/davecgh/go-spew/spew
	"net/http"
	"fmt"
	"log"
	"reflect"
	"encoding/json"
)


var myBlockchain blockchain.Blockchain


func main() {

	// create genesis block 
	myBlockchain = append(myBlockchain, blockchain.Block{0, "genesis", "genesis", "genesis", "genesis"})

	http.HandleFunc("/addBlock", handlePost)
	http.HandleFunc("/getBlockchain", handleGet)

	fmt.Println("Serving on: http//localhost:7000/")
	log.Fatal(http.ListenAndServe("localhost:7000", nil))
}

func handlePost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	receivedData := req.Form.Get("data")

	fmt.Println("Raw data:", req.Form)
	fmt.Println("Received data:", receivedData)
	fmt.Println("Type of data:", reflect.TypeOf(receivedData))

	// add received data to blockchain
	err := myBlockchain.Append(receivedData)
	if err != nil {
		fmt.Println(err)
	}

	// range returns index, value pair. Index is ignored in this case
	for _, val := range myBlockchain {
		spew.Dump(val)
	}
}

func handleGet(w http.ResponseWriter, req *http.Request) {

	// create JSON from blockchain struct
	b, err := json.Marshal(myBlockchain)
	if err != nil {
		fmt.Println("error:", err)
	}

	blockchainJSONString := string(b)
	fmt.Fprintf(w, blockchainJSONString)
}