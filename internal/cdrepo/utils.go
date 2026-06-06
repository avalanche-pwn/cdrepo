package cdrepo

import (
	"os"
	"path"
)
const configPath string = "cdrepo"

type FuzzySearcher interface {
	Add(s string)
	Read(s string)
	Save(s string)
	Search(s string)
}

func init() {
	home := os.Getenv("HOME");
	fullConfPath := path.Join(home, ".config", configPath)
	os.Mkdir(fullConfPath, os.ModePerm);
}

func Register(s string) error {
	return nil;
}
