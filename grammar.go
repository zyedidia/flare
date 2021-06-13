package flare

import (
	"github.com/zyedidia/gpeg/charset"
	"github.com/zyedidia/gpeg/isa"
	p "github.com/zyedidia/gpeg/pattern"
)

var (
	alpha  = p.Set(charset.Range('A', 'Z').Add(charset.Range('a', 'z')))
	alnum  = p.Set(charset.Range('A', 'Z').Add(charset.Range('a', 'z')).Add(charset.Range('0', '9')))
	digit  = p.Set(charset.Range('0', '9'))
	xdigit = p.Set(charset.Range('0', '9').Add(charset.Range('A', 'F')).Add(charset.Range('a', 'f')))
	space  = p.Set(charset.New([]byte{9, 10, 11, 12, 13, ' '}))

	dec_num = p.Plus(digit)
	hex_num = p.Concat(
		p.Literal("0"),
		p.Set(charset.New([]byte{'x', 'X'})),
		p.Plus(xdigit),
	)
	oct_num = p.Concat(
		p.Literal("0"),
		p.Plus(p.Set(charset.New([]byte{'0', '7'}))),
	)

	integer = p.Concat(
		p.Optional(p.Set(charset.New([]byte{'+', '-'}))),
		p.Or(
			hex_num,
			oct_num,
			dec_num,
		),
	)
	float = p.Concat(
		p.Optional(p.Set(charset.New([]byte{'+', '-'}))),
		p.Or(
			p.Concat(
				p.Concat(
					p.Or(
						p.Concat(
							p.Star(digit),
							p.Literal("."),
							p.Plus(digit),
						),
						p.Concat(
							p.Plus(digit),
							p.Literal("."),
							p.Star(digit),
						),
					),
				),
				p.Optional(p.Concat(
					p.Set(charset.New([]byte{'e', 'E'})),
					p.Optional(p.Set(charset.New([]byte{'+', '-'}))),
					p.Plus(digit),
				)),
			),
			p.Concat(
				p.Plus(digit),
				p.Set(charset.New([]byte{'e', 'E'})),
				p.Optional(p.Set(charset.New([]byte{'+', '-'}))),
				p.Plus(digit),
			),
		),
	)

	word = p.Concat(
		p.Or(alpha, p.Literal("_")),
		p.Star(p.Or(alnum, p.Literal("_"))),
	)
)

func wordMatch(words ...string) p.Pattern {
	m := make(map[string]struct{})

	for _, w := range words {
		m[w] = struct{}{}
	}

	return p.Check(word, isa.MapChecker(m))
}
