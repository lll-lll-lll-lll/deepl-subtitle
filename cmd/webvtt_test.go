package cmd

import (
	"log"
	"os/exec"
	"testing"
)

func TestGetVtt(t *testing.T) {
	url := "https://www.youtube.com/watch?v=YS4e4q9oBaU&t=3764s"
	filename := "testvtt"
	cmd := exec.Command(
		"yt-dlp", "--skip-download", "--sub-format",
		"vtt", "--write-subs", "--sub-langs", "en", "-o", filename, url)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func TestWebVTTStruct(t *testing.T) {
	t.Run("webvtt starttime", func(t *testing.T) {
		webVtt := WebVtt{}
		vttElement := &VTTElement{
			StartTime: "00:00:06.649",
		}
		webVtt.AppendVttElement(vttElement)

		got := webVtt.VttElements[0]
		want := VTTElement{StartTime: "00:00:06.649"}
		if got.StartTime != want.StartTime {
			t.Errorf("got %s want %s", got.StartTime, want.StartTime)
		}
	})
	t.Run("append webvtt element to WebVtt", func(t *testing.T) {
		webVtt := WebVtt{}
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
