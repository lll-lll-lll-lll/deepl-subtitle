// Package webvtt parse file of vtt extension.
//
// This package doesn't guarantee parsing incorrectly formatted vtt file
// Correct formatted vtt file example following:
//
// 00:00:00.350 --> 00:00:01.530 position:63% line:0%
// - Yo what is going on guys,

// 00:00:01.530 --> 00:00:02.770 position:63% line:0%
// welcome back to the channel.
package webvtt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lll-lll-lll-lll/format-webvtt/sub"
)

var (
	ErrVTTEXT = errors.New("invalid file extension")
)

// To distinguish from string
type FileName string

type WebVtt struct {
	File string `json:"file"`
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

func New(file FileName) *WebVtt {
	f := string(file)
	scanner := bufio.NewScanner(strings.NewReader(f))
	header := NewHeader()
	return &WebVtt{File: f, Scanner: scanner, Header: header}
}

func NewHeader() *Header {
	return &Header{}
}

func (wv *WebVtt) NewElement() *Element {
	return &Element{}
}

// AppendElement append VTTElement to WebVtt
func (wv *WebVtt) AppendElement(vtt *Element) {
	wv.Elements = append(wv.Elements, vtt)
}

func ScanSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	t := string(token)
	// CheckTimeRegexpFlagでtrueが走るとその行を空白で単語区切りにする。トークン区切りになった他の`-->`や`position...`を他のフラグで検索
	if sub.CheckStartOrEndTime(t) || sub.CheckSeparator(t) || sub.CheckPosition(t) || sub.CheckLine(t) {
		{
			advance, token, err = bufio.ScanWords(data, atEOF)
			return
		}
	}
	return
}

// Read use when WebVTT struct is initialized.
func Read(filename string) (FileName, error) {
	ext := filepath.Ext(filename)
	if ext != ".vtt" {
		return "", fmt.Errorf("%w. input extension is %v", ErrVTTEXT, ext)
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if len(b) == 0 {
		return "", fmt.Errorf("file content is empty")
	}
	return FileName(b), nil
}

func PrintlnJson(elements []*Element) {
	for _, e := range elements {
		var out bytes.Buffer
		b, _ := json.Marshal(e)
		err := json.Indent(&out, b, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(out.String())
	}
}

// UnifyText Updates EndTime by passing Text to the previous structure if it contains `.` or `?`.
func (wv *WebVtt) UnifyText() {
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
func (wv *WebVtt) DeleteElementOfEmptyText() {
	var i int
	f := true
	for f {
		if wv.Elements[i].Text == "" {
			wv.Elements = append(wv.Elements[:i], wv.Elements[i+1:]...)
			i--
		}
		i++
		if len(wv.Elements) == i {
			f = false
		}
	}
}

// RecursiveSearchTerminalPoint SearchTerminalTokenRegexp メソッドで文末トークンが見つかるまでの構造体の個数を返す
func RecursiveSearchTerminalPoint(vs []*Element, untilTerminalCnt int) int {
	if untilTerminalCnt == len(vs)-1 {
		return untilTerminalCnt
	}
	e := vs[untilTerminalCnt].Text
	locs := sub.SearchTerminalToken(e)
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
	e := wv.NewElement()
	wv.Scanner.Split(splitFunc)
	var isStartOrEndTime int

	for wv.Scanner.Scan() {
		line := wv.Scanner.Text()
		switch {
		case sub.CheckHeader(line):
			if wv.Header.Head != "" && wv.Header.Note != "" {
				continue
			}
			if line == "WEBVTT" {
				wv.Header.Head = line
			} else {
				wv.Header.Note = line
			}
		case sub.CheckStartOrEndTime(line):
			if isStartOrEndTime == 0 {
				isStartOrEndTime++
				e.StartTime = line
			} else {
				e.EndTime = line
				isStartOrEndTime--
			}

		case sub.CheckSeparator(line):
			e.Separator = line

		case sub.CheckPosition(line):
			e.Position = line

		case sub.CheckLine(line):
			e.Line = line

		case line == "":
			wv.AppendElement(e)
			e = wv.NewElement()
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

func (wv *WebVtt) ToFile(onlyFileName string) {
	const (
		emptyRow = "\n"
		empty    = " "
	)

	f, err := os.Create(onlyFileName)
	if err != nil {
		log.Fatal(err)
	}

	// Header
	_, err = f.WriteString(wv.Header.Head + emptyRow)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(wv.Header.Note + emptyRow)
	if err != nil {
		log.Fatal(err)
	}

	// Body
	for _, e := range wv.Elements {
		// 空行
		_, err = f.WriteString(emptyRow)
		if err != nil {
			log.Fatal(err)
		}
		// timelineの部分
		_, err = f.WriteString(e.StartTime + empty + e.Separator + empty +
			e.EndTime + empty + e.Position + empty + e.Line + emptyRow)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.WriteString(e.Text + emptyRow)
	}
	if err != nil {
		log.Fatal(err)
	}
}
