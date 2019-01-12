package main

import (
	"github.com/przemekBielak/blockchain"
	"github.com/davecgh/go-spew/spew" 
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto" 
	ma "github.com/multiformats/go-multiaddr"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-net"
	"context"
	"fmt"
	"crypto/rand"
	"log"
	"flag"
	"bufio"
	"os"
	"encoding/json"
	"strings"
)


var myBlockchain blockchain.Blockchain

func handleStream(s net.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			newBlockchain := make(blockchain.Blockchain, 0)
			err := json.Unmarshal([]byte(str), &newBlockchain)
			if err != nil {
				fmt.Println("Parsing json failed")
			}

			blockchain.Verify(&myBlockchain, &newBlockchain)

			for _, val := range myBlockchain {
				spew.Dump(val)
			}
			fmt.Print("> ")
		}
	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("\x1b[32m> \x1b[0m")
		readData, err := stdReader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		// add read data to blockchain
		readData = strings.Replace(readData, "\n", "", -1)
		err = myBlockchain.Append(readData)
		if err != nil {
			fmt.Println(err)
		}

		// create JSON from blockchain struct
		b, err := json.Marshal(myBlockchain)
		if err != nil {
			fmt.Println("error:", err)
		}

		blockchainJSONString := string(b)

		rw.WriteString(fmt.Sprintf("%s\n", blockchainJSONString))
		rw.Flush()
	}

}

func createHost() {
	sourcePortPtr := flag.Int("sp", 0, "This peer source port")
	destAddrPtr := flag.String("d", "", "Destination multiaddr string")
	flag.Parse()

	// priv, pub, err
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		panic(err)
	}

	SrcMultiAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", *sourcePortPtr))
	if err != nil {
		panic(err)
	}

	// The context governs the lifetime of the libp2p node
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	host, err := libp2p.New(
		ctx,
		libp2p.Identity(priv),
		libp2p.ListenAddrs(SrcMultiAddr),
	)
	if err != nil {
		panic(err)
	}

	// If destination address is not provided, act as a host
	if *destAddrPtr == "" {
		host.SetStreamHandler("/blockchain/0.0.1", handleStream)

		// Get the actual TCP port from listen multiaddr
		// 0 - random available port
		// otherwise *sourcePortPtr from cli argument
		var port string
		for _, addr := range host.Addrs() {
			p, err := addr.ValueForProtocol(ma.P_TCP) 
			if err == nil {
				port = p
				break
			}
		}

		if port == "" {
			panic("was not able to find actual local port")
		}

		fmt.Printf("Run 'go run node.go -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, host.ID().Pretty())
		fmt.Println("You can replace 127.0.0.1 with public IP as well.")
		fmt.Printf("\nWaiting for incoming connection\n\n")

		// Hang forever
		<-make(chan struct{})
	} else {
		DstMultiAddr, err := ma.NewMultiaddr(*destAddrPtr)
		if err != nil {
			panic(err)
		}

		peerInfo, err := peerstore.InfoFromP2pAddr(DstMultiAddr)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("This nodes multiaddress:")
		for _, addr := range host.Addrs() {
			fmt.Println(addr)
		}
		fmt.Println(host.ID().Pretty());
		fmt.Println()
	
		fmt.Println("Peer multiaddress:")
		fmt.Println(peerInfo.Addrs)
		fmt.Println(peerInfo.ID.Pretty())

		host.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)

		s, err := host.NewStream(context.Background(), peerInfo.ID, "/blockchain/0.0.1")
		if err != nil {
			fmt.Println(err)
		}

		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		
		go writeData(rw)
		go readData(rw)

		select {}
	}	
}

func main() {

	// create genesis block 
	myBlockchain = append(myBlockchain, blockchain.Block{0, "genesis", "genesis", "genesis", "genesis"})

	createHost()
}


