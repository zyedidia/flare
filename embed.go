//go:build !flare_custom

package flare

import (
	"embed"
	"path/filepath"
)

//go:embed languages/*.lang
var langs embed.FS

// BuiltinLoader loads the grammar for the given language from the highlighters
// that are built into Flare. It has the signature expected by SetLoader, so a
// custom loader can fall back to the built-in highlighters by calling it:
//
//	flare.SetLoader(func(name string) ([]byte, error) {
//		if data, err := myLoader(name); err == nil {
//			return data, nil
//		}
//		return flare.BuiltinLoader(name)
//	})
//
// If Flare is built with the 'flare_custom' tag no grammars are embedded and
// BuiltinLoader always returns an error.
func BuiltinLoader(name string) ([]byte, error) {
	return langs.ReadFile(filepath.Join("languages", name+".lang"))
}

func init() {
	SetLoader(BuiltinLoader)
}
