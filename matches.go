package flare

import (
	"sort"
)

// Matches is a structure that stores highlighting results corresponding to a
// certain range of the document.  It supports querying the highlight group of
// any index within the range.
type Matches struct {
	tokens []token
	last   int

	start, end int
}

func (m *Matches) add(idx int, group string) {
	m.tokens = append(m.tokens, token{
		idx:   idx,
		group: group,
	})
}

// InRange returns true if the given offset is in this Matches structure.
func (m *Matches) InRange(p int) bool {
	return p >= m.start && p < m.end
}

// Group returns the highlighting group at the specified index.
func (m *Matches) Group(idx int) string {
	if !m.InRange(idx) {
		return ""
	}

	if m.last >= 0 && m.last < len(m.tokens) {
		if idx >= m.tokens[m.last].idx && m.last >= len(m.tokens)-1 {
			return m.tokens[m.last].group
		}
		if idx >= m.tokens[m.last].idx && idx < m.tokens[m.last+1].idx {
			return m.tokens[m.last].group
		}
	}

	i := sort.Search(len(m.tokens), func(i int) bool { return m.tokens[i].idx > idx })
	m.last = i - 1
	return m.tokens[i-1].group
}

type token struct {
	idx   int
	group string
}
