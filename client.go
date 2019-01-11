package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"encoding/json"
	"bytes"
)

const serverAddress string = "http://localhost:7000"

func postToBlockchain(data string) {
	resp, err := http.PostForm(serverAddress + "/addBlock", url.Values{"data": {data}})
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func getBlockchain() {
	resp, err := http.Get(serverAddress + "/getBlockchain")
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// indent JSON file to be easy to read
	var out bytes.Buffer
	json.Indent(&out, body, "-->", "\t")

	out.WriteTo(os.Stdout)
}

func main() {
	// Join command line arguments as a block data
	data := strings.Join(os.Args[1:], " ")

	postToBlockchain(data)
	getBlockchain()
}	
