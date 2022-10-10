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
	VttFile     string         `json:"file"`
	VttElements []*VTTElement  `json:"vtt_elements"`
	VTTHeader   *VTTHeader     `json:"header"`
	VTTScanner  *bufio.Scanner `json:"scanner"`
}

type VTTElement struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Position  string `json:"position"`
	Line      string `json:"line"`
	Text      string `json:"text"`
	Separator string `json:"separator"`
}

type VTTHeader struct {
	Head string `json:"head"`
	Note string `json:"note"`
}

func NewWebVtt(file WebVttString) *WebVtt {
	f := string(file)
	scanner := bufio.NewScanner(strings.NewReader(f))
	header := NewVTTHeader()
	return &WebVtt{VttFile: f, VTTScanner: scanner, VTTHeader: header}
}

func NewVTTHeader() *VTTHeader {
	return &VTTHeader{}
}

func (wv *WebVtt) NewVttElement() *VTTElement {
	return &VTTElement{}
}

// AppendVttElement append VTTElement to WebVtt
func (wv *WebVtt) AppendVttElement(vtt *VTTElement) {
	wv.VttElements = append(wv.VttElements, vtt)
}

func ScanSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	tokenStr := string(token)
	// CheckTimeRegexpFlagでtrueが走るとその行を空白で単語区切りにする。トークン区切りになった他の`-->`や`position...`を他のフラグで検索
	if CheckStartOrEndTimeFlag(tokenStr) || CheckSeparatorFlag(tokenStr) || CheckPositionFlag(tokenStr) || CheckLineFlag(tokenStr) {
		{
			advance, token, err = bufio.ScanWords(data, atEOF)
			return
		}
	}

	return
}

// ReadVTTFile use when WebVTT struct is initialized.
func ReadVTTFile(filename string) (WebVttString, error) {
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
