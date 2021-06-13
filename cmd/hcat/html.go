package main

import (
	"fmt"
	"strings"

	"github.com/zyedidia/flare/theme"
)

type HTMLStyler struct {
	theme theme.Theme
	name  string
}

func (st *HTMLStyler) Style(s, group string) string {
	if s != "" && group != "" {
		return fmt.Sprintf("<span class=\"%s\">%s</span>", strings.ReplaceAll(group, ".", "-"), s)
	}
	return s
}

func (st *HTMLStyler) Pre() string {
	css := "<style>"
	for k, style := range st.theme {
		internal := ""
		if style.Fg != nil {
			internal += fmt.Sprintf("color:rgb(%d,%d,%d);", style.Fg.R, style.Fg.G, style.Fg.B)
		}
		if style.Attr&theme.AttrBold != 0 {
			internal += "font-weight:bold;"
		}
		css += fmt.Sprintf(".%s { %s } ", strings.ReplaceAll(k, ".", "-"), internal)
	}

	if st.theme["default"].Fg != nil {
		fg := st.theme["default"].Fg
		css += fmt.Sprintf("#%s { color:rgb(%d,%d,%d); }", "hcat", fg.R, fg.G, fg.B)
	}
	if st.theme["default"].Bg != nil {
		bg := st.theme["default"].Bg
		css += fmt.Sprintf("#%s { background-color:rgb(%d,%d,%d); }", "hcat", bg.R, bg.G, bg.B)
	}

	css += "</style>"
	return css + fmt.Sprintf("\n<pre id=\"%s\" class=\"%s\">\n", "hcat", "hcat-"+st.name)
}
func (st *HTMLStyler) Post() string {
	return "</pre>\n"
}
