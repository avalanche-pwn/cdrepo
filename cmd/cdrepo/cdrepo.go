package main

import (
	"flag"
	"github.com/avalanche-pwn/cdrepo/internal/core"
	"github.com/avalanche-pwn/cdrepo/internal/view"
)

func main() {
	register := flag.Bool("register",
		false, "Change directory and register if it's a git repo")
	flag.Parse()
	if *register {
		core.Register()
		return
	}
	view.Run()
}
