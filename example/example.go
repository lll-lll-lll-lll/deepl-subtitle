package main

import (
	"log"

	"github.com/lll-lll-lll-lll/vtt-formatter/vtt"
)

func main_() {
	filename := "data/example.vtt"
	f, err := vtt.ReadFileContents(filename)
	if err != nil {
		log.Fatal(err)
	}
	vtt := vtt.New(f)
	vtt.Format()
}
