package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/lll-lll-lll-lll/deepl-subtitle/webvtt"
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
		vttfile webvtt.WebVttString
		path    string
		err     error
		pj      bool
	)
	webVtt := &webvtt.WebVtt{}

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
		vttfile, err = webvtt.ReadVTT(file)
		if err != nil {
			log.Fatal(err)
		}
		webVtt = webvtt.New(vttfile)
		webVtt.ScanLines(webvtt.ScanSplitFunc)
		w := webvtt.UnifyText(webVtt)
		webvtt.DeleteElementOfEmptyText(w)
	}

	if path != "" {
		webVtt.ToFile(path)
	}

	if pj {
		webvtt.PrintlnJson(webVtt.Elements)
	}

	return ExitCodeOk
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
