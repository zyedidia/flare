package main

import (
	"fmt"
	"html"
	"strings"

	"github.com/zyedidia/flare/theme"
)

type HTMLStyler struct {
	theme theme.Theme
	name  string
}

func (st *HTMLStyler) Style(s, group string) string {
	style := st.theme.Style(group)

	css := ""
	if style.Fg != nil {
		css += fmt.Sprintf("color:rgb(%d,%d,%d);", style.Fg.R, style.Fg.G, style.Fg.B)
	}
	if style.Attr&theme.AttrBold != 0 {
		css += "font-weight:bold;"
	}

	if s != "" && group != "" {
		return fmt.Sprintf("<span class=\"%s\" style=\"%s\">%s</span>", strings.ReplaceAll(group, ".", "-"), css, html.EscapeString(s))
	}
	return s
}

func (st *HTMLStyler) Pre() string {
	css := "<style>"

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
