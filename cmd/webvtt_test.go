package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	webvttstructtest := []struct {
		StartTime string
		EndTime   string
		Position  string
		Line      string
		Text      string
	}{

		{
			StartTime: "00:00:00.350",
			EndTime:   "00:00:01.530",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "- Yo what is going on guys,",
		},
		{
			StartTime: "00:00:01.530",
			EndTime:   "00:00:02.770",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "welcome back to the channel.",
		},
		{
			StartTime: "00:00:02.770",
			EndTime:   "00:00:05.240",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "My name's Sonny and todayI'm gonna teach you all about",
		},
		{
			StartTime: "00:00:05.240",
			EndTime:   "00:00:06.730",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "the useEffect Hook",
		},
		{
			StartTime: "00:00:06.730",
			EndTime:   "00:00:08.840",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "and why it has transformed",
		},
		{
			StartTime: "00:00:08.840",
			EndTime:   "00:00:11.110",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "the way that we usefunctional components and why",
		},
		{
			StartTime: "00:00:11.110",
			EndTime:   "00:00:12.158",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "you need to know it.♪ I know ♪",
		},
	}
	filename := "testvtt.en-ehkg1hFWq8A.vtt"
	t.Run("webVtt startTime", func(t *testing.T) {
		webVtt := WebVtt{}
		vttElement := &VTTElement{
			StartTime: "00:00:06.649",
			EndTime:   "00:00:10.690",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "So I wanted to start out with an introduction to the go language itself. Now I know that",
		}
		webVtt.AppendVttElement(vttElement)

		got := webVtt.VttElements[0]
		want := VTTElement{
			StartTime: "00:00:06.649",
			EndTime:   "00:00:10.690",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "So I wanted to start out with an introduction to the go language itself. Now I know that",
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
		f, err := ReadFile(filename)
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
		f, err := ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := NewWebVtt(f)
		webVtt.ScanLines(ScanTimeLineSplitFunc)
		got := webVtt.VttElements[1]
		want := &VTTElement{
			StartTime: "00:00:06.649",
			EndTime:   "00:00:10.690",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "- Yo what is going on guys,",
		}
		if got.Text != want.Text {
			t.Errorf("got %s want %s", got.Text, want.Text)
		}
	})

	t.Run("test", func(t *testing.T) {
		f, err := ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := NewWebVtt(f)
		webVtt.ScanLines(ScanTimeLineSplitFunc)
		elements := webVtt.VttElements
		for i, tt := range webvttstructtest {
			d := elements[i]
			if tt.StartTime != d.StartTime {
				t.Errorf("got %s want %s", d.StartTime, tt.StartTime)
			}
			if tt.EndTime != d.EndTime {
				t.Errorf("got %s want %s", d.EndTime, tt.EndTime)
			}
			if tt.Position != d.Position {
				t.Errorf("got %s want %s", d.Position, tt.Position)
			}
			if tt.Line != d.Line {
				t.Errorf("got %s want %s", d.Line, tt.Line)
			}
			if tt.Text != d.Text {
				t.Errorf("got %s want %s", d.Text, tt.Text)
			}
			fmt.Println(d, tt)

		}
	})
	t.Run("Println Json VTTElement", func(t *testing.T) {
		f, err := ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := NewWebVtt(f)
		webVtt.ScanLines(ScanTimeLineSplitFunc)
		elements := webVtt.VttElements
		for _, e := range elements {
			var out bytes.Buffer
			b, _ := json.Marshal(e)
			err = json.Indent(&out, b, "", "  ")
			if err != nil {
				panic(err)
			}
			t.Log(out.String())
		}
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
