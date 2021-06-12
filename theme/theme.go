package theme

import (
	"embed"
	"errors"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:embed themes/*.yaml
var builtin embed.FS
var custom map[string]func() ([]byte, error)

type Theme map[string]Style

func AddTheme(name string, loader func() ([]byte, error)) {
	custom[name] = loader
}

func LoadTheme(name string) (Theme, error) {
	data, err := loadData(name)
	if err != nil {
		return nil, err
	}
	return NewThemeFromYaml(data)
}

func loadData(name string) ([]byte, error) {
	if loader, ok := custom[name]; ok {
		return loader()
	}
	return builtin.ReadFile(filepath.Join("themes", name+".yaml"))
}

func NewThemeFromYaml(data []byte) (Theme, error) {
	var t Theme
	err := yaml.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return NewTheme(t)
}

func NewTheme(rules map[string]Style) (Theme, error) {
	_, ok := rules["default"]
	if !ok {
		return nil, errors.New("no default style found")
	}
	return rules, nil
}

func (t Theme) Style(group string) Style {
	if r, ok := t[group]; ok {
		return r
	}
	i := strings.LastIndexByte(group, '.')
	if i == -1 {
		return t["default"]
	}
	return t.Style(group[:i])
}
