package cli

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/lll-lll-lll-lll/vtt-formatter/vtt"
)

const Version string = "v0.1.0"

const (
	ExitCodeOk             = 0
	ExitCodeParseFlagError = 1
)

type CLI struct {
	outStream, errStream io.Writer
}

func New(outStream io.Writer, errStream io.Writer) *CLI {
	return &CLI{outStream: outStream, errStream: errStream}
}

func (c *CLI) Run(args []string) int {
	var (
		version bool
		file    string
		vttfile string
		path    string
		err     error
		isPrint bool
	)
	webVtt := &vtt.WebVtt{}

	flags := flag.NewFlagSet("vttreader", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.Usage = func() {
		fmt.Fprintf(c.errStream, usage, "vttreader")
	}
	flags.BoolVar(&version, "version", false, "print version")
	flags.StringVar(&file, "file", "", "vtt file")
	flags.StringVar(&path, "path", "", "save path")
	flags.BoolVar(&isPrint, "p", false, "print vtt elements json format")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	if version {
		fmt.Fprintf(c.errStream, "vttreader version %v\n", Version)
		return ExitCodeOk
	}
	if file != "" {
		vttfile, err = vtt.ReadFileContents(file)
		if err != nil {
			log.Fatal(err)
		}
		webVtt = vtt.New(vttfile)
		webVtt.Format()
	}

	if path != "" {
		f, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if _, err := webVtt.WriteTo(f); err != nil {
			log.Fatal(err)
		}
	}

	if isPrint {
		webVtt.WriteTo(os.Stdout)
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
  -p                        print vtt elements json format
`
