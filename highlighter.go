package flare

import (
	"io"

	"github.com/zyedidia/gpeg/input"
	"github.com/zyedidia/gpeg/memo"
	"github.com/zyedidia/gpeg/pattern"
	p "github.com/zyedidia/gpeg/pattern"
	"github.com/zyedidia/gpeg/vm"
)

type Highlighter struct {
	code     vm.Code
	captures map[int]string
}

var empty Highlighter = Highlighter{
	code: vm.Encode(pattern.MustCompile(p.Star(p.Any(1)))),
}

func (h *Highlighter) Highlight(r io.ReaderAt, tbl memo.Table, draw func(text []byte, group string), i *vm.Interval) {
	match, _, ast, _ := h.code.ExecInterval(r, tbl, i)
	if !match {
		return
	}

	if draw != nil {
		colorize(0, ast, h.captures, r, draw)
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
