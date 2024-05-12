package vtt

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type WebVtt struct {
	fileContent string `json:"file"`
	// Collection of vtt structure
	//
	// Example:
	// 00:00:00.350 --> 00:00:01.530 position:63% line:0%
	// Yo what is going on guys,
	//
	// 00:00:01.530 --> 00:00:02.770 position:63% line:0%
	// welcome back to the channel.
	//
	Elements []*Element     `json:"vtt_elements"`
	Header   *Header        `json:"header"`
	Scanner  *bufio.Scanner `json:"scanner"`
}

// A pattern from the vtt file dropped into the structure
//
// Example:
// 00:00:00.350(StartTime) -->(Separator) 00:00:01.530(EndTime) position:63%(Position) line:0%(Line)
// Yo what is going on guys,(Text)
type Element struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Position  string `json:"position"`
	Line      string `json:"line"`
	Text      string `json:"text"`
	Separator string `json:"separator"`
}

// First two lines of the vtt file
//
// Example:
// WEBVTT
// Kind: captions
type Header struct {
	Head string `json:"head"`
	Note string `json:"note"`
}

func New(file string) *WebVtt {
	scanner := bufio.NewScanner(strings.NewReader(file))
	header := &Header{}
	return &WebVtt{fileContent: file, Scanner: scanner, Header: header}
}

func (wv *WebVtt) Format() {
	wv.ScanLines(scanSplitFunc)
	wv.unifyText()
	wv.deleteElementOfEmptyText()
}

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

// ReadFileContents is used when the WebVTT struct is initialized.
func ReadFileContents(filename string) (string, error) {
	ext := filepath.Ext(filename)
	if ext != ".vtt" {
		return "", fmt.Errorf("invalid file extension: %v, expected .vtt", ext)
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v, error: %w", filename, err)
	}
	if len(b) == 0 {
		return "", fmt.Errorf("file content is empty: %v", filename)
	}
	return string(b), nil
}

// unifyText Updates EndTime by passing Text to the previous structure if it contains `.` or `?`.
func (wv *WebVtt) unifyText() {
	for i := 0; i < len(wv.Elements)-1; i++ {
		untilTerminalPointCnt := RecursiveSearchTerminalPoint(wv.Elements, i)
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

// DeleteElementOfEmptyText
// Loop until all structures with empty text are deleted
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

// RecursiveSearchTerminalPoint is a method that recursively searches for
// the terminal point in the given slice of VTT elements.
// It returns the index of the element that contains the terminal point.
func RecursiveSearchTerminalPoint(vs []*Element, untilTerminalCnt int) int {
	if untilTerminalCnt == len(vs)-1 {
		return untilTerminalCnt
	}
	e := vs[untilTerminalCnt].Text
	locs := searchTerminalToken(e)
	f := func(locs []int) bool {
		return len(locs) == 0
	}
	if f(locs) {
		untilTerminalCnt++
		return RecursiveSearchTerminalPoint(vs, untilTerminalCnt)
	}
	return untilTerminalCnt
}

// ScanLines 一行ずつ読み込んで構造体を作成するメソッド
func (wv *WebVtt) ScanLines(splitFunc bufio.SplitFunc) {
	e := &Element{}
	wv.Scanner.Split(splitFunc)
	var isStartOrEndTime int

	for wv.Scanner.Scan() {
		line := wv.Scanner.Text()
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

	if err := wv.Scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	// Skip head element header
	wv.Elements = wv.Elements[1:]
}

func (wv *WebVtt) WriteToFile(fileName string) error {
	const (
		emptyRow = "\n"
		empty    = " "
	)

	var builder strings.Builder

	// Header
	builder.WriteString(wv.Header.Head + emptyRow)
	builder.WriteString(wv.Header.Note + emptyRow)

	// Body
	for _, e := range wv.Elements {
		// Empty line
		builder.WriteString(emptyRow)
		// Timeline part
		builder.WriteString(e.StartTime + empty + e.Separator + empty +
			e.EndTime + empty + e.Position + empty + e.Line + emptyRow)
		builder.WriteString(e.Text + emptyRow)
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(builder.String())
	if err != nil {
		return err
	}

	return nil
}

func Print(elements []*Element) error {
	d, err := json.Marshal(elements)
	if err != nil {
		return err
	}
	fmt.Println(string(d))
	return nil
}
