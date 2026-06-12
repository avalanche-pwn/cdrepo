package core

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/avalanche-pwn/cdrepo/internal/bk_tree"
	"github.com/avalanche-pwn/cdrepo/internal/searchif"
	pb "github.com/avalanche-pwn/cdrepo/internal/daemon_pb"

	"google.golang.org/grpc"
)

type daemon struct {
	pb.UnimplementedDaemonServer
	search    searchif.FuzzySearcher
	grpcSrv   *grpc.Server
	keepAlive chan bool
}

func searchFactory() searchif.FuzzySearcher {
	return &bk_tree.BKTree{}
}

func (d *daemon) register(path string) {
	if isRepo(path) {
		fmt.Printf("Adding path %s\n", path)
		d.search.Add(path)
	}
}

func (d *daemon) init() {
	d.search = searchFactory()
	binDir := binPath()
	if _, err := os.Stat(binDir); err == nil {
		d.search.Read(binDir)
		fmt.Printf("Read search from %s\n", binDir)
	}
	d.keepAlive = make(chan bool)
}

func (d *daemon) stop() {
	d.search.Save(binPath())
	d.grpcSrv.Stop()
}

func (d *daemon) Register(_ context.Context, in *pb.RegisterRequest) (
	*pb.RegisterResponse, error) {
	d.keepAlive <- true
	d.register(in.Path)
	return &pb.RegisterResponse{Success: true}, nil
}

func (d *daemon) Search(_ context.Context, in *pb.SearchRequest) (
	*pb.SearchResponse, error) {
	api_res := d.search.Search(in.Query)
	results := make([]*pb.SearchResult, len(api_res))

	for i, val := range api_res {
		results[i] = &pb.SearchResult{Score: int32(val.Score), Value: val.Value}
	}

	rsp := &pb.SearchResponse{Results: results}
	return rsp, nil
}

func serve(initialPath string) daemon {
	lis, err := net.Listen("unix", "/tmp/cdrepo.sock")
	if err != nil {
		println("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	var d daemon
	d.init()
	d.register(initialPath)
	d.grpcSrv = s

	pb.RegisterDaemonServer(s, &d)
	println("daemon listening at %v", lis.Addr())
	go s.Serve(lis)
	return d
}
