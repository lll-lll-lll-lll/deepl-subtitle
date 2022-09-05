package deeplyoutubesubtitle

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// func SplitByCommaAndQuestion(token string) (string, string, string, string, bool) {
// 	locs := SearchTerminalTokenRegexp(token)
// 	if len(locs) == 0 {
// 		return token, "", "", "", true
// 	}
// 	terminal := token[locs[0]:locs[1]]
// 	previousString := token[:locs[0]]
// 	backString := token[locs[1]:]
// 	return token, previousString, backString, terminal, false
// }

func haveTerminalPoint(locs []int) bool {
	if len(locs) == 0 {
		return true
	}
	return false
}

func RecursiveSearchTerminalPoint(vs []*VTTElement, untilTerminalCnt int) int {
	e := vs[untilTerminalCnt].Text
	locs := SearchTerminalTokenRegexp(e)
	if haveTerminalPoint(locs) == true {
		untilTerminalCnt++
		return RecursiveSearchTerminalPoint(vs, untilTerminalCnt)
	}
	return untilTerminalCnt
}

func UnifyTextByTerminalPoint(webVtt *WebVtt) *WebVtt {
	es := webVtt.VttElements
	for i := 0; i < len(es)-1; i++ {
		untilTerminalPointCnt := RecursiveSearchTerminalPoint(es, i)
		for j := untilTerminalPointCnt; j > i; j-- {
			currentToken := es[j].Text
			currentEndTime := es[j].EndTime
			es[j-1].Text += " " + currentToken
			es[j-1].EndTime = currentEndTime
			es[j].Text = ""
		}
		if untilTerminalPointCnt != 0 {
			i = untilTerminalPointCnt
		}
	}
	webVtt.VttElements = es
	return webVtt
}

func DeleteVTTElementStructOfEmptyText(webVtt *WebVtt) *WebVtt {
	var i int
	f := true
	es := webVtt.VttElements

	for f {
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

func PrintlnJson(elements []*VTTElement) {
	for _, e := range elements {
		var out bytes.Buffer
		b, _ := json.Marshal(e)
		err := json.Indent(&out, b, "", "  ")
		if err != nil {
			panic(err)

		}
		fmt.Println(out.String())
	}
}
