package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/zyedidia/flare"
	"github.com/zyedidia/ftdetect"
	"github.com/zyedidia/gpeg/memo"
)

var theme = map[string]*color.Color{
	"":           color.New(),
	"whitespace": color.New(),
	"class":      color.New(color.FgBlue),
	"keyword":    color.New(color.FgRed),
	"type":       color.New(color.FgBlue),
	"function":   color.New(color.FgGreen),
	"identifier": color.New(),
	"string":     color.New(color.FgYellow),
	"comment":    color.New(color.FgBlack),
	"number":     color.New(color.FgMagenta),
	"constant":   color.New(color.FgMagenta),
	"preproc":    color.New(color.FgRed),
	"annotation": color.New(color.FgCyan),
	"operator":   color.New(),
	"special":    color.New(color.FgMagenta),
	"other":      color.New(),
}

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

	h, err := flare.LoadHighlighter(*lang)
	if err != nil {
		fatal(err)
	}

	tbl := memo.NewTreeTable(128)

	buf := &bytes.Buffer{}
	h.Highlight(f, tbl, func(text []byte, group string) {
		clr := theme[group]
		clr.Fprint(buf, string(text))
	})
	fmt.Print(buf.String())
}
