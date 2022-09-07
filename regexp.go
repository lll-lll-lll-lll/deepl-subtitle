package deeplyoutubesubtitle

import (
	"regexp"
)

//CheckTimeRegexpFlag Check if the first 2 characters are 0~9 of int
func CheckTimeRegexpFlag(token string) bool {
	return CheckRegexp(`^[0-9]+`, token)
}

//CheckSeparatorFlag // Check if the first character is `-->`
func CheckSeparatorFlag(token string) bool {
	return CheckRegexp(`^-->`, token)
}

//CheckPositionFlag Check if the first character is `position:...`
func CheckPositionFlag(token string) bool {
	return CheckRegexp(`^position:[0-9]+%`, token)
}

//CheckLineFlag Check if the first character is `line:...`
func CheckLineFlag(token string) bool {
	return CheckRegexp(`^line:[0-9]+%`, token)
}

//CheckTerminalFlag Check if the en character is `.` or `?`
func CheckTerminalFlag(token string) bool {
	return CheckRegexp(`.`, token) || CheckRegexp(`?`, token)
}

//CheckRegexp Pattern detection of regular expression things method
func CheckRegexp(pattern, str string) bool {
	return regexp.MustCompile(pattern).Match([]byte(str))
}

//SearchTerminalTokenRegexp 「.」か「?」を含んでいるか
func SearchTerminalTokenRegexp(token string) []int {
	r, _ := regexp.Compile("[.?]")
	locs := r.FindStringIndex(token)
	return locs
}

//CheckHeaderFlag headerならtrue
func CheckHeaderFlag(token string) bool {
	return CheckRegexp(`WEBVTT`, token) || CheckRegexp(`Kind: captions`, token)
}
