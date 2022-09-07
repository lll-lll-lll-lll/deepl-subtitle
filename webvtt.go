package deeplyoutubesubtitle

import (
	"bufio"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type WebVtt struct {
	VttFile     string         `json:"file"`
	VttElements []*VTTElement  `json:"vtt_elements"`
	VTTHeader   *VTTHeader     `json:"header"`
	VTTScanner  *bufio.Scanner `json:"scanner"`
}

func NewWebVtt(file string) *WebVtt {
	scanner := bufio.NewScanner(strings.NewReader(file))
	header := NewVTTHeader()
	return &WebVtt{VttFile: file, VTTScanner: scanner, VTTHeader: header}
}

type VTTHeader struct {
	Head string `json:"head"`
	Note string `json:"note"`
}

func NewVTTHeader() *VTTHeader {
	return &VTTHeader{}
}

type VTTElement struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Position  string `json:"position"`
	Line      string `json:"line"`
	Text      string `json:"text"`
	Separator string `json:"separator"`
}

//AppendVttElement append VTTElement to WebVtt
func (wv *WebVtt) AppendVttElement(vtt *VTTElement) {
	wv.VttElements = append(wv.VttElements, vtt)
}

// NewVttElement initialize VTTElement
func (wv *WebVtt) NewVttElement() *VTTElement {
	return &VTTElement{}
}

//ScanHeaderSplitFunc default split func
func ScanHeaderSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	return
}

func ScanSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	tokenStr := string(token)
	if CheckTimeRegexpFlag(tokenStr) || CheckSeparatorFlag(tokenStr) || CheckPositionFlag(tokenStr) || CheckLineFlag(tokenStr) {
		{
			advance, token, err = bufio.ScanWords(data, atEOF)
			return
		}
	}

	return
}

//ReadVTTFile use when WebVTT struct is initialized.
func ReadVTTFile(filename string) (string, error) {
	ext := filepath.Ext(filename)
	if ext != ".vtt" {
		return "", errors.New("your input file extension is not `.vtt`. check your file extension")
	}

	bytesFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", nil
	}
	if string(bytesFile) == "" {
		return "", errors.New("file content is empty")
	}
	return string(bytesFile), nil
}
