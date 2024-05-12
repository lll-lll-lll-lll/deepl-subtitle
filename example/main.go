package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lll-lll-lll-lll/vtt-formatter/vtt"
)

func main() {
	filename := "../data/example.vtt"
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("ファイルを開くのに失敗しました:", err)
		return
	}
	defer f.Close()
	wvtt := &vtt.WebVtt{}
	n, err := wvtt.ReadFrom(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("test", n)
	fff, err := os.Create("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fff.Close()
	if _, err := wvtt.WriteTo(fff); err != nil {
		log.Fatal(err)
	}
}
