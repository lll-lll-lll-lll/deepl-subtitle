package cmd

import (
	"log"
	"os/exec"
	"testing"
)

func TestGetVtt(t *testing.T) {
	url := "https://www.youtube.com/watch?v=YS4e4q9oBaU&t=3764s"
	filename := "testvtt"
	cmd := exec.Command("yt-dlp", "--skip-download", "--sub-format", "vtt", "--write-subs", "--sub-langs", "en", "-o", filename, url)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func TestWebVTTStruct(t *testing.T) {
	webVtt := WebVtt{}
	t.Run("webVtt startTime", func(t *testing.T) {
		vttElement := &VTTElement{
			StartTime: "00:00:06.649",
			EndTime:   "00:00:10.690",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "So I wanted to start out with an introduction\nto the go language itself. Now I know that",
		}
		webVtt.AppendVttElement(vttElement)

		got := webVtt.VttElements[0]
		want := VTTElement{
			StartTime: "00:00:06.649",
			EndTime:   "00:00:10.690",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "So I wanted to start out with an introduction\nto the go language itself. Now I know that",
		}
		if got.StartTime != want.StartTime {
			t.Errorf("got %s want %s", got.StartTime, want.StartTime)
		}
		if got.EndTime != want.EndTime {
			t.Errorf("got %s want %s", got.EndTime, want.EndTime)
		}
		if got.Line != want.Line {
			t.Errorf("got %s want %s", got.Line, want.Line)
		}
		if got.Position != want.Position {
			t.Errorf("got %s want %s", got.Position, want.Position)
		}
		if got.Text != want.Text {
			t.Errorf("got %s want %s", got.Text, want.Text)
		}

	})

	t.Run("append webVtt element to WebVtt", func(t *testing.T) {
		vttElement1 := &VTTElement{
			StartTime: "00:00:06.649",
		}
		vttElement2 := &VTTElement{
			StartTime: "00:00:15.620",
		}
		webVtt.AppendVttElement(vttElement1)
		webVtt.AppendVttElement(vttElement2)
		for _, vtt := range webVtt.VttElements {
			log.Println(vtt)
		}
	})
}
