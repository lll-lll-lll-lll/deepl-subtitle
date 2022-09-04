package cmd

import (
	"log"
	"testing"
)

func TestTextSegment(t *testing.T) {
	token := "you need to know it? ♪ I know ♪"
	filename := "testvtt.en-ehkg1hFWq8A.vtt"
	t.Run("get 「.」and 「?」", func(t *testing.T) {
		got := CheckTerminalFlag(token)
		want := true
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
		t.Log(got)
	})
	t.Run("", func(t *testing.T) {
		wantpr := token[:19] // you need to know it
		wantba := token[20:] // ♪ I know ♪
		wantterminal := "?"
		pr, ba, ter := SplitByCommaAndQuestion(token)

		if pr != wantpr {
			t.Errorf("got %s want %s", pr, wantpr)
		}
		if ba != wantba {
			t.Errorf("got %s want %s", ba, wantba)
		}

		if ter != wantterminal {
			t.Errorf("got %s want %s", "?", wantterminal)
		}
		t.Log(pr, wantterminal, ba)
	})

	t.Run("Unify Text", func(t *testing.T) {
		f, err := ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := NewWebVtt(f)
		webVtt.ScanLines(ScanTimeLineSplitFunc)

	})
}
