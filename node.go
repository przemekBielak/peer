package main

import (
	"github.com/przemekBielak/blockchain"
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
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
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
		host.SetStreamHandler("/chat/1.0.0", handleStream)

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

		fmt.Printf("Run './chat -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, host.ID().Pretty())
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

		s, err := host.NewStream(context.Background(), peerInfo.ID, "/chat/1.0.0")
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

	fmt.Println(myBlockchain)

	createHost()
}


