package flare

import (
	"embed"
	"path/filepath"

	"github.com/zyedidia/flare/syntax"
	"github.com/zyedidia/gpeg/pattern"
	p "github.com/zyedidia/gpeg/pattern"
	"github.com/zyedidia/gpeg/vm"
)

//go:embed languages/*.lang
var builtin embed.FS
var custom map[string]func() ([]byte, error)

func AddLanguage(name string, loader func() ([]byte, error)) {
	custom[name] = loader
}

func LoadHighlighter(name string) (*Highlighter, error) {
	data, err := loadData(name)
	if err != nil {
		return nil, err
	}
	return loadHighlighter(data)
}

func loadData(name string) ([]byte, error) {
	if loader, ok := custom[name]; ok {
		return loader()
	}
	return builtin.ReadFile(filepath.Join("languages", name+".lang"))
}

func loadHighlighter(data []byte) (*Highlighter, error) {
	capid := 0
	caps := make(map[int]string)

	capfn := func(patt p.Pattern, group string) p.Pattern {
		patt = p.Cap(patt, capid)
		caps[capid] = group
		capid++
		return patt
	}
	wordsfn := func(words ...string) p.Pattern {
		return wordMatch(words...)
	}

	var includefn func(lang string) p.Pattern
	includefn = func(lang string) p.Pattern {
		data, _ := loadData(lang)
		token, _ := syntax.Compile(string(data), syntax.CustomFns{
			Cap:     capfn,
			Words:   wordsfn,
			Include: includefn,
		})
		return token
	}

	token, err := syntax.Compile(string(data), syntax.CustomFns{
		Cap:     capfn,
		Words:   wordsfn,
		Include: includefn,
	})
	if err != nil {
		return nil, err
	}

	grammar := map[string]p.Pattern{
		"alpha":   alpha,
		"alnum":   alnum,
		"digit":   digit,
		"xdigit":  xdigit,
		"space":   space,
		"dec_num": dec_num,
		"hex_num": hex_num,
		"oct_num": oct_num,
		"integer": integer,
		"float":   float,
		"word":    word,

		"top": p.Star(p.Memo(p.Or(
			token,
			p.Concat(
				p.Any(1),
				p.Star(p.Concat(
					p.Not(token),
					p.Any(1),
				)),
			),
		))),
	}

	prog, err := pattern.Compile(p.Grammar("top", grammar))
	if err != nil {
		return nil, err
	}

	return &Highlighter{
		code:     vm.Encode(prog),
		captures: caps,
	}, nil
}
