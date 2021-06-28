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

		var styler Styler
		if opts.Html {
			styler = &HTMLStyler{
				theme: th,
				name:  lang,
			}
		} else {
			styler = &GChalkStyler{
				theme: th,
				lvl:   gchalk.GetLevel(),
			}
		}

		fmt.Print(highlight(h, th, f, styler))
	}
}

type Styler interface {
	Pre() string
	Style(s, group string) string
	Post() string
}

func highlight(h *flare.Highlighter, th theme.Theme, f io.ReaderAt, st Styler) string {
	buf := &bytes.Buffer{}
	buf.WriteString(st.Pre())
	h.HighlightFunc(f, memo.NoneTable{}, func(text []byte, group string) {
		fmt.Fprint(buf, st.Style(string(text), group))
	}, nil)
	buf.WriteString(st.Post())
	return buf.String()
}
