package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	EXTVTT = ".vtt"
)

// UnifyText `.` か `?`を含んでいたら１つ前の構造体にTextを渡し、EndTimeを更新するメソッド
func UnifyText(webVtt *WebVtt) *WebVtt {
	ves := webVtt.Elements
	for i := 0; i < len(ves)-1; i++ {
		// どこまでのテキストを繋げてよいかを表す値を取得
		untilTerminalPointCnt := RecursiveSearchTerminalPoint(ves, i)
		for j := untilTerminalPointCnt; j > i; j-- {
			t := ves[j].Text
			e := ves[j].EndTime
			ves[j-1].Text += " " + t
			ves[j-1].EndTime = e
			ves[j].Text = ""
		}
		// 文末を表現するトークンを見つけた位置まで移動
		if untilTerminalPointCnt > 0 {
			i = untilTerminalPointCnt
		}
	}
	webVtt.Elements = ves
	return webVtt
}

// DeleteElementOfEmptyText テキストが空の構造体を削除するメソッド
func DeleteElementOfEmptyText(webVtt *WebVtt) {
	var i int
	f := true
	es := webVtt.Elements
	// 空のテキストを持つ構造体を削除し切るまでループ
	for f {
		// ここで空のテキストを持つ構造体を削除
		if es[i].Text == "" {
			es = append(es[:i], es[i+1:]...)
			i--
		}
		i++
		if len(es) == i {
			f = false
		}
	}
	webVtt.Elements = es
}

// RecursiveSearchTerminalPoint SearchTerminalTokenRegexp メソッドで文末トークンが見つかるまでの構造体の個数を返す
func RecursiveSearchTerminalPoint(vs []*Element, untilTerminalCnt int) int {
	if untilTerminalCnt == len(vs)-1 {
		return untilTerminalCnt
	}
	e := vs[untilTerminalCnt].Text
	locs := SearchTerminalToken(e)
	f := func(locs []int) bool {
		if len(locs) == 0 {
			return true
		}
		return false
	}
	if f(locs) == true {
		untilTerminalCnt++
		return RecursiveSearchTerminalPoint(vs, untilTerminalCnt)
	}
	return untilTerminalCnt
}

// ScanLines 一行ずつ読み込んで構造体を作成するメソッド
func (wv *WebVtt) ScanLines(splitFunc bufio.SplitFunc) {
	vttElement := wv.NewElement()
	wv.Scanner.Split(splitFunc)
	var vttElementFlag int

	for wv.Scanner.Scan() {
		line := wv.Scanner.Text()
		switch {
		case CheckHeader(line):
			if wv.Header.Head != "" && wv.Header.Note != "" {
				continue
			}
			if line == "WEBVTT" {
				wv.Header.Head = line
			} else {
				wv.Header.Note = line
			}
		case CheckStartOrEndTime(line):
			if vttElementFlag == 0 {
				vttElementFlag++
				vttElement.StartTime = line
			} else {
				vttElement.EndTime = line
				vttElementFlag--
			}

		case CheckSeparator(line):
			vttElement.Separator = line

		case CheckPosition(line):
			vttElement.Position = line

		case CheckLine(line):
			vttElement.Line = line

		case line == "":
			wv.AppendElement(vttElement)
			vttElement = wv.NewElement()
		default:
			vttElement.Text += line
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
	check(err)
	_, err = f.WriteString(wv.Header.Note + emptyRow)
	check(err)

	// Body
	for _, e := range wv.Elements {
		// 空行
		_, err = f.WriteString(emptyRow)
		// timelineの部分
		_, err = f.WriteString(e.StartTime + empty + e.Separator + empty +
			e.EndTime + empty + e.Position + empty + e.Line + emptyRow)
		_, err = f.WriteString(e.Text + emptyRow)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
