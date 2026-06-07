package main

import (
	"flag"
	"github.com/avalanche-pwn/cdrepo/internal/core"
	"github.com/avalanche-pwn/cdrepo/internal/bk_tree"
)

func main() {
	flag.BoolFunc("register",
		"Change directory and register if it's a git repo", core.Register)
	flag.Parse()
	var tree bk_tree.BKTree
	tree.Read("test")
	// tree.Add("dupa")
	// tree.Add("zephyr")
	for _, a := range tree.Search("zephyr") {
		println(a)
	}
	tree.Save("test")
}
