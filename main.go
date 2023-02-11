package main

import (
	"os"

	"github.com/lll-lll-lll-lll/format-webvtt/cli"
)

func main() {
	cli := cli.New(os.Stdout, os.Stderr)
	os.Exit(cli.Run(os.Args))
}
