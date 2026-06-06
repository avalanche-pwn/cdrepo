package main

import (
	"flag"
	"github.com/avalanche-pwn/cdrepo/internal/cdrepo"
)

func main() {
	flag.BoolFunc("register",
		"Change directory and register if it's a git repo", cdrepo.Register)
	flag.Parse()
	var tree cdrepo.BKTree
	tree.Read("test")
	for _, a := range tree.Search("files") {
		println(a)
	}
	tree.Save("test")
}
