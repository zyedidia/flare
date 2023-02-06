package main

var opts struct {
	Html  bool   `long:"html" description:"output HTML"`
	Lang  string `short:"l" long:"lang" description:"language to use for highlighting (autodetect if omitted)"`
	Theme string `long:"theme" description:"color theme to use" default:"monokai"`
	// List    bool   `long:"list" description:"list available themes and languages"`
	Help    bool `short:"h" long:"help" description:"display this help and exit"`
	Version bool `long:"version" description:"output version information and exit"`
}
