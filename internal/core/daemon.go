package core

import (
	"net"
	"context"
	"google.golang.org/grpc"
	pb "github.com/avalanche-pwn/cdrepo/internal/daemon_pb"
)

type server struct {
	pb.UnimplementedDaemonServer
}

func (s *server) Register(_ context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{Success: true}, nil
}

func Serve() {
	lis, err := net.Listen("unix", "/tmp/cdrepo.sock")
	if err != nil {
		println("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDaemonServer(s, &server{})
	println("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		println("failed to serve: %v", err)
	}
}
