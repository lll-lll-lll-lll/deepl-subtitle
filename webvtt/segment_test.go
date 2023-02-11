package webvtt

import (
	"log"
	"testing"

	"github.com/lll-lll-lll-lll/sn-formater/sub"
)

func TestTextSegment(t *testing.T) {
	token := "you need to know it, ♪ I know ♪"
	filename := "../data/example.vtt"
	t.Run("get 「.」and 「?」", func(t *testing.T) {
		got := sub.CheckTerminal(token)
		want := true
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
		t.Log(got)
	})

	t.Run("`UnifyTextByTerminalPoint` method test ", func(t *testing.T) {
		f, err := Read(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := New(f)
		webVtt.ScanLines(ScanSplitFunc)
		webVtt.UnifyText()
		got := webVtt.Elements[0].Text
		want := "- Yo what is going on guys, welcome back to the channel."
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
	t.Run("`DeleteEmptyTextVTTElementStruct` method test", func(t *testing.T) {
		f, err := Read(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := New(f)
		webVtt.ScanLines(ScanSplitFunc)
		webVtt.UnifyText()
		webVtt.DeleteElementOfEmptyText()
		got := len(webVtt.Elements)
		want := 3
		if got != want {
			t.Errorf("got %d, want %d ", got, want)
		}
	})
}
