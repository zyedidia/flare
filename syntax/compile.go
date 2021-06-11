package syntax

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/zyedidia/gpeg/charset"
	"github.com/zyedidia/gpeg/memo"
	"github.com/zyedidia/gpeg/pattern"
	"github.com/zyedidia/gpeg/vm"
)

var parser vm.VMCode

func init() {
	prog := pattern.MustCompile(pattern.Grammar("Pattern", grammar))
	parser = vm.Encode(prog)
}

type CustomFns struct {
	Cap     func(p pattern.Pattern, group string) pattern.Pattern
	Words   func(words ...string) pattern.Pattern
	Include func(lang string) pattern.Pattern
}

func compile(root *memo.Capture, s string, fns CustomFns) pattern.Pattern {
	var p pattern.Pattern
	switch root.Id() {
	case idPattern:
		p = compile(root.Child(0), s, fns)
	case idGrammar:
		nonterms := make(map[string]pattern.Pattern)
		it := root.ChildIterator(0)
		for c := it(); c != nil; c = it() {
			k, v := compileDef(c, s, fns)
			nonterms[k] = v
		}
		p = pattern.Grammar("token", nonterms)
	case idExpression:
		alternations := make([]pattern.Pattern, 0, root.NumChildren())
		it := root.ChildIterator(0)
		for c := it(); c != nil; c = it() {
			alternations = append(alternations, compile(c, s, fns))
		}
		p = pattern.Or(alternations...)
	case idSequence:
		concats := make([]pattern.Pattern, 0, root.NumChildren())
		it := root.ChildIterator(0)
		for c := it(); c != nil; c = it() {
			concats = append(concats, compile(c, s, fns))
		}
		p = pattern.Concat(concats...)
	case idPrefix:
		c := root.Child(0)
		switch c.Id() {
		case idAND:
			p = pattern.And(compile(root.Child(1), s, fns))
		case idNOT:
			p = pattern.Not(compile(root.Child(1), s, fns))
		default:
			p = compile(root.Child(0), s, fns)
		}
	case idSuffix:
		if root.NumChildren() == 2 {
			c := root.Child(1)
			switch c.Id() {
			case idQUESTION:
				p = pattern.Optional(compile(root.Child(0), s, fns))
			case idSTAR:
				p = pattern.Star(compile(root.Child(0), s, fns))
			case idPLUS:
				p = pattern.Plus(compile(root.Child(0), s, fns))
			}
		} else {
			p = compile(root.Child(0), s, fns)
		}
	case idPrimary:
		switch root.Child(0).Id() {
		case idCAP:
			cpatt := compile(root.Child(1), s, fns)
			group := literal(root.Child(2), s)
			p = fns.Cap(cpatt, group)
		case idINCLUDE:
			lang := literal(root.Child(1), s)
			p = fns.Include(lang)
		case idWORDS:
			it := root.ChildIterator(0)
			words := make([]string, root.NumChildren()-1)
			for c := it(); c != nil; c = it() {
				if c.Id() != idLiteral {
					continue
				}
				words = append(words, literal(c, s))
			}
			p = fns.Words(words...)
		case idIdentifier, idLiteral, idClass:
			p = compile(root.Child(0), s, fns)
		case idOPEN:
			p = compile(root.Child(1), s, fns)
		case idDOT:
			p = pattern.Any(1)
		}
	case idLiteral:
		p = pattern.Literal(literal(root, s))
	case idClass:
		var set charset.Set
		if root.NumChildren() <= 0 {
			break
		}
		complement := false
		if root.Child(0).Id() == idCARAT {
			complement = true
		}
		it := root.ChildIterator(0)
		i := 0
		for c := it(); c != nil; c = it() {
			if i == 0 && complement {
				i++
				continue
			}
			set = set.Add(compileSet(c, s))
		}
		if complement {
			set = set.Complement()
		}
		p = pattern.Set(set)
	case idIdentifier:
		p = pattern.NonTerm(parseId(root, s))
	}
	return p
}

var special = map[byte]byte{
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'\'': '\'',
	'"':  '"',
	'[':  '[',
	']':  ']',
	'\\': '\\',
	'-':  '-',
}

func parseChar(char string) byte {
	switch char[0] {
	case '\\':
		for k, v := range special {
			if char[1] == k {
				return v
			}
		}

		i, _ := strconv.ParseInt(string(char[1:]), 8, 8)
		return byte(i)
	default:
		return char[0]
	}
}

func parseId(root *memo.Capture, s string) string {
	ident := &bytes.Buffer{}
	it := root.ChildIterator(0)
	for c := it(); c != nil; c = it() {
		ident.WriteString(s[c.Start():c.End()])
	}
	return ident.String()
}

func literal(root *memo.Capture, s string) string {
	lit := &bytes.Buffer{}
	it := root.ChildIterator(0)
	for c := it(); c != nil; c = it() {
		lit.WriteByte(parseChar(s[c.Start():c.End()]))
	}
	return lit.String()
}

func compileDef(root *memo.Capture, s string, fns CustomFns) (string, pattern.Pattern) {
	id := root.Child(0)
	exp := root.Child(1)
	return parseId(id, s), compile(exp, s, fns)
}

func compileSet(root *memo.Capture, s string) charset.Set {
	switch root.NumChildren() {
	case 1:
		c := root.Child(0)
		return charset.New([]byte{parseChar(s[c.Start():c.End()])})
	case 2:
		c1, c2 := root.Child(0), root.Child(1)
		return charset.Range(parseChar(s[c1.Start():c1.End()]), parseChar(s[c2.Start():c2.End()]))
	}
	return charset.Set{}
}

func Compile(s string, fns CustomFns) (pattern.Pattern, error) {
	match, n, ast, errs := parser.Exec(strings.NewReader(s), memo.NoneTable{})
	if len(errs) != 0 {
		return nil, errs[0]
	}
	if !match {
		return nil, fmt.Errorf("Invalid PEG: failed at %d", n)
	}

	return compile(ast.Child(0), s, fns), nil
}

func MustCompile(s string, fns CustomFns) pattern.Pattern {
	p, err := Compile(s, fns)
	if err != nil {
		panic(err)
	}
	return p
}
