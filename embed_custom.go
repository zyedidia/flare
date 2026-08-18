//go:build flare_custom

package flare

import "fmt"

// BuiltinLoader always fails when Flare is built with the 'flare_custom' tag,
// because no grammars are embedded in that configuration. See the
// documentation of BuiltinLoader in embed.go.
func BuiltinLoader(name string) ([]byte, error) {
	return nil, fmt.Errorf("no built-in highlighter for language: %s (built with flare_custom)", name)
}
