
# Specification
https://www.w3.org/TR/webvtt1/

This is a tool for formatting VTT (Video Text Tracks) files. It follows the specification provided by W3C, which can be found [here](https://www.w3.org/TR/webvtt1/).

## Installation

To install `vtt-formatter`, you need to tap into our Homebrew repository and then install the package. Run the following commands in your terminal:

```sh
brew tap lll-lll-lll-lll/vtt-formatter
brew install vtt-formatter
```

## Usage
You can use VTT Formatter from the command line with the following options:

Options:

- `-h` or `--help`: Display help information
- `-v` or `--version`: Display the current version of VTT Formatter
- `-i` or `--input=<filename.vtt>`: Specify the VTT file to be formatted
- `-o` or `--output=<filename.vtt>`: Specify the destination file for the formatted VTT
- `-p` or `--print`: Print the formatted VTT to the console

## CLI Example

```sh
vtt-formatter --input example.vtt --print --output formattedfile.vtt


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
