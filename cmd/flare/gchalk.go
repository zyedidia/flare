package main

import (
	"github.com/jwalton/gchalk"
	"github.com/zyedidia/flare/theme"
)

type GChalkStyler struct {
	lvl   gchalk.ColorLevel
	theme theme.Theme
}

func (st *GChalkStyler) Style(s, group string) string {
	style := st.theme.Style(group)
	gc := gchalk.New(gchalk.ForceLevel(st.lvl))
	if style.Fg != nil {
		gc = gc.WithRGB(style.Fg.R, style.Fg.G, style.Fg.B)
	}
	return gc.StyleMust()(s)
}

func (st *GChalkStyler) Pre() string  { return "" }
func (st *GChalkStyler) Post() string { return "" }
