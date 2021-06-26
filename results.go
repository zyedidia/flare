package flare

import (
	"sort"
)

type HighlightResults struct {
	tokens []token
	last   int
}

func (r *HighlightResults) add(idx int, group string) {
	r.tokens = append(r.tokens, token{
		idx:   idx,
		group: group,
	})
}

func (r *HighlightResults) Group(idx int) string {
	if r.last >= 0 && r.last < len(r.tokens) {
		if idx >= r.tokens[r.last].idx && r.last == len(r.tokens) {
			return r.tokens[r.last].group
		}
		if idx >= r.tokens[r.last].idx && idx < r.tokens[r.last+1].idx {
			return r.tokens[r.last].group
		}
	}

	i := sort.Search(len(r.tokens), func(i int) bool { return r.tokens[i].idx > idx })
	r.last = i - 1
	return r.tokens[i-1].group
}

type token struct {
	idx   int
	group string
}
