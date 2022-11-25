package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type WebVttString string

type WebVtt struct {
	File     string         `json:"file"`
	Elements []*Element     `json:"vtt_elements"`
	Header   *Header        `json:"header"`
	Scanner  *bufio.Scanner `json:"scanner"`
}

type Element struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Position  string `json:"position"`
	Line      string `json:"line"`
	Text      string `json:"text"`
	Separator string `json:"separator"`
}

type Header struct {
	Head string `json:"head"`
	Note string `json:"note"`
}

func New(file WebVttString) *WebVtt {
	f := string(file)
	scanner := bufio.NewScanner(strings.NewReader(f))
	header := NewHeader()
	return &WebVtt{File: f, Scanner: scanner, Header: header}
}

func NewHeader() *Header {
	return &Header{}
}

func (wv *WebVtt) NewElement() *Element {
	return &Element{}
}

// AppendElement append VTTElement to WebVtt
func (wv *WebVtt) AppendElement(vtt *Element) {
	wv.Elements = append(wv.Elements, vtt)
}

func ScanSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	tokenStr := string(token)
	// CheckTimeRegexpFlagでtrueが走るとその行を空白で単語区切りにする。トークン区切りになった他の`-->`や`position...`を他のフラグで検索
	if CheckStartOrEndTime(tokenStr) || CheckSeparator(tokenStr) || CheckPosition(tokenStr) || CheckLine(tokenStr) {
		{
			advance, token, err = bufio.ScanWords(data, atEOF)
			return
		}
	}

	return
}

// ReadVTT use when WebVTT struct is initialized.
func ReadVTT(filename string) (WebVttString, error) {
	ext := filepath.Ext(filename)
	if ext != ".vtt" {
		return "", errors.New("your input file extension is not `.vtt`. check your file extension")
	}

	bytesFile, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if string(bytesFile) == "" {
		return "", errors.New("file content is empty")
	}
	return WebVttString(bytesFile), nil
}
