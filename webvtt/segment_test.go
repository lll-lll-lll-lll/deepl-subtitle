package webvtt

import (
	"log"
	"testing"

	"github.com/lll-lll-lll-lll/deepl-subtitle/sub"
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
		f, err := ReadVTT(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := New(f)
		webVtt.ScanLines(ScanSplitFunc)
		w := UnifyText(webVtt)
		got := w.Elements[0].Text
		want := "- Yo what is going on guys, welcome back to the channel."
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
		//PrintlnJson(w.VttElements)
	})
	t.Run("`DeleteEmptyTextVTTElementStruct` method test", func(t *testing.T) {
		f, err := ReadVTT(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := New(f)
		webVtt.ScanLines(ScanSplitFunc)
		w := UnifyText(webVtt)
		DeleteElementOfEmptyText(w)
		got := len(w.Elements)
		want := 3
		if got != want {
			t.Errorf("got %d, want %d ", got, want)
		}
		PrintlnJson(w.Elements)
	})
}
