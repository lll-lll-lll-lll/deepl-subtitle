# unify-vtt-text

# Specification
https://www.w3.org/TR/webvtt1/


## Install
```sh
brew tap lll-lll-lll-lll/vtt-formatter
```

```sh
brew install vtt-formatter
```


```md
Usage: vreader [options] 

Options:
  -help or h 	 		        help
  -version            		 now version
  -filepath=<{filename}.vtt>    vtt file name
  -path=<{filename}.vtt>    File name of destination
  -p                        print json in console
```

### CLI Example

```sh
 vtt-formatter  -filepath example.vtt -pj -path shapedfile.vtt
```



```
WEBVTT
Kind: captions

00:00:00.350 --> 00:00:02.770 position:63% line:0%
- Yo what is going on guys, welcome back to the channel.

00:00:02.770 --> 00:00:08.840 position:63% line:0%
My name's Sonny and todayI'm gonna teach you all about the useEffect Hook and why it has transformed.

00:00:08.840 --> 00:00:12.158 position:63% line:0%
the way that we usefunctional components and why you need to know it.♪ I know ♪

```
</div></details>
<br>
<br>

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
}
```
</div></details>