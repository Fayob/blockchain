package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fayob/blockchain/node"
	"github.com/fayob/blockchain/proto"
	"google.golang.org/grpc"
)

func main()  {
	makeNode(":3000", []string{})
	makeNode(":4000", []string{":3000"})

	// go func() {
	// 	for {
	// 		time.Sleep(2 * time.Second)
	// 		makeTransaction()
	// 	}
	// }()

	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()
	go n.Start(listenAddr)
	if len(bootstrapNodes) > 0 {
		if err := n.BootstrapNetwork(bootstrapNodes); err != nil {
			log.Fatal(err)
		}
	}
	return n
}

func makeTransaction()  {
	client, err := grpc.Dial(":4000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	nodeClient := proto.NewNodeClient(client)
	
	version := &proto.Version{
		Version: "blocker-0.1",
		Height: 12,
		ListenAddr: ":4001",
	}

	_, err = nodeClient.Handshake(context.Background(), version)
	if err != nil {
		log.Fatal(err)
	}
	
}