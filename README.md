# deepl-subtitle

### Example

<details><summary>example.vtt</summary><div>

```
WEBVTT
Kind: captions

00:00:00.350 --> 00:00:01.530 position:63% line:0%
- Yo what is going on guys,

00:00:01.530 --> 00:00:02.770 position:63% line:0%
welcome back to the channel.

00:00:02.770 --> 00:00:05.240 position:63% line:0%
My name's Sonny and today
I'm gonna teach you all about

00:00:05.240 --> 00:00:06.730 position:63% line:0%
the useEffect Hook

00:00:06.730 --> 00:00:08.840 position:63% line:0%
and why it has transformed.

00:00:08.840 --> 00:00:11.110 position:63% line:0%
the way that we use
functional components and why

00:00:11.110 --> 00:00:12.158 position:63% line:0%
you need to know it.
♪ I know ♪
```
</div></details>

```go
package main

import (
	ds "github.com/lll-lll-lll-lll/deepl-subtitle"
	"log"
)

func main() {
	filename := "example.vtt"
	f, err := ds.ReadVTTFile(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	webVtt := ds.NewWebVtt(f)
	webVtt.ScanLines(ds.ScanTimeLineSplitFunc)
	w := ds.UnifyTextByTerminalPoint(webVtt)
	a := ds.DeleteVTTElementOfEmptyText(w)
	// console
	ds.PrintlnJson(a.VttElements)
}

```

<details><summary>output json</summary><div>

```json

{
  "start_time": "00:00:00.350",
  "end_time": "00:00:02.770",
  "position": "position:63%",
  "line": "line:0%",
  "text": "- Yo what is going on guys, welcome back to the channel.",
  "separator": "--\u003e"
},
{
  "start_time": "00:00:02.770",
  "end_time": "00:00:08.840",
  "position": "position:63%",
  "line": "line:0%",
  "text": "My name's Sonny and todayI'm gonna teach you all about the useEffect Hook and why it has transformed.",
  "separator": "--\u003e"
},
{
  "start_time": "00:00:08.840",
  "end_time": "00:00:12.158",
  "position": "position:63%",
  "line": "line:0%",
  "text": "the way that we usefunctional components and why you need to know it.♪ I know ♪",
  "separator": "--\u003e"
}]
```
</div></details>