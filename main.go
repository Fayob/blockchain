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
	node := node.NewNode()

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal(err)
	}
	proto.RegisterNodeServer(grpcServer, node)
	fmt.Println("node running on port:", ":4000")

	go func() {
		for {
			time.Sleep(2 * time.Second)
			makeTransaction()
		}
	}()

	grpcServer.Serve(ln)
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
	}

	_, err = nodeClient.Handshake(context.Background(), version)
	if err != nil {
		log.Fatal(err)
	}
	
}