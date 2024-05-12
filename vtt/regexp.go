package vtt

import (
	"regexp"
)

type TokenType int

const (
	StartOrEndTime TokenType = iota
	Separator
	Position
	Line
	Terminal
)

// patterns maps each TokenType to its corresponding regular expression pattern.
var patterns = map[TokenType]string{
	StartOrEndTime: `^[0-9]+`,           // Check if the first 2 characters are 0~9 of int
	Separator:      `^-->`,              // Check if the first character is `-->`
	Position:       `^position:[0-9]+%`, // Check if the first character is `position:...`
	Line:           `^line:[0-9]+%`,     // Check if the first character is `line:...`
	Terminal:       `.`,                 // Check if the en character is `.` or `?`
}

// checkToken checks the token based on the given type.
// Example: checkToken(StartOrEndTime, "10:00") will return true if "10:00" matches the StartOrEndTime pattern.
func checkToken(tokenType TokenType, token string) bool {
	pattern, ok := patterns[tokenType]
	if !ok {
		return false
	}
	return checkReg(pattern, token)
}

// checkReg is a method for pattern detection of regular expressions.
// Example: checkReg(`^[0-9]+`, "10:00") will return true as "10:00" matches the pattern.
func checkReg(pattern, token string) bool {
	return regexp.MustCompile(pattern).Match([]byte(token))
}

// searchTerminalToken checks if the token contains a "." or "?".
// Example: searchTerminalToken("Hello?") will return [5, 6] as "?" is found at index 5.
func searchTerminalToken(token string) []int {
	r, _ := regexp.Compile("[.?]")
	locs := r.FindStringIndex(token)
	return locs
}

func checkHeader(token string) bool {
	return checkReg(`WEBVTT`, token) || checkReg(`Kind: captions`, token)
}
