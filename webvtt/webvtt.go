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

	"github.com/lll-lll-lll-lll/deepl-subtitle/sub"
)

const (
	EXTVTT = ".vtt"
)

type WebVttString string

type WebVtt struct {
	File     string         `json:"file"`
	Elements []*Element     `json:"vtt_elements"`
	Header   *Header        `json:"header"`
	Scanner  *bufio.Scanner `json:"scanner"`
}

type Element struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Position  string `json:"position"`
	Line      string `json:"line"`
	Text      string `json:"text"`
	Separator string `json:"separator"`
}

type Header struct {
	Head string `json:"head"`
	Note string `json:"note"`
}

func New(file WebVttString) *WebVtt {
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
func Read(filename string) (WebVttString, error) {
	ext := filepath.Ext(filename)
	if ext != ".vtt" {
		return "", errors.New("your input file extension is not `.vtt`. check your file extension")
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if string(b) == "" {
		return "", errors.New("file content is empty")
	}
	return WebVttString(b), nil
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

// UnifyText `.` か `?`を含んでいたら１つ前の構造体にTextを渡し、EndTimeを更新するメソッド
func (wv *WebVtt) UnifyText() {
	for i := 0; i < len(wv.Elements)-1; i++ {
		// どこまでのテキストを繋げてよいかを表す値を取得
		untilTerminalPointCnt := RecursiveSearchTerminalPoint(wv.Elements, i)
		for j := untilTerminalPointCnt; j > i; j-- {
			t := wv.Elements[j].Text
			e := wv.Elements[j].EndTime
			wv.Elements[j-1].Text += " " + t
			wv.Elements[j-1].EndTime = e
			wv.Elements[j].Text = ""
		}
		// 文末を表現するトークンを見つけた位置まで移動
		if untilTerminalPointCnt > 0 {
			i = untilTerminalPointCnt
		}
	}
}

// DeleteElementOfEmptyText テキストが空の構造体を削除するメソッド
func (wv *WebVtt) DeleteElementOfEmptyText() {
	var i int
	f := true
	// 空のテキストを持つ構造体を削除し切るまでループ
	for f {
		// ここで空のテキストを持つ構造体を削除
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

// ToFile 文字列をファイルに戻すメソッド.
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
