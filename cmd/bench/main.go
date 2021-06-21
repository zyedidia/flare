package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"time"

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

func randbytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

var display = flag.Bool("display", false, "display the final highlighted text")
var nedits = flag.Int("n", 1000, "number of edits to perform")
var mthreshold = flag.Int("mthreshold", 128, "memoization entry size threshold")
var lang = flag.String("lang", "", "language to use for highlighting (autodetect if empty)")
var mem = flag.Bool("mem", false, "print memory usage")

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

	var intrvl *vm.Interval
	for i := 0; i < *nedits; i++ {
		loc := rand.Intn(r.Len())
		text := randbytes(4)
		edit := memo.Edit{
			Start: loc,
			End:   loc,
			Len:   len(text),
		}

		r.Insert(loc, text)

		start := time.Now()
		tbl.ApplyEdit(edit)
		h.Highlight(r, tbl, nil, intrvl)
		fmt.Println(time.Since(start).Microseconds())
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

	if *mem {
		PrintMemUsage()
		fmt.Printf("memo table size: %d\n", tbl.Size())
	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
