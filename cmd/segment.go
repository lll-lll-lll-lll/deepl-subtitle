package cmd

import "regexp"

func SearchTerminalTokenRegexp(token string) []int {
	r, _ := regexp.Compile("[.?]")
	locs := r.FindStringIndex(token)
	return locs
}

func SplitByCommaAndQuestion(token string) (string, string, string) {
	locs := SearchTerminalTokenRegexp(token)
	terminal := token[locs[0]:locs[1]]
	previousString := token[:locs[0]]
	backString := token[locs[1]:]
	return previousString, backString, terminal
}
