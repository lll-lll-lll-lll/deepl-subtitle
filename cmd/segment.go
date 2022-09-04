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

func UnifyTextByTerminalPoint(webVtt *WebVtt) *WebVtt {
	es := webVtt.VttElements

	for i := 0; i < len(es)-1; i++ {
		curt := es[i].Text
		locs := SearchTerminalTokenRegexp(curt)

		// no `?` or `.` in text
		if haveTerminalPoint(locs) {
			continue
		}

		// after terminal point text
		pt := es[i].Text[:locs[0]]
		bt := es[i].Text[locs[1]:]
		next := es[i+1].Text
		locs = SearchTerminalTokenRegexp(next)
		if haveTerminalPoint(locs) {
			es[i+1].Text = bt + next
			es[i].Text = ""
		}
		// previous terminal point text
		es[i-1].Text += pt
	}
	webVtt.VttElements = es
	return webVtt
}
