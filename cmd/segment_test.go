package cmd

import (
	"testing"
)

func TestTextSegment(t *testing.T) {
	t.Run("get 「.」and 「?」", func(t *testing.T) {
		token := "you need to know it? ♪ I know ♪"
		got := CheckTerminalFlag(token)
		want := true
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
		t.Log(got)
	})
	t.Run("", func(t *testing.T) {
		token := "you need to know it? ♪ I know ♪"
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
}
