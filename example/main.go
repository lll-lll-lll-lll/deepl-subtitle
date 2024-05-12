package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lll-lll-lll-lll/vtt-formatter/vtt"
)

func main() {
	filename := "../data/go_learn.vtt"
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("ファイルを開くのに失敗しました:", err)
		return
	}
	defer f.Close()

	wvtt := vtt.New(f)
	wvtt.Format()
	fff, err := os.Create("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fff.Close()
	wvtt.WriteTo(os.Stdout)
}
