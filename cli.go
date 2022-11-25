package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
)

const Version string = "v0.1.0"

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
		vttfile WebVttString
		path    string
		err     error
		pj      bool
	)
	webVtt := &WebVtt{}

	flags := flag.NewFlagSet("vttreader", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.Usage = func() {
		fmt.Fprintf(c.errStream, usage, "vttreader")
	}
	flags.BoolVar(&version, "version", false, "print version")
	flags.StringVar(&file, "file", "", "vtt file")
	flags.StringVar(&path, "path", "", "save path")
	flags.BoolVar(&pj, "pj", false, "print vtt elements json format")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	if version {
		fmt.Fprintf(c.errStream, "vttreader version %v\n", Version)
		return ExitCodeOk
	}
	if file != "" {
		vttfile, err = ReadVTT(file)
		if err != nil {
			log.Fatal(err)
		}
		webVtt = New(vttfile)
		webVtt.ScanLines(ScanSplitFunc)
		w := UnifyText(webVtt)
		DeleteElementOfEmptyText(w)
	}

	if path != "" {
		webVtt.ToFile(path)
	}

	if pj {
		PrintlnJson(webVtt.Elements)
	}

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
  -file=<{filename}.vtt>    vtt file name
  -path=<{filename}.vtt>    shaped vtt file path
  -pj                       print json console
`
