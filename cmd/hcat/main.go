package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/jwalton/gchalk"
	"github.com/zyedidia/flare"
	"github.com/zyedidia/flare/theme"
	"github.com/zyedidia/ftdetect"
	"github.com/zyedidia/gpeg/memo"
)

var Version = "v0.0.0-unknown"

func fatal(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}

func main() {
	flagparser := flags.NewParser(&opts, flags.PassDoubleDash|flags.PrintErrors)
	flagparser.Usage = "[OPTION]... [FILE]..."
	args, err := flagparser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println("hcat version", Version)
	}

	if opts.Help {
		flagparser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	th, err := theme.LoadTheme(opts.Theme)
	if err != nil {
		fatal(err)
	}

	lvl := gchalk.GetLevel()
	if opts.ColorLvl != -1 {
		lvl = gchalk.ColorLevel(opts.ColorLvl)
	}

	if len(args) == 0 {
		args = []string{"-"}
	}

	for _, inputf := range args {
		var data []byte
		var err error

		if inputf == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = ioutil.ReadFile(inputf)
		}

		if err != nil {
			fatal(err)
		}
		f := bytes.NewReader(data)

		lang := opts.Lang

		if opts.Lang == "" {
			ds := ftdetect.LoadDefaultDetectors()
			first, _ := bufio.NewReader(f).ReadSlice('\n')
			first = bytes.TrimSpace(first)
			d := ds.Detect(inputf, first)

			if d != nil {
				lang = d.Name
			}
		}

		h, err := flare.LoadHighlighter(lang, false)

		if err != nil {
			fatal(err)
		}

		buf := &bytes.Buffer{}
		h.Highlight(f, memo.NoneTable{}, func(text []byte, group string) {
			style := th.Style(group)
			fmt.Fprint(buf, stylize(string(text), style, lvl))
		})
		fmt.Print(buf.String())
	}
}

func stylize(s string, style theme.Style, lvl gchalk.ColorLevel) string {
	gc := gchalk.New(gchalk.ForceLevel(lvl))
	if style.Fg != nil {
		gc = gc.WithRGB(style.Fg.R, style.Fg.G, style.Fg.B)
	}
	return gc.StyleMust()(s)
}
