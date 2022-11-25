package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"
)

//func TestGetVtt(t *testing.T) {
//	//url := "https://www.youtube.com/watch?v=YS4e4q9oBaU&t=3764s"
//	shortMovie := "https://www.youtube.com/watch?v=UVhIMwHDS7k"
//	filename := "testvtt"
//	cmd := exec.Command("yt-dlp", "--skip-download", "-o", filename, "--sub-format", "vtt", "--write-subs", shortMovie)
//	err := cmd.Run()
//	if err != nil {
//		log.Fatalln(err)
//		return
//	}
//}

// TestWebVTTStruct Test method and property of WebVtt struct
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
			Text:      "and why it has transformed.",
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
	filename := "example.vtt"
	t.Run("webVtt startTime", func(t *testing.T) {
		webVtt := WebVtt{}
		vttElement := &Element{
			StartTime: "00:00:06.649",
			EndTime:   "00:00:10.690",
			Position:  "position:63%",
			Line:      "line:0%",
			Text:      "So I wanted to start out with an introduction to the go language itself. Now I know that",
		}
		webVtt.AppendElement(vttElement)

		got := webVtt.Elements[0]
		want := Element{
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
		vttElement1 := &Element{
			StartTime: "00:00:06.649",
		}
		vttElement2 := &Element{
			StartTime: "00:00:15.620",
		}
		webVtt.AppendElement(vttElement1)
		webVtt.AppendElement(vttElement2)

		if len(webVtt.Elements) != 2 {
			t.Errorf("didn't append VttElement. \n got %d want 2", len(webVtt.Elements))
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

		//bytesFile, err := ioutil.ReadVTTFile(filename)
		//if err != nil {
		//	t.Errorf("err is %s", err)
		//}

		err = os.Remove("copy" + filename)
		if err != nil {
			t.Errorf("remove err %s", err)
		}
	})

	t.Run("ScanTimeLine method test. create VTT Element struct", func(t *testing.T) {
		f, err := ReadVTT(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := New(f)
		webVtt.ScanLines(ScanSplitFunc)
		got := webVtt.Elements[0]
		want := &Element{
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
		f, err := ReadVTT(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := New(f)
		webVtt.ScanLines(ScanSplitFunc)
		elements := webVtt.Elements
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
		}
	})
	t.Run("Println Json VTTElement", func(t *testing.T) {
		f, err := ReadVTT(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := New(f)
		webVtt.ScanLines(ScanSplitFunc)
		elements := webVtt.Elements
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
	t.Run("perfect all processed", func(t *testing.T) {
		allProcessedDone := []struct {
			StartTime string
			EndTime   string
			Position  string
			Line      string
			Text      string
		}{
			{
				StartTime: "00:00:00.350",
				EndTime:   "00:00:02.770",
				Position:  "position:63%",
				Line:      "line:0%",
				Text:      "- Yo what is going on guys, welcome back to the channel.",
			},
			{
				StartTime: "00:00:02.770",
				EndTime:   "00:00:08.840",
				Position:  "position:63%",
				Line:      "line:0%",
				Text:      "",
			},
			{
				StartTime: "00:00:08.840",
				EndTime:   "00:00:12.158",
				Position:  "position:63%",
				Line:      "line:0%",
				Text:      "",
			},
		}
		f, err := ReadVTT(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := New(f)
		webVtt.ScanLines(ScanSplitFunc)
		w := UnifyText(webVtt)
		DeleteVTTElementOfEmptyText(w)
		e := w.Elements
		for i, tt := range allProcessedDone {
			d := e[i]
			if tt.StartTime != d.StartTime {
				t.Errorf("got %s want %s", tt.StartTime, d.StartTime)
			}
			if tt.EndTime != d.EndTime {
				t.Errorf("got %s want %s", tt.EndTime, d.EndTime)
			}
		}
		PrintlnJson(e)
	})

	t.Run("`ToFile method test`", func(t *testing.T) {
		f, err := ReadVTT(filename)
		if err != nil {
			log.Fatal(err)
		}
		webVtt := New(f)
		webVtt.ScanLines(ScanSplitFunc)
		w := UnifyText(webVtt)
		DeleteVTTElementOfEmptyText(w)
		w.ToFile("test")
	})

}
