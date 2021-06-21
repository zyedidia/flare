package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jwalton/gchalk"
	"github.com/zyedidia/flare"
	"github.com/zyedidia/flare/theme"
	"github.com/zyedidia/ftdetect"
	"github.com/zyedidia/gpeg/memo"
	"github.com/zyedidia/gpeg/vm"
	"github.com/zyedidia/rope"
)

func fatal(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}

var letters = []byte("\n \tabcdefghijklmnopqrstuvwxyz")

func randbytes(b []byte) {
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
}

var display = flag.Bool("display", false, "display the final highlighted text")
var nedits = flag.Int("n", 1000, "number of edits to perform")
var mthreshold = flag.Int("mthreshold", 128, "memoization entry size threshold")
var lang = flag.String("lang", "", "language to use for highlighting (autodetect if empty)")
var summary = flag.Bool("summary", false, "print performance summary")
var file = flag.Bool("file", false, "send output data to file")

func main() {
	flag.Parse()

	if len(flag.Args()) <= 0 {
		fmt.Println("Usage: bench [OPTIONS]... FILE")
		fmt.Println()
		flag.PrintDefaults()
		os.Exit(0)
	}

	inputf := flag.Args()[0]

	data, err := ioutil.ReadFile(inputf)
	if err != nil {
		fatal(err)
	}

	r := rope.New(data)

	hlang := *lang
	if hlang == "" {
		ds := ftdetect.LoadDefaultDetectors()
		first, _ := bufio.NewReader(bytes.NewReader(data)).ReadSlice('\n')
		first = bytes.TrimSpace(first)
		d := ds.Detect(inputf, first)

		if d != nil {
			hlang = d.Name
		}
	}

	h, err := flare.LoadHighlighter(hlang, true)
	if err != nil {
		fatal(err)
	}

	tbl := memo.NewTreeTable(*mthreshold)

	var membuf io.WriteCloser
	var timebuf io.WriteCloser

	if *file {
		membuf, err = os.Create("mem.dat")
		if err != nil {
			fatal(err)
		}
		timebuf, err = os.Create("time.dat")
		if err != nil {
			fatal(err)
		}
	} else {
		membuf = os.Stdout
		timebuf = os.Stdout
	}

	defer membuf.Close()
	defer timebuf.Close()

	var intrvl *vm.Interval = &vm.Interval{0, 4000}
	text := make([]byte, 4)

	h.Highlight(r, tbl, nil, intrvl)

	tottime := 0.0
	totmem := 0.0

	for i := 0; i < *nedits; i++ {
		loc := rand.Intn(r.Len())
		randbytes(text)
		edit := memo.Edit{
			Start: loc,
			End:   loc,
			Len:   len(text),
		}

		r.Insert(loc, text)

		start := time.Now()
		tbl.ApplyEdit(edit)
		h.Highlight(r, tbl, nil, intrvl)
		t := time.Since(start).Microseconds()
		tottime += float64(t)
		totmem += float64(memusage())
		fmt.Fprintln(timebuf, t)
		fmt.Fprintln(membuf, bToMb(memusage()))
	}

	if *display {
		th, err := theme.LoadTheme("monokai")
		if err != nil {
			fatal(err)
		}

		buf := &bytes.Buffer{}
		h.Highlight(r, tbl, func(text []byte, group string) {
			style := th.Style(group)
			gc := gchalk.New()
			if style.Fg != nil {
				gc = gc.WithRGB(style.Fg.R, style.Fg.G, style.Fg.B)
			}
			fmt.Fprintf(buf, gc.StyleMust()(string(text)))
		}, intrvl)
		fmt.Print(buf.String())
	}

	if *summary {
		fmt.Printf("Summary:\n")
		fmt.Printf("memo table size: %d\n", tbl.Size())
		fmt.Printf("final document size: %s\n", humanize.Bytes(uint64(r.Len())))
		fmt.Printf("avg time: %v\n", time.Microsecond*time.Duration(tottime/float64(*nedits)))
		fmt.Printf("avg mem: %s\n", humanize.Bytes(uint64(totmem/float64(*nedits))))
	}
}

func memusage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
