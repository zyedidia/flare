package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jwalton/gchalk"
	"github.com/zyedidia/flare"
	"github.com/zyedidia/flare/theme"
	"github.com/zyedidia/ftdetect"
	"github.com/zyedidia/gpeg/memo"
)

func fatal(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}

func main() {
	lang := flag.String("lang", "", "language")

	flag.Parse()

	inputf := flag.Args()[0]
	data, err := ioutil.ReadFile(inputf)
	if err != nil {
		fatal(err)
	}
	f := bytes.NewReader(data)

	if *lang == "" {
		ds := ftdetect.LoadDefaultDetectors()
		first, _ := bufio.NewReader(f).ReadSlice('\n')
		first = bytes.TrimSpace(first)
		d := ds.Detect(inputf, first)

		if d != nil {
			*lang = d.Name
		}
	}

	fmt.Printf("Filetype: %s\n", *lang)

	h, err := flare.LoadHighlighter(*lang, false)
	if err != nil {
		fatal(err)
	}

	th, err := theme.LoadTheme("monokai")
	if err != nil {
		fatal(err)
	}

	buf := &bytes.Buffer{}
	h.Highlight(f, memo.NoneTable{}, func(text []byte, group string) {
		style := th.Style(group)
		fmt.Fprint(buf, stylize(string(text), style))
	})
	fmt.Print(buf.String())
}

func stylize(s string, style theme.Style) string {
	gc := gchalk.New()
	if style.Fg != nil {
		gc = gc.WithRGB(style.Fg.R, style.Fg.G, style.Fg.B)
	}
	return gc.StyleMust()(s)
}
