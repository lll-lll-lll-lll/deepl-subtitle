package main

import (
	ds "github.com/lll-lll-lll-lll/deepl-subtitle"
	"log"
)

func main() {
	filename := ""
	f, err := ds.ReadVTTFile(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	webVtt := ds.NewWebVtt(f)
	webVtt.ScanLines(ds.ScanTimeLineSplitFunc)
	w := ds.UnifyTextByTerminalPoint(webVtt)
	a := ds.DeleteVTTElementOfEmptyText(w)
	ds.PrintlnJson(a.VttElements)
}
