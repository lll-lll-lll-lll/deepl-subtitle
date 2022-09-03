package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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
}

//Scan scan and bind one block of vtt file.
func (wv *WebVtt) Scan() *VTTElement {
	return nil
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
func (wv *WebVtt) SkipHeader(file string) {
	var lineNum = 0
	scanner := bufio.NewScanner(strings.NewReader(file))
	scanner.Split(ScanHeader)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println("text")
		log.Fatal(text)
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
}

func ScanHeader(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanWords(data, atEOF)
	return
}

func ScanTimeLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanWords(data, atEOF)

	if CheckRegexp(`^[0-9]+`, string(token)) {
		token = []byte("number match")
	}

	return
}

//CreateFile use when WebVTT struct is initialized.
func CreateFile(filename string) (string, error) {
	file, err := os.Open(filename)

	if err != nil {
		return "", err
	}
	defer file.Close()

	copyFile, err := os.Create("copy" + filename)

	if err != nil {
		return "", err
	}

	_, err = io.Copy(copyFile, file)

	if err != nil {
		return "", err
	}
	bytesFile, err := ioutil.ReadAll(file)
	return string(bytesFile), nil
}

//CheckRegexp Pattern detection of regular expression things method
func CheckRegexp(pattern, str string) bool {
	return regexp.MustCompile(pattern).Match([]byte(str))
}
