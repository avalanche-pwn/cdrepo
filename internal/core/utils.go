package core

import (
	"context"
	"fmt"
	"os"
	"path"
	"syscall"
	"time"

	pb "github.com/avalanche-pwn/cdrepo/internal/daemon_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const configPath string = "cdrepo"
const searchFile string = "cdrepo.bin"
const lockPath string = "/tmp/cdrepo.lock"
const daemonCheckTimeout time.Duration = 10 * time.Second
const daemonTimeout time.Duration = 60 * time.Second

type FuzzySearcher interface {
	Add(s string)
	Read(s string)
	Save(s string)
	Search(s string) []string
}

func init() {
	home := os.Getenv("HOME")
	fullConfPath := path.Join(home, ".config", configPath)
	os.Mkdir(fullConfPath, os.ModePerm)
}

func binPath() string {
	home := os.Getenv("HOME")
	return path.Join(home, ".config", configPath, searchFile)
}

func requestRegister() {
	conn, err := grpc.NewClient("unix:/tmp/cdrepo.sock",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
	}
	defer conn.Close()
	c := pb.NewDaemonClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	p, _ := os.Getwd()
	r, err := c.Register(ctx, &pb.RegisterRequest{Path: p})
	if err != nil {
		fmt.Printf("could not greet: %v\n", err)
	}
	if r.Success != true {
		println("Sth broke")
	}
}

func waitWhileActive(keepAlive chan bool) {
	lastUpdate := time.Now()
	for time.Since(lastUpdate) < daemonTimeout {
		select {
		case <-keepAlive:
			lastUpdate = time.Now()
		case <-time.After(daemonCheckTimeout):
		}
	}
}

func Register(s string) error {
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		os.Exit(1)
	}
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		requestRegister()
		os.Exit(0)
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	defer f.Close()

	p, _ := os.Getwd()
	d := serve(p)
	waitWhileActive(d.keepAlive)
	d.stop()
	return nil
}
