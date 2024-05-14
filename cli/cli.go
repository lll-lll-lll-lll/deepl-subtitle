package cli

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/lll-lll-lll-lll/vtt-formatter/vtt"
)

const Version string = "v0.2.0"

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
		version  bool
		filePath string
		path     string
		isPrint  bool
		help     bool // Add help flag
	)
	webVtt := &vtt.WebVtt{}

	flags := flag.NewFlagSet("vtt-formatter", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.Usage = func() {
		fmt.Fprintf(c.errStream, usage, "vtt-formatter")
	}
	flags.BoolVar(&version, "version", false, "print version")
	flags.StringVar(&filePath, "filepath", "", "path to the VTT file")
	flags.StringVar(&path, "path", "", "path to save the formatted VTT file")
	flags.BoolVar(&isPrint, "p", false, "print VTT elements in JSON format")
	flags.BoolVar(&help, "help", false, "print help")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	if help { // Check if help flag is set
		flags.Usage()
		return ExitCodeOk
	}

	if version {
		fmt.Fprintf(c.errStream, "vtt-formatter version %v\n", Version)
		return ExitCodeOk
	}
	if filePath != "" {
		if ext := filepath.Ext(filePath); ext != ".vtt" {
			fmt.Fprintf(c.errStream, "file extension is not vtt: %s\n", ext)
			return ExitCodeParseFlagError
		}
		vttfile, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer vttfile.Close()
		if _, err := webVtt.ReadFrom(vttfile); err != nil {
			log.Fatal(err)
		}
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
  -help or h                    Print help
  -version                      Print version
  -filepath=<{filename}.vtt>    Path to the VTT file
  -path=<{filename}.vtt>        Path to save the formatted VTT file
  -p                            Print VTT elements in JSON format
`
