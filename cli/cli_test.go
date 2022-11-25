package cli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestCli(t *testing.T) {
	t.Run("version test cli", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		cli := &CLI{outStream: outStream, errStream: errStream}
		args := strings.Split("vttreader --version", " ")
		status := cli.Run(args)
		if status != ExitCodeOk {
			t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOk)
		}
		expected := fmt.Sprintf("vttreader version %s\n", Version)
		if !strings.Contains(errStream.String(), expected) {
			t.Errorf("Output=%q, want %q", errStream.String(), expected)
		}
	})
	t.Run("version test cli", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		cli := &CLI{outStream: outStream, errStream: errStream}
		args := strings.Split("vttreader --file ../data/example.vtt", " ")
		status := cli.Run(args)
		if status != ExitCodeOk {
			t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOk)
		}
		expected := ""
		if !strings.Contains(errStream.String(), expected) {
			t.Errorf("Output=%q, want %q", errStream.String(), expected)
		}
	})

	t.Run("print json test cli", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		cli := &CLI{outStream: outStream, errStream: errStream}
		args := strings.Split("vttreader --file ../data/example.vtt --pj", " ")
		status := cli.Run(args)
		if status != ExitCodeOk {
			t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOk)
		}
		expected := ""
		if !strings.Contains(errStream.String(), expected) {
			t.Errorf("Output=%q, want %q", errStream.String(), expected)
		}
	})
	t.Run("save recreated vttfile test cli", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		cli := &CLI{outStream: outStream, errStream: errStream}
		args := strings.Split("vttreader --file ../data/example.vtt --path test.vtt ", " ")
		status := cli.Run(args)
		if status != ExitCodeOk {
			t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOk)
		}
		expected := ""
		if !strings.Contains(errStream.String(), expected) {
			t.Errorf("Output=%q, want %q", errStream.String(), expected)
		}
	})

}
