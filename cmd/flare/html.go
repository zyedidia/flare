package main

import (
	"fmt"
	"html"
	"strings"

	"github.com/zyedidia/flare/theme"
)

type HTMLStyler struct {
	theme theme.Theme
	class bool
	name  string
}

func (st *HTMLStyler) Style(s, group string) string {
	if s != "" && group != "" {
		style := st.theme.Style(group)

		css := ""
		if style.Fg != nil {
			css += fmt.Sprintf("color:%s;", style.Fg.Hex())
		}
		if style.Attr&theme.AttrBold != 0 {
			css += "font-weight:bold;"
		}
		if style.Attr&theme.AttrUnderline != 0 {
			css += "text-decoration:underline;"
		}

		class := ""
		if st.class {
			class = fmt.Sprintf("class=\"%s\"", strings.ReplaceAll(group, ".", "-"))
		}

		if style != st.theme.Style("default") {
			return fmt.Sprintf("<span %s style=\"%s\">%s</span>", class, css, html.EscapeString(s))
		}
	}
	return html.EscapeString(s)
}

func (st *HTMLStyler) Pre() string {
	css := ""
	if st.theme["default"].Fg != nil {
		fg := st.theme["default"].Fg
		css += fmt.Sprintf("color:%s;", fg.Hex())
	}
	if st.theme["default"].Bg != nil {
		bg := st.theme["default"].Bg
		css += fmt.Sprintf("background-color:%s;", bg.Hex())
	}
	return fmt.Sprintf("<pre id=\"%s\" class=\"%s\" style=\"%s\">\n", "flare", "flare-"+st.name, css)
}
func (st *HTMLStyler) Post() string {
	return "</pre>\n"
}
