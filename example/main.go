package main

import (
	"fmt"
	ds "github.com/lll-lll-lll-lll/deepl-subtitle"
	"log"
)

func main() {
	filename := "example.vtt"
	f, err := ds.ReadVTTFile(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	webVtt := ds.NewWebVtt(f)
	webVtt.ScanLines(ds.ScanSplitFunc)
	w := ds.UnifyTextByTerminalPoint(webVtt)
	a := ds.DeleteVTTElementOfEmptyText(w)
	ds.PrintlnJson(a.VttElements)
	fmt.Println(a.VTTHeader)
}
