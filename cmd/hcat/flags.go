package main

var opts struct {
	// original cat flags
	All         bool `short:"A" long:"show-all" description:"equivalent to -vET"`
	NumNonBlank bool `short:"b" long:"number-nonblank" description:"number nonempty output lines, overrides -n"`
	Ends        bool `short:"E" long:"show-ends" description:"display $ at end of each line"`
	Num         bool `short:"n" long:"number" description:"number all output lines"`
	Squeeze     bool `short:"s" long:"squeeze-blank" description:"suppress repeated empty output lines"`
	Tabs        bool `short:"T" long:"show-tabs" description:"display TAB characters as ^I"`
	_           bool `short:"u" description:"(ignored)"`
	NonPrint    bool `short:"v" long:"show-nonprinting" description:"use ^ and M- notation, except for LFD and TAB"`
	Help        bool `long:"help" description:"display this help and exit"`
	Version     bool `long:"version" description:"output version information and exit"`

	ColorLvl int    `short:"c" long:"color" description:"color level: none (0), 16-color (1), 256-color (2), true-color (3)" default:"-1" default-mask:"-"`
	Lang     string `short:"l" long:"lang" description:"language to use for highlighting (autodetect if omitted)"`
	Theme    string `long:"theme" description:"color theme to use" default:"monokai"`
}
