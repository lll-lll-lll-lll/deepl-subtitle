package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"

	ds "github.com/lll-lll-lll-lll/deepl-subtitle"
)

const (
	ExitCodeOk             = 0
	ExitCodeParseFlagError = 1
)

type CLI struct {
	outStream, errStream io.Writer
}

func NewCLI(outStream io.Writer, errStream io.Writer) *CLI {
	return &CLI{outStream: outStream, errStream: errStream}
}

func (c *CLI) Run(args []string) int {
	var (
		version bool
		file    string
		vttfile ds.WebVttString
		path    string
		err     error
	)

	flags := flag.NewFlagSet("vttreader", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.Usage = func() {
		fmt.Fprintf(c.errStream, usage, "vttreader")
	}
	flags.BoolVar(&version, "version", false, "print version")
	flags.StringVar(&file, "file", "", "vtt file")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	if version {
		fmt.Fprintf(c.errStream, "vttreader version %v\n", Version)
		return ExitCodeOk
	}
	if file != "" {
		vttfile, err = ds.ReadVTTFile(file)
		if err != nil {
			log.Fatal(err)
		}
	}
	webVtt := ds.NewWebVtt(vttfile)
	webVtt.ScanLines(ds.ScanSplitFunc)
	w := ds.UnifyTextByTerminalPoint(webVtt)
	vtt := ds.DeleteVTTElementOfEmptyText(w)

	if path != "" {
		vtt.ToFile(path)
		return ExitCodeOk
	}

	defer ds.PrintlnJson(vtt.VttElements)

	return ExitCodeOk
}

func printJsonAny[S ~[]e, e any](c *CLI, f S) {
	for _, e := range f {
		b, _ := json.Marshal(e)
		var out bytes.Buffer
		if err := json.Indent(&out, b, "", "  "); err != nil {
			fmt.Fprintf(c.errStream, "filter err %s\n", err)
		}
		fmt.Fprintf(c.errStream, "%s\n", out.String())
	}
}

const usage = `
Usage: %s [options] slug path
  
Options:
  -help or h 	 		    help
  -version            		now version
  -file=<option>      		vtt file name
`
