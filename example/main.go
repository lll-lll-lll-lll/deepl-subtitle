package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lll-lll-lll-lll/deepl-subtitle/webvtt"
)

func main() {
	filename := "data/example.vtt"
	fmt.Println("start reading file.")
	f, err := webvtt.ReadVTT(filename)
	if err != nil {
		log.Fatal(err)
	}
	webVtt := webvtt.New(f)
	start := time.Now()
	fmt.Println("start scanning file")
	webVtt.ScanLines(webvtt.ScanSplitFunc)
	fmt.Println("start unify text by terminal point")
	w := webvtt.UnifyText(webVtt)
	fmt.Println("start delete empty text of vtt element")
	webvtt.DeleteElementOfEmptyText(w)
	// a.ToFile("testoutput")
	webvtt.PrintlnJson(w.Elements)
	fmt.Println("start calculate untile end")
	fmt.Println(time.Since(start).Seconds())
}
