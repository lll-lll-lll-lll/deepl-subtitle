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
	webVtt := webvtt.New(f)
	webVtt.ScanLines(webvtt.ScanSplitFunc)
	webVtt.UnifyText()
	webVtt.DeleteElementOfEmptyText()
	// a.ToFile("testoutput")
	webvtt.PrintlnJson(webVtt.Elements)
}
