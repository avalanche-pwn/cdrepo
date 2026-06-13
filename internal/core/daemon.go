package core

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	pb "github.com/avalanche-pwn/cdrepo/internal/daemon_pb"
	"github.com/avalanche-pwn/cdrepo/internal/invidx"
	"github.com/avalanche-pwn/cdrepo/internal/searchif"

	"google.golang.org/grpc"
)

type daemon struct {
	pb.UnimplementedDaemonServer
	search    searchif.FuzzySearcher
	grpcSrv   *grpc.Server
	keepAlive chan bool
}

type stringSearch string

func (s stringSearch) Key() string {
	return string(s)
}

func searchFactory() searchif.FuzzySearcher {
	searcher := invidx.InvIdx{}
	searcher.Init(nil)
	return &searcher
}

func (d *daemon) register(path string) {
	if isRepo(path) {
		fmt.Printf("Adding path %s\n", path)
		d.search.Add(stringSearch(path))
	}
}

func init() {
	gob.Register(stringSearch(""))
}

func read(s string) searchif.FuzzySearcher {
	f, _ := os.Open(s)
	defer f.Close()
	dec := gob.NewDecoder(f)
	var test searchif.FileFormat
	err := dec.Decode(&test)
	if err != nil || test.Version != searchif.FileVersion {
		println ("Invalid bin file")
	}
	return test.Value.Decode()
}

func save(s string, searcher searchif.FuzzySearcher) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)

	enc.Encode(searchif.FileFormat{Value: searcher.Encode(),
		Version: searchif.FileVersion})
	os.WriteFile(s, network.Bytes(), 0644)
}

func (d *daemon) init() {
	d.search = searchFactory()
	binDir := binPath()
	if _, err := os.Stat(binDir); err == nil {
		d.search = read(binDir)
		fmt.Printf("Read search from %s\n", binDir)
	}
	d.keepAlive = make(chan bool)
}

func (d *daemon) stop() {
	save(binPath(), d.search)
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
		res_node := val.Value.(stringSearch)
		results[i] = &pb.SearchResult{
			Score: int32(val.Score), Value: string(res_node)}
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
