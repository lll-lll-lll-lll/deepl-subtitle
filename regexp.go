package main

import (
	"regexp"
)

// CheckStartOrEndTime Check if the first 2 characters are 0~9 of int
func CheckStartOrEndTime(token string) bool {
	return CheckReg(`^[0-9]+`, token)
}

// CheckSeparator // Check if the first character is `-->`
func CheckSeparator(token string) bool {
	return CheckReg(`^-->`, token)
}

// CheckPosition Check if the first character is `position:...`
func CheckPosition(token string) bool {
	return CheckReg(`^position:[0-9]+%`, token)
}

// CheckLine Check if the first character is `line:...`
func CheckLine(token string) bool {
	return CheckReg(`^line:[0-9]+%`, token)
}

// CheckTerminal Check if the en character is `.` or `?`
func CheckTerminal(token string) bool {
	return CheckReg(`.`, token) || CheckReg(`?`, token)
}

// CheckReg Pattern detection of regular expression things method
func CheckReg(pattern, str string) bool {
	return regexp.MustCompile(pattern).Match([]byte(str))
}

// SearchTerminalToken 「.」か「?」を含んでいるか
func SearchTerminalToken(token string) []int {
	r, _ := regexp.Compile("[.?]")
	locs := r.FindStringIndex(token)
	return locs
}

// CheckHeader headerならtrue
func CheckHeader(token string) bool {
	return CheckReg(`WEBVTT`, token) || CheckReg(`Kind: captions`, token)
}
