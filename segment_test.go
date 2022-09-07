package deeplyoutubesubtitle

import (
	"log"
	"testing"
)

func TestTextSegment(t *testing.T) {
	token := "you need to know it, ♪ I know ♪"
	filename := "example.vtt"
	t.Run("get 「.」and 「?」", func(t *testing.T) {
		got := CheckTerminalFlag(token)
		want := true
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
		t.Log(got)
	})
	// t.Run("", func(t *testing.T) {
	// 	wantpr := token[:19] // you need to know it
	// 	wantba := token[20:] // ♪ I know ♪
	// 	wantterminal := "?"
	// 	locs := SearchTerminalTokenRegexp(token)
	// 	fmt.Print(locs)
	// 	_, pr, ba, ter, flag := SplitByCommaAndQuestion(token)
	// 	if flag {
	// 		t.Log("no word in text. `?` or `.`")
	// 	} else {
	// 		if pr != wantpr {
	// 			t.Errorf("got %s want %s", pr, wantpr)
	// 		}
	// 		if ba != wantba {
	// 			t.Errorf("got %s want %s", ba, wantba)
	// 		}

	// 		if ter != wantterminal {
	// 			t.Errorf("got %s want %s", "?", wantterminal)
	// 		}
	// 	}
	// })

	t.Run("`UnifyTextByTerminalPoint` method test ", func(t *testing.T) {
		f, err := ReadVTTFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := NewWebVtt(f)
		webVtt.ScanLines(ScanSplitFunc)
		w := UnifyTextByTerminalPoint(webVtt)
		got := w.VttElements[0].Text
		want := "- Yo what is going on guys, welcome back to the channel."
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
		//PrintlnJson(w.VttElements)
	})
	t.Run("`DeleteEmptyTextVTTElementStruct` method test", func(t *testing.T) {
		f, err := ReadVTTFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := NewWebVtt(f)
		webVtt.ScanLines(ScanSplitFunc)
		w := UnifyTextByTerminalPoint(webVtt)
		w = DeleteVTTElementOfEmptyText(w)
		got := len(w.VttElements)
		want := 3
		if got != want {
			t.Errorf("got %d, want %d ", got, want)
		}
		PrintlnJson(w.VttElements)
	})
}
