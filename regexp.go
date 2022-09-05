package deeplyoutubesubtitle

import (
	"io/ioutil"
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

//ReadFile use when WebVTT struct is initialized.
func ReadFile(filename string) (string, error) {
	bytesFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", nil
	}
	return string(bytesFile), nil
}

//CheckRegexp Pattern detection of regular expression things method
func CheckRegexp(pattern, str string) bool {
	return regexp.MustCompile(pattern).Match([]byte(str))
}
