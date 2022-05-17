//go:build !flare_custom

package flare

import (
	"embed"
	"path/filepath"
)

//go:embed languages/*.lang
var langs embed.FS

func init() {
	SetLoader(func(name string) ([]byte, error) {
		return langs.ReadFile(filepath.Join("languages", name+".lang"))
	})
}
