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

// SearchTerminalTokenRegexpで返ってきたlocs内に「.」もしくは「?」の位置を示す値が入っているか
func haveTerminalPoint(locs []int) bool {
	if len(locs) == 0 {
		return true
	}
	return false
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
