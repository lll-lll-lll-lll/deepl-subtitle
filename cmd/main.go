package main

import (
	ds "github.com/lll-lll-lll-lll/deepl-subtitle"
	"os"
)

func main() {
	cli := ds.NewCLI(os.Stdout, os.Stderr)
	os.Exit(cli.Run(os.Args))
}
