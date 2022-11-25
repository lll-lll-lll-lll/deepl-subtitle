package main

import (
	"os"
)

func main() {
	cli := cli.NewCLI(os.Stdout, os.Stderr)
	os.Exit(cli.Run(os.Args))
}
