package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/zyedidia/flare"
	"github.com/zyedidia/gpeg/input"
	"github.com/zyedidia/gpeg/memo"
	"github.com/zyedidia/gpeg/pattern"
	"github.com/zyedidia/gpeg/vm"
)

var theme = map[string]*color.Color{
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

func colorize(start int, token *memo.Capture, captures map[int]string, theme map[string]*color.Color, in io.ReaderAt) string {
	buf := &bytes.Buffer{}

	off := 0
	var clr *color.Color
	if token.Dummy() {
		clr = color.New()
	} else {
		clr = theme[captures[token.Id()]]
	}
	if clr == nil {
		clr = color.New()
	}

	it := token.ChildIterator(0)

	for ch := it(); ch != nil; ch = it() {
		chstart := ch.Start()
		buf.WriteString(clr.Sprint(string(input.Slice(in, start+off, chstart))))
		off = chstart - start
		s := colorize(chstart, ch, captures, theme, in)
		buf.WriteString(s)
		off += ch.Len()
	}
	buf.WriteString(clr.Sprint(string(input.Slice(in, start+off, start+token.Len()))))

	// for ch, choff := it(); ch != nil; ch, choff = it() {
	// 	buf.WriteString(clr.Sprint(string(input.Slice(in, start+off, start+choff))))
	// 	off = choff
	// 	s := colorize(start+choff, ch, theme, in)
	// 	buf.WriteString(s)
	// 	off += ch.Len()
	// }
	// buf.WriteString(clr.Sprint(string(input.Slice(in, start+off, start+token.Len()))))

	return buf.String()
}

func main() {
	lang := flag.String("lang", "", "language file")

	flag.Parse()

	if *lang == "" {
		fatal("No langauge provided")
	}

	ldata, err := ioutil.ReadFile(*lang)
	if err != nil {
		fatal(err)
	}

	l, err := flare.LoadLanguage(ldata)
	if err != nil {
		fatal(err)
	}
	h, err := l.Highlighter()
	if err != nil {
		fatal(err)
	}

	prog, err := pattern.Compile(h.Grammar)
	if err != nil {
		fatal(err)
	}
	code := vm.Encode(prog)

	data, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		fatal(err)
	}
	f := bytes.NewReader(data)

	start := time.Now()
	tbl := memo.NewTreeTable(128)
	match, _, ast, _ := code.Exec(f, tbl)
	fmt.Println("Parse time:", time.Since(start))

	if !match {
		fatal("error parsing document")
	}

	PrintMemUsage()

	// start = time.Now()
	// loc := rand.Intn(f.Len() - 1)
	// // loc := f.Len() - 1
	// data = insert(data, loc, []byte{'a'})
	// fmt.Println("EDIT AT", loc)
	// edit := memo.Edit{
	// 	Start: loc,
	// 	End:   loc,
	// 	Len:   1,
	// }
	// f = bytes.NewReader(data)
	// tbl.ApplyEdit(edit)
	// match, _, ast, _ = code.Exec(f, tbl)
	// fmt.Println("Reparse", match, time.Since(start))
	// fmt.Println(ast.NumChildren())

	fmt.Print(colorize(0, ast, h.Captures, theme, f))
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// from slice tricks
func insert(s []byte, k int, vs []byte) []byte {
	if n := len(s) + len(vs); n <= cap(s) {
		s2 := s[:n]
		copy(s2[k+len(vs):], s[k:])
		copy(s2[k:], vs)
		return s2
	}
	s2 := make([]byte, len(s)+len(vs))
	copy(s2, s[:k])
	copy(s2[k:], vs)
	copy(s2[k+len(vs):], s[k:])
	return s2
}
