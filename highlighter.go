package flare

import (
	"io"

	"github.com/zyedidia/gpeg/input"
	"github.com/zyedidia/gpeg/memo"
	"github.com/zyedidia/gpeg/pattern"
	p "github.com/zyedidia/gpeg/pattern"
	"github.com/zyedidia/gpeg/vm"
)

// A Highlighter corresponds to a compiled GPeg program and a mapping from
// captures to highlight groups. It can be used to generate highlight results
// from a io.ReaderAt.
type Highlighter struct {
	code     vm.Code
	captures map[int]string
	backrefs map[int]string
}

var empty Highlighter = Highlighter{
	code: vm.Encode(pattern.MustCompile(p.Star(p.Any(1)))),
}

// HighlightMatches performs syntax highlighting on the given ReaderAt with a
// possibly pre-filled memoization table, over a specified interval. It
// generates a Matches structure that can be queried for the highlight results.
// If the specified interval is nil, the full file is used.
//
// NOTE: if using an interval, make sure it is as small as possible. Using nil
// and the full file range are not equivalent. Parse time is proportional to
// the size of the window, unless nil is used. However, if nil is used memory
// usage may be high if the input is large because the full parse tree is
// returned.
func (h *Highlighter) HighlightMatches(r io.ReaderAt, tbl memo.Table, i *vm.Interval) *Matches {
	match, _, ast, _ := h.code.ExecInterval(r, tbl, i)
	if !match {
		return nil
	}

	results := &Matches{}
	results.start, results.end = ast.Start(), ast.Start()+ast.Len()
	colorizeMatches(0, ast, h.captures, r, results)
	return results
}

// HighlightFunc is similar to HighlightMatches but calls a function for each
// match instead of storing it in a table.
func (h *Highlighter) HighlightFunc(r io.ReaderAt, tbl memo.Table, draw func(text []byte, group string), i *vm.Interval) {
	match, _, ast, _ := h.code.ExecInterval(r, tbl, i)
	if !match {
		return
	}

	if draw != nil {
		colorize(0, ast, h.captures, r, draw)
	}
}

func colorizeMatches(start int, token *memo.Capture, captures map[int]string, in io.ReaderAt, results *Matches) {
	off := 0
	var group string

	if !token.Dummy() {
		group = captures[token.Id()]
	}

	it := token.ChildIterator(0)

	for ch := it(); ch != nil; ch = it() {
		chstart := ch.Start()

		if start+off != chstart {
			results.add(start+off, group)
		}

		off = chstart - start
		colorizeMatches(chstart, ch, captures, in, results)
		off += ch.Len()
	}
	if start+off != start+token.Len() {
		results.add(start+off, group)
	}
}

func colorize(start int, token *memo.Capture, captures map[int]string, in io.ReaderAt, draw func(text []byte, group string)) {
	off := 0
	var group string

	if !token.Dummy() {
		group = captures[token.Id()]
	}

	it := token.ChildIterator(0)

	for ch := it(); ch != nil; ch = it() {
		chstart := ch.Start()

		if start+off != chstart {
			draw(input.Slice(in, start+off, chstart), group)
		}

		off = chstart - start
		colorize(chstart, ch, captures, in, draw)
		off += ch.Len()
	}
	if start+off != start+token.Len() {
		draw(input.Slice(in, start+off, start+token.Len()), group)
	}
}
