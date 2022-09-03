package cmd

import (
	"bufio"
	"io"
	"os"
)

type WebVtt struct {
	VttFile     *os.File       `json:"file"`
	VttElements []*VTTElement  `json:"vtt_elements"`
	VTTHeader   *VTTHeader     `json:"header"`
	VTTScanner  *bufio.Scanner `json:"scanner"`
}

func NewWebVtt(f *os.File) *WebVtt {
	scanner := bufio.NewScanner(f)
	return &WebVtt{VttFile: f, VTTScanner: scanner}
}

type VTTHeader struct {
	Head string `json:"head"`
	Note string `json:"note"`
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
func (wv *WebVtt) SkipHeader() *VTTHeader {
	return &VTTHeader{}
}

//createFileObject use when WebVTT struct is initialized.
func createFileObject(filename string) (*os.File, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	copyFile, err := os.Create("copy" + filename)

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(copyFile, file)

	if err != nil {
		return nil, err
	}

	return copyFile, nil
}
