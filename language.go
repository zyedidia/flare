package flare

import (
	"fmt"
	"strconv"

	p "github.com/zyedidia/gpeg/pattern"
	"github.com/zyedidia/gpeg/re"
	"gopkg.in/yaml.v2"
)

type Language struct {
	Filetype string
	Rules    map[string]Rule
	Tokens   []string
}

type Rule struct {
	Pattern string
	Capture string
	Words   []string
}

func LoadLanguage(data []byte) (*Language, error) {
	l := new(Language)
	err := yaml.Unmarshal([]byte(data), l)
	return l, err
}

func (l *Language) Highlighter() (*Highlighter, error) {
	caps := make(map[int]string)
	capid := 0

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
	}

	for name, rule := range l.Rules {
		var patt p.Pattern
		if rule.Words == nil {
			var err error
			patt, err = re.Compile(rule.Pattern)
			if err != nil {
				return nil, fmt.Errorf("%s: %s: %w", name, strconv.Quote(rule.Pattern), err)
			}
		} else {
			patt = wordMatch(rule.Words...)
		}

		if rule.Capture != "" {
			patt = p.Cap(patt, capid)
			caps[capid] = rule.Capture
			capid++
		}
		grammar[name] = patt
	}

	var tokens []p.Pattern
	for _, t := range l.Tokens {
		tokens = append(tokens, p.NonTerm(t))
	}

	grammar["top"] = p.Star(p.Memo(p.Or(
		p.NonTerm("token"),
		p.Concat(
			p.Any(1),
			p.Star(p.Concat(
				p.Not(p.NonTerm("token")),
				p.Any(1),
			)),
		),
	)))
	grammar["token"] = p.Or(tokens...)

	return &Highlighter{
		Grammar:  p.Grammar("top", grammar),
		Captures: caps,
	}, nil
}
