package flare

import (
	"embed"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/zyedidia/gpeg/pattern"
	p "github.com/zyedidia/gpeg/pattern"
	"github.com/zyedidia/gpeg/re"
	"github.com/zyedidia/gpeg/vm"
	"gopkg.in/yaml.v2"
)

type Language struct {
	Rules  map[string]Rule
	Tokens []string
}

type Rule struct {
	Pattern string
	Capture string
	Words   []string
	Include string
}

//go:embed languages/*.yaml
var builtin embed.FS
var custom map[string]func() ([]byte, error)

func AddLanguage(name string, loader func() ([]byte, error)) {
	custom[name] = loader
}

func LoadLanguage(name string) (*Language, error) {
	_, ok := custom[name]
	if ok {
		return loadCustomLanguage(name)
	}
	return loadBuiltinLanguage(name)
}

func loadBuiltinLanguage(name string) (*Language, error) {
	f, err := builtin.ReadFile(filepath.Join("languages", name+".yaml"))
	if err != nil {
		return nil, err
	}
	return loadYaml(f)
}

// assumes a custom language with 'name' exists
func loadCustomLanguage(name string) (*Language, error) {
	loader := custom[name]
	data, err := loader()
	if err != nil {
		return nil, err
	}
	return loadYaml(data)
}

func loadYaml(data []byte) (*Language, error) {
	l := new(Language)
	err := yaml.Unmarshal([]byte(data), l)
	return l, err
}

func (l *Language) Highlighter() (*Highlighter, error) {
	token, caps, err := l.pattern(0)
	if err != nil {
		return nil, err
	}

	grammar := make(map[string]p.Pattern)
	grammar["top"] = p.Star(p.Memo(p.Or(
		token,
		p.Concat(
			p.Any(1),
			p.Star(p.Concat(
				p.Not(token),
				p.Any(1),
			)),
		),
	)))

	prog, err := pattern.Compile(p.Grammar("top", grammar))
	if err != nil {
		return nil, err
	}

	return &Highlighter{
		code:     vm.Encode(prog),
		captures: caps,
	}, nil
}

func (l *Language) pattern(capid int) (p.Pattern, map[int]string, error) {
	caps := make(map[int]string)

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
		if rule.Pattern != "" {
			var err error
			patt, err = re.Compile(rule.Pattern)
			if err != nil {
				return nil, nil, fmt.Errorf("%s: %s: %w", name, strconv.Quote(rule.Pattern), err)
			}
		} else if rule.Include != "" {
			lang, err := LoadLanguage(rule.Include)
			if err != nil {
				return nil, nil, fmt.Errorf("%s: invalid include: %s: %w", name, rule.Include, err)
			}
			grammar, icaps, err := lang.pattern(capid)
			if err != nil {
				return nil, nil, fmt.Errorf("during include of %s: %w", rule.Include, err)
			}
			for range icaps {
				caps[capid] = icaps[capid]
				capid++
			}
			patt = grammar
		} else if rule.Words != nil {
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

	grammar["token"] = p.Or(tokens...)

	return p.Grammar("token", grammar), caps, nil
}
