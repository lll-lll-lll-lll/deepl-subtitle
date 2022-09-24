package main

import (
	"fmt"
	"log"

	ds "github.com/lll-lll-lll-lll/deepl-subtitle"
)

func main() {
	filename := "example.vtt"
	f, err := ds.ReadVTTFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	webVtt := ds.NewWebVtt(f)
	webVtt.ScanLines(ds.ScanSplitFunc)
	w := ds.UnifyTextByTerminalPoint(webVtt)
	a := ds.DeleteVTTElementOfEmptyText(w)
	// a.ToFile("testoutput")
	ds.PrintlnJson(a.VttElements)
	fmt.Println(a.VTTHeader)
}
