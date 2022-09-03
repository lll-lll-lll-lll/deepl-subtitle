package cmd

import (
	"io"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestGetVtt(t *testing.T) {
	url := "https://www.youtube.com/watch?v=YS4e4q9oBaU&t=3764s"
	//shortMovie := "https://www.youtube.com/watch?v=UVhIMwHDS7k"
	filename := "testvtt"
	cmd := exec.Command("yt-dlp", "--skip-download", "-o", filename, "--sub-format", "vtt", "--write-subs", url)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
		return
	}
}

//TestWebVTTStruct Test method and property of WebVtt struct
func TestWebVTTStruct(t *testing.T) {
	filename := "testvtt.en-ehkg1hFWq8A.vtt"
	t.Run("webVtt startTime", func(t *testing.T) {
		webVtt := WebVtt{}
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
		webVtt := WebVtt{}
		vttElement1 := &VTTElement{
			StartTime: "00:00:06.649",
		}
		vttElement2 := &VTTElement{
			StartTime: "00:00:15.620",
		}
		webVtt.AppendVttElement(vttElement1)
		webVtt.AppendVttElement(vttElement2)

		if len(webVtt.VttElements) != 2 {
			t.Errorf("didn't append VttElement. \n got %d want 2", len(webVtt.VttElements))
		}
	})

	t.Run("open vtt file", func(t *testing.T) {
		file, err := os.Open(filename)
		defer file.Close()

		copyFile, err := os.Create("copy" + filename)
		if err != nil {
			t.Errorf("%s", err)
		}

		_, err = io.Copy(copyFile, file)
		if err != nil {
			t.Errorf("file doesn't exist.")
		}

		//bytesFile, err := ioutil.ReadFile(filename)
		//if err != nil {
		//	t.Errorf("err is %s", err)
		//}

		err = os.Remove("copy" + filename)
		if err != nil {
			t.Errorf("remove err %s", err)
		}
	})

	t.Run("SkipHeader method test", func(t *testing.T) {
		f, err := CreateFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := NewWebVtt(f)
		webVtt.SkipHeader(ScanHeaderSplitFunc)
		want := VTTHeader{
			Head: "WEBVTT",
			Note: "Kind: captions",
		}

		if webVtt.VTTHeader.Head != want.Head {
			t.Errorf("got %s want %s", webVtt.VTTHeader.Head, want.Head)
		}
		if webVtt.VTTHeader.Note != want.Note {
			t.Errorf("got %s want %s", webVtt.VTTHeader.Note, want.Note)
		}
		t.Log(webVtt.VTTHeader.Head)
		t.Log(webVtt.VTTHeader.Note)
	})

	t.Run("ScanTimeLine method test. create VTT Element struct", func(t *testing.T) {
		f, err := CreateFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := NewWebVtt(f)
		webVtt.ScanTimeLine(ScanTimeLineSplitFunc)
		//got := webVtt.VttElements[0]
		//want := &VTTElement{
		//	StartTime: "00:00:06.649",
		//	EndTime:   "00:00:10.690",
		//	Position:  "position:63%",
		//	Line:      "line:0%",
		//	Text:      "So I wanted to start out with an introduction\nto the go language itself. Now I know that",
		//}
		//if got != want {
		//	t.Errorf("got %s want %s", got.Text, want.Text)
		//}
	})

	//t.Run("scan vtt file", func(t *testing.T) {
	//	webVtt := WebVtt{}
	//
	//	got := webVtt.Scanner()
	//	want := &VTTElement{
	//		StartTime: "00:00:06.649",
	//		EndTime:   "00:00:10.690",
	//		Position:  "position:63%",
	//		Line:      "line:0%",
	//		Text:      "So I wanted to start out with an introduction\nto the go language itself. Now I know that",
	//	}
	//
	//	if got != want {
	//		t.Errorf("didn't scan one block vtt")
	//	}
	//})
}
