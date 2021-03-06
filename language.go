package flare

import (
	"fmt"

	"github.com/zyedidia/flare/syntax"
	"github.com/zyedidia/gpeg/isa"
	p "github.com/zyedidia/gpeg/pattern"
	"github.com/zyedidia/gpeg/vm"
)

var builtin func(name string) ([]byte, error)
var custom map[string]func() ([]byte, error)

func SetLoader(loader func(name string) ([]byte, error)) {
	builtin = loader
}

// AddLanguage adds support for a new language. The 'loader' should return the
// highlighting grammar when called.
func AddLanguage(name string, loader func() ([]byte, error)) {
	custom[name] = loader
}

// LoadHighlighterBuiltin compiles and loads the given highlighter. Use 'memo'
// if the highlighter will be used in an incremental setting (editor).
func LoadHighlighterBuiltin(name string, memo bool) (*Highlighter, error) {
	data, err := loadData(name)
	if err != nil {
		return &empty, err
	}
	return LoadHighlighter(name, data, memo)
}

func loadData(name string) ([]byte, error) {
	if loader, ok := custom[name]; ok {
		return loader()
	}
	if builtin != nil {
		return builtin(name)
	}
	return nil, fmt.Errorf("no highlighter for language: %s", name)
}

// LoadHighlighter loads a highlighter from some data.
func LoadHighlighter(lang string, data []byte, memo bool) (*Highlighter, error) {
	capid := 0
	refid := 0
	caps := make(map[int]string)
	refs := make(map[string]int)

	capfn := func(patt p.Pattern, group string) p.Pattern {
		patt = p.Cap(patt, capid)
		caps[capid] = group
		capid++
		return patt
	}
	wordsfn := func(words ...string) p.Pattern {
		return wordMatch(words...)
	}
	br := isa.NewBackRef()
	reffn := func(patt p.Pattern, group string) p.Pattern {
		patt = p.CheckFlags(patt, br, refid, int(isa.RefDef))
		refs[group] = refid
		refid++
		return patt
	}
	backfn := func(ref string) p.Pattern {
		return p.CheckFlags(&p.EmptyNode{}, br, refs[ref], int(isa.RefUse))
	}
	imports := map[string]p.Pattern{
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
	}

	var includefn func(lang string) p.Pattern
	includefn = func(lang string) p.Pattern {
		data, _ := loadData(lang)
		token, _ := syntax.Compile(lang+"_", string(data), syntax.CustomFns{
			Cap:     capfn,
			Words:   wordsfn,
			Include: includefn,
			Ref:     reffn,
			Back:    backfn,
			Imports: imports,
		})
		return token
	}

	token, err := syntax.Compile(lang+"_", string(data), syntax.CustomFns{
		Cap:     capfn,
		Words:   wordsfn,
		Include: includefn,
		Ref:     reffn,
		Back:    backfn,
		Imports: imports,
	})
	if err != nil {
		return &empty, err
	}

	top := p.Or(
		token,
		p.Concat(
			p.Any(1),
			p.Star(p.Concat(
				p.Not(token),
				p.Any(1),
			)),
		),
	)
	if memo {
		top = p.Memo(top)
	}

	top = p.Star(top)

	grammar := map[string]p.Pattern{
		"top": top,
	}

	prog, err := p.Compile(p.Grammar("top", grammar))
	if err != nil {
		return &empty, err
	}

	return &Highlighter{
		code:     vm.Encode(prog),
		captures: caps,
	}, nil
}
