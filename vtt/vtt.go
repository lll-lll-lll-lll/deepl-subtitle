package vtt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// WebVtt represents a WebVTT file.
type WebVtt struct {
	Elements []*Element     `json:"vtt_elements"` // Collection of vtt structure
	Header   *Header        `json:"header"`       // First two lines of the vtt file
	Scanner  *bufio.Scanner `json:"scanner"`      // Scanner for reading the file
	builder  *strings.Builder
	sumBytes int
}

// Element represents a pattern from the vtt file.
type Element struct {
	StartTime string `json:"start_time"` // Start time of the element
	EndTime   string `json:"end_time"`   // End time of the element
	Position  string `json:"position"`   // Position of the element
	Line      string `json:"line"`       // Line of the element
	Text      string `json:"text"`       // Text of the element
	Separator string `json:"separator"`  // Separator of the element
}

// Header represents the first two lines of the vtt file.
type Header struct {
	Head string `json:"head"` // First line of the header
	Note string `json:"note"` // Second line of the header
}

// scanSplitFunc is a split function for the scanner.
func scanSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	t := string(token)
	// CheckTimeRegexpFlagでtrueが走るとその行を空白で単語区切りにする。トークン区切りになった他の`-->`や`position...`を他のフラグで検索
	if checkToken(StartOrEndTime, t) || checkToken(Separator, t) || checkToken(Position, t) || checkToken(Line, t) {
		{
			advance, token, err = bufio.ScanWords(data, atEOF)
			return
		}
	}
	return
}

// unifyText updates EndTime by passing Text to the previous structure if it contains `.` or `?`.
func (wv *WebVtt) unifyText() {
	for i := 0; i < len(wv.Elements)-1; i++ {
		untilTerminalPointCnt := recursiveSearchTerminalPoint(wv.Elements, i)
		for j := untilTerminalPointCnt; j > i; j-- {
			wv.Elements[j-1].Text += " " + wv.Elements[j].Text
			wv.Elements[j-1].EndTime = wv.Elements[j].EndTime
			wv.Elements[j].Text = ""
		}
		// Move to the position where the token representing the end of the sentence is found
		if untilTerminalPointCnt > 0 {
			i = untilTerminalPointCnt
		}
	}
}

// deleteElementOfEmptyText deletes all structures with empty text.
func (wv *WebVtt) deleteElementOfEmptyText() {
	var i int
	for {
		if wv.Elements[i].Text == "" {
			wv.Elements = append(wv.Elements[:i], wv.Elements[i+1:]...)
			i--
		}
		i++
		if len(wv.Elements) == i {
			break
		}
	}
}

// recursiveSearchTerminalPoint recursively searches for the terminal point in the given slice of VTT elements.
// It returns the index of the element that contains the terminal point.
func recursiveSearchTerminalPoint(elements []*Element, untilTerminalCnt int) int {
	if untilTerminalCnt == len(elements)-1 {
		return untilTerminalCnt
	}
	e := elements[untilTerminalCnt].Text
	locs := searchTerminalToken(e)
	f := func(locs []int) bool {
		return len(locs) == 0
	}
	if f(locs) {
		untilTerminalCnt++
		return recursiveSearchTerminalPoint(elements, untilTerminalCnt)
	}
	return untilTerminalCnt
}

// scanLines reads the file line by line and creates the structure.
func (wv *WebVtt) scanLines(splitFunc bufio.SplitFunc) {
	e := &Element{}
	wv.Scanner.Split(splitFunc)
	var isStartOrEndTime int
	var sumBytes int

	for wv.Scanner.Scan() {
		line := wv.Scanner.Text()
		sumBytes += len(line) + 1

		// Skip empty lines
		switch {
		case checkHeader(line):
			if wv.Header.Head != "" && wv.Header.Note != "" {
				continue
			}
			if line == "WEBVTT" {
				wv.Header.Head = line
			} else {
				wv.Header.Note = line
			}
		case checkToken(StartOrEndTime, line):
			if isStartOrEndTime == 0 {
				isStartOrEndTime++
				e.StartTime = line
			} else {
				e.EndTime = line
				isStartOrEndTime--
			}

		case checkToken(Separator, line):
			e.Separator = line

		case checkToken(Position, line):
			e.Position = line

		case checkToken(Line, line):
			e.Line = line

		case line == "":
			wv.Elements = append(wv.Elements, e)
			e = &Element{}
		default:
			e.Text += line
		}
	}
	wv.sumBytes = sumBytes
	if err := wv.Scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

// read reads the WebVTT file.
func (wv *WebVtt) read() {
	wv.scanLines(scanSplitFunc)
	wv.unifyText()
	wv.deleteElementOfEmptyText()
}

// Read reads the WebVTT file from a byte slice.
func (wv *WebVtt) Read(p []byte) (n int, err error) {
	wv.Header = &Header{}
	wv.Scanner = bufio.NewScanner(bytes.NewReader(p))
	wv.builder = &strings.Builder{}
	wv.read()
	return int(wv.sumBytes), nil
}

// ReadFrom reads the WebVTT file from an io.Reader.
func (wv *WebVtt) ReadFrom(r io.Reader) (n int64, err error) {
	wv.Header = &Header{}
	wv.Scanner = bufio.NewScanner(r)
	wv.builder = &strings.Builder{}
	wv.read()
	return int64(wv.sumBytes), nil
}

// WriteTo writes the WebVTT file to an io.Writer.
func (wv *WebVtt) WriteTo(w io.Writer) (int64, error) {
	if wv.builder.Len() == 0 {
		wv.write()
	}
	n, err := w.Write([]byte(wv.builder.String()))
	return int64(n), err
}

// write writes the WebVTT file.
func (wv *WebVtt) write() {
	const (
		emptyRow = "\n"
		empty    = " "
	)

	// Header
	wv.builder.WriteString(wv.Header.Head + emptyRow)
	wv.builder.WriteString(wv.Header.Note + emptyRow)

	// Body
	for _, e := range wv.Elements {
		// Empty line
		wv.builder.WriteString(emptyRow)
		// Timeline part
		wv.builder.WriteString(e.StartTime + empty + e.Separator + empty +
			e.EndTime + empty + e.Position + empty + e.Line + emptyRow)
		wv.builder.WriteString(e.Text + emptyRow)
	}
}
