package main

import (
	"github.com/przemekBielak/blockchain"
	// "github.com/libp2p/go-libp2p-host"
	"fmt"
	"crypto/rand"
	"crypto/rsa"
	"golang.org/x/crypto/ed25519" 
)


var myBlockchain blockchain.Blockchain


func createHost(port int) {
	fmt.Println(port)

	// random number
	r := rand.Reader

	// generate private/public keys
	privateKey, err := rsa.GenerateKey(r, 2048)
	if err != nil {
		fmt.Println(err.Error)
	}
	publicKey := &privateKey.PublicKey

	fmt.Println(privateKey)
	fmt.Println(publicKey)

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	fmt.Println("priv:", priv)
	fmt.Println("pub:", pub)

	data := []byte{123, 123, 123, 124}
	sign := ed25519.Sign(priv, data)

	fmt.Println("signature:", sign)

	fakedSign := []byte{123, 231, 143}
	fmt.Println(ed25519.Verify(pub, data, fakedSign))
	
}

func main() {

	// create genesis block 
	myBlockchain = append(myBlockchain, blockchain.Block{0, "genesis", "genesis", "genesis", "genesis"})

	fmt.Println(myBlockchain)

	createHost(123)
}

