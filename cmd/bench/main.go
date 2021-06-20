package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/jwalton/gchalk"
	"github.com/zyedidia/flare"
	"github.com/zyedidia/flare/theme"
	"github.com/zyedidia/gpeg/memo"
	"github.com/zyedidia/rope"
)

func fatal(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}

var letters = []byte("\n \tabcdefghijklmnopqrstuvwxyz")

func randbytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

func main() {
	flag.Parse()
	inputf := flag.Args()[0]

	data, err := ioutil.ReadFile(inputf)
	if err != nil {
		fatal(err)
	}

	r := rope.New(data)

	h, err := flare.LoadHighlighter("java", true)
	if err != nil {
		fatal(err)
	}

	tbl := memo.NewTreeTable(128)

	const nedits = 1000

	for i := 0; i < nedits; i++ {
		loc := rand.Intn(r.Len())
		text := randbytes(4)
		edit := memo.Edit{
			Start: loc,
			End:   loc,
			Len:   len(text),
		}

		r.Insert(loc, text)

		tbl.ApplyEdit(edit)

		start := time.Now()
		h.Highlight(r, tbl, nil, nil)
		fmt.Println(time.Since(start))
	}

	th, err := theme.LoadTheme("monokai")
	if err != nil {
		fatal(err)
	}

	buf := &bytes.Buffer{}
	h.Highlight(r, tbl, func(text []byte, group string) {
		// fmt.Println(strconv.Quote(string(text)), group)
		style := th.Style(group)
		gc := gchalk.New()
		if style.Fg != nil {
			gc = gc.WithRGB(style.Fg.R, style.Fg.G, style.Fg.B)
		}
		fmt.Fprintf(buf, gc.StyleMust()(string(text)))
	}, nil)
	fmt.Print(buf.String())
}
