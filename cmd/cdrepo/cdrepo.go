package main

import (
	"flag"
	"github.com/avalanche-pwn/cdrepo/internal/core"
)

func main() {
	flag.BoolFunc("register",
		"Change directory and register if it's a git repo", core.Register)
	flag.Parse()
}
