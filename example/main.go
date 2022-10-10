package main

import (
	"fmt"
	"log"
	"time"

	ds "github.com/lll-lll-lll-lll/deepl-subtitle"
)

func main() {
	filename := "example.vtt"
	fmt.Println("start reading file.")
	f, err := ds.ReadVTTFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	webVtt := ds.NewWebVtt(f)
	start := time.Now()
	fmt.Println("start scanning file")
	webVtt.ScanLines(ds.ScanSplitFunc)
	fmt.Println("start unify text by terminal point")
	w := ds.UnifyTextByTerminalPoint(webVtt)
	fmt.Println("start delete empty text of vtt element")
	ds.DeleteVTTElementOfEmptyText(w)
	// a.ToFile("testoutput")
	ds.PrintlnJson(w.VttElements)
	fmt.Println("start calculate untile end")
	fmt.Println(time.Since(start).Seconds())
}
