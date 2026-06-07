package core

import (
	"os"
	"path"
)

func isRepo(p string) bool {
	git_dir := path.Join(p, ".git")
	_, err := os.Stat(git_dir)
	return err == nil
}
