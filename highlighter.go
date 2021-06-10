package flare

import p "github.com/zyedidia/gpeg/pattern"

type Highlighter struct {
	Grammar  p.Pattern
	Captures map[int]string
}
