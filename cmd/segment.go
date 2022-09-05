package cmd

import "regexp"

func SearchTerminalTokenRegexp(token string) []int {
	r, _ := regexp.Compile("[.?]")
	locs := r.FindStringIndex(token)
	return locs
}

func SplitByCommaAndQuestion(token string) (string, string, string, string, bool) {
	locs := SearchTerminalTokenRegexp(token)
	if len(locs) == 0 {
		return token, "", "", "", true
	}
	terminal := token[locs[0]:locs[1]]
	previousString := token[:locs[0]]
	backString := token[locs[1]:]
	return token, previousString, backString, terminal, false
}

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
		untilTerminalCnt := RecursiveSearchTerminalPoint(es, i)
		for j := untilTerminalCnt; j > i; j-- {
			curt := es[j].Text
			es[j-1].Text += " " + curt
			es[j].Text = ""
		}
		if untilTerminalCnt != 0 {
			i = untilTerminalCnt
		}
	}
	webVtt.VttElements = es
	return webVtt
}
