package core

import (
	"context"
	"net"
	"fmt"

	"github.com/avalanche-pwn/cdrepo/internal/bk_tree"
	pb "github.com/avalanche-pwn/cdrepo/internal/daemon_pb"
	"google.golang.org/grpc"
)

type daemon struct {
	pb.UnimplementedDaemonServer
	search FuzzySearcher
}

func searchFactory() FuzzySearcher {
	return bk_tree.BKTree{}
}

func (d *daemon) register(path string) {
	fmt.Printf("Adding path %s\n", path)
	d.search.Add(path)
}

func (d *daemon) init() {
	d.search = searchFactory()
}

func (d *daemon) Register(_ context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	d.register(in.Path)
	return &pb.RegisterResponse{Success: true}, nil
}

func serve(initialPath string) error {
	lis, err := net.Listen("unix", "/tmp/cdrepo.sock")
	if err != nil {
		println("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	var d daemon
	d.init()
	d.register(initialPath)

	pb.RegisterDaemonServer(s, &d)
	println("daemon listening at %v", lis.Addr())
	return s.Serve(lis)
}
