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
	f, err := ds.ReadVTT(filename)
	if err != nil {
		log.Fatal(err)
	}
	webVtt := ds.New(f)
	start := time.Now()
	fmt.Println("start scanning file")
	webVtt.ScanLines(ds.ScanSplitFunc)
	fmt.Println("start unify text by terminal point")
	w := ds.UnifyText(webVtt)
	fmt.Println("start delete empty text of vtt element")
	ds.DeleteElementOfEmptyText(w)
	// a.ToFile("testoutput")
	ds.PrintlnJson(w.Elements)
	fmt.Println("start calculate untile end")
	fmt.Println(time.Since(start).Seconds())
}
