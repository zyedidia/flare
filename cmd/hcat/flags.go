package main

var opts struct {
	ColorLvl int    `short:"c" long:"color" description:"color level: none (0), 16-color (1), 256-color (2), true-color (3)" default:"-1" default-mask:"-"`
	Lang     string `short:"l" long:"lang" description:"language to use for highlighting (autodetect if omitted)"`
	Theme    string `long:"theme" description:"color theme to use" default:"monokai"`
	Help     bool   `short:"h" long:"help" description:"display this help and exit"`
	Version  bool   `long:"version" description:"output version information and exit"`
}
