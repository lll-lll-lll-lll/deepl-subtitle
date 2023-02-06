package main

import (
	"fmt"
	"log"

	"github.com/lll-lll-lll-lll/deepl-subtitle/webvtt"
)

func main() {
	filename := "data/example.vtt"
	fmt.Println("start reading file.")
	f, err := webvtt.Read(filename)
	if err != nil {
		log.Fatal(err)
	}
	vtt := webvtt.New(f)
	vtt.ScanLines(webvtt.ScanSplitFunc)
	vtt.UnifyText()
	vtt.DeleteElementOfEmptyText()
	// a.ToFile("testoutput")
	webvtt.PrintlnJson(vtt.Elements)
}
