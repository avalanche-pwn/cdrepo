package core

import (
	"os"
	"path"
	"syscall"
)
const configPath string = "cdrepo"
const lockPath string = "/tmp/cdrepo.lock"

type FuzzySearcher interface {
	Add(s string)
	Read(s string)
	Save(s string)
	Search(s string) []string
}

func init() {
	home := os.Getenv("HOME");
	fullConfPath := path.Join(home, ".config", configPath)
	os.Mkdir(fullConfPath, os.ModePerm);
}

func Register(s string) error {
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		os.Exit(1)
	}
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		os.Exit(0)
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	defer f.Close()
	
	p, _ := os.Getwd()
	return serve(p)
}
