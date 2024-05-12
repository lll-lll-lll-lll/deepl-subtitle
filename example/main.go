package main

import (
	"log"
	"os"

	"github.com/lll-lll-lll-lll/vtt-formatter/vtt"
)

func main() {
	filename := "../data/go_learn.vtt"
	f, err := vtt.ReadFileContents(filename)
	if err != nil {
		log.Fatal(err)
	}
	wvtt := vtt.New(f)
	wvtt.Format()
	fff, err := os.Create("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fff.Close()
	wvtt.WriteTo(fff)

}
