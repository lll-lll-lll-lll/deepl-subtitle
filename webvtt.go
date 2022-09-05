package deeplyoutubesubtitle

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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

//SkipHeader ignore header of vtt file.
func (wv *WebVtt) SkipHeader() {
	var lineNum = 0
	for wv.VTTScanner.Scan() {
		text := wv.VTTScanner.Text()
		switch lineNum {
		case 0:
			wv.VTTHeader.Head = text
		case 1:
			wv.VTTHeader.Note = text
		default:
			break
		}
		lineNum++
	}
	if err := wv.VTTScanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

//ScanLines Read the contents of the file and structure the data
func (wv *WebVtt) ScanLines(splitFunc bufio.SplitFunc) {
	vttElement := wv.NewVttElement()
	wv.VTTScanner.Split(splitFunc)
	var vttElementFlag int

	for wv.VTTScanner.Scan() {
		line := wv.VTTScanner.Text()
		switch {

		case CheckTimeRegexpFlag(line):
			if vttElementFlag == 0 {
				vttElementFlag++
				vttElement.StartTime = line
			} else {
				vttElement.EndTime = line
				vttElementFlag--
			}

		case CheckSeparatorFlag(line):
			vttElement.Separator = line

		case CheckPositionFlag(line):
			vttElement.Position = line

		case CheckLineFlag(line):
			vttElement.Line = line

		case line == "":
			wv.AppendVttElement(vttElement)
			vttElement = wv.NewVttElement()
		default:
			vttElement.Text += line
		}
	}

	if err := wv.VTTScanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	// Skip head element header
	wv.VttElements = wv.VttElements[1:]
}

//ScanHeaderSplitFunc default split func
func ScanHeaderSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	return
}

func ScanTimeLineSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
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
	return string(bytesFile), nil
}
