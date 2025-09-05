package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/fayob/blockchain/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version    string
	listenAddr string
	logger     *zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[proto.NodeClient]proto.Version

	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()
	return &Node{
		peers:   make(map[proto.NodeClient]proto.Version),
		version: "blocker-0.1",
		logger: logger.Sugar(),
	}
}

func (n *Node) addPeer(client proto.NodeClient, v proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.logger.Debugw("new peer connected", "addr", v.ListenAddr, "height", v.Height)
	// fmt.Printf("[%s] new peer connected (%s) - height (%d)\n", n.listenAddr, v.ListenAddr, v.Height)

	n.peers[client] = v
}

func (n *Node) deletePeer(client proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, client)
}

func (n *Node) BootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		c, err := makeNodeClient(addr)
		if err != nil {
			return err
		}

		version, err := c.Handshake(context.Background(), n.getVersion())
		if err != nil {
			n.logger.Error("handshake error:", err)
			// log.Println("handshake failed with peer:", addr)
			continue
		}

		n.addPeer(c, version)
	}

	return nil
}

func (n *Node) Start(listenAddr string) error {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	proto.RegisterNodeServer(grpcServer, n)

	fmt.Println("node running on port:", listenAddr)

	return grpcServer.Serve(ln)
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	c, err := makeNodeClient(v.ListenAddr)
	if err != nil {
		return nil, err
	}

	n.addPeer(c, *v)

	return n.getVersion(), nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	fmt.Println("received transaction from ", peer.Addr)
	return &proto.Ack{}, nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "blocker-0.1",
		Height:     0,
		ListenAddr: n.listenAddr,
	}
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	c, err := grpc.Dial(listenAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return proto.NewNodeClient(c), nil
}
