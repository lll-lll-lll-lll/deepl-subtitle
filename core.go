package deeplyoutubesubtitle

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	EXTVTT = ".vtt"
)

//UnifyTextByTerminalPoint `.` か `?`を含んでいたら１つ前の構造体にTextを渡し、EndTimeを更新するメソッド
func UnifyTextByTerminalPoint(webVtt *WebVtt) *WebVtt {
	es := webVtt.VttElements
	for i := 0; i < len(es)-1; i++ {
		// どこまでのテキストを繋げてよいかを表す値を取得
		cnt := RecursiveSearchTerminalPoint(es, i)
		for j := cnt; j > i; j-- {
			ct := es[j].Text
			cet := es[j].EndTime
			es[j-1].Text += " " + ct
			es[j-1].EndTime = cet
			es[j].Text = ""
		}
		// 文末を表現するトークンを見つけた位置まで移動
		if cnt > 0 {
			i = cnt
		}
	}
	webVtt.VttElements = es
	return webVtt
}

//DeleteVTTElementOfEmptyText テキストが空の構造体を削除するメソッド
func DeleteVTTElementOfEmptyText(webVtt *WebVtt) *WebVtt {
	var i int
	f := true
	es := webVtt.VttElements
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
	webVtt.VttElements = es
	return webVtt
}

//RecursiveSearchTerminalPoint SearchTerminalTokenRegexp メソッドで文末トークンが見つかるまでの構造体の個数を返す
func RecursiveSearchTerminalPoint(vs []*VTTElement, untilTerminalCnt int) int {
	if untilTerminalCnt == len(vs)-1 {
		return untilTerminalCnt
	}
	e := vs[untilTerminalCnt].Text
	locs := SearchTerminalTokenRegexp(e)
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

//ScanLines 一行ずつ読み込んで構造体を作成するメソッド
func (wv *WebVtt) ScanLines(splitFunc bufio.SplitFunc) {
	vttElement := wv.NewVttElement()
	wv.VTTScanner.Split(splitFunc)
	var vttElementFlag int

	for wv.VTTScanner.Scan() {
		line := wv.VTTScanner.Text()
		switch {
		case CheckHeaderFlag(line):
			if wv.VTTHeader.Head != "" && wv.VTTHeader.Note != "" {
				continue
			}
			if line == "WEBVTT" {
				wv.VTTHeader.Head = line
			} else {
				wv.VTTHeader.Note = line
			}
		case CheckStartOrEndTimeFlag(line):
			if vttElementFlag == 0 {
				vttElementFlag++
				vttElement.StartTime = line
			} else {
				vttElement.EndTime = line
				vttElementFlag--
			}

		case CheckSeparatorFlag(line):
			vttElement.Separator = line

		case CheckPositionFlag(line):
			vttElement.Position = line

		case CheckLineFlag(line):
			vttElement.Line = line

		case line == "":
			wv.AppendVttElement(vttElement)
			vttElement = wv.NewVttElement()
		default:
			vttElement.Text += line
		}
	}

	if err := wv.VTTScanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	// Skip head element header
	wv.VttElements = wv.VttElements[1:]
}

//ToFile 文字列をファイルに戻すメソッド.
func (wv *WebVtt) ToFile(onlyFileName string) {
	const (
		emptyRow = "\n"
		empty    = " "
	)

	f, err := os.Create(onlyFileName + EXTVTT)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Header
	_, err = f.WriteString(wv.VTTHeader.Head + emptyRow)
	check(err)
	_, err = f.WriteString(wv.VTTHeader.Note + emptyRow)
	check(err)

	// Body
	for _, e := range wv.VttElements {
		// 空行
		_, err = f.WriteString(emptyRow)
		// timelineの部分
		_, err = f.WriteString(e.StartTime + empty + e.Separator + empty +
			e.EndTime + empty + e.Position + empty + e.Line + emptyRow)
		_, err = f.WriteString(e.Text + emptyRow)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
