package theme

import (
	"fmt"
)

type Style struct {
	Fg, Bg *Color
	Attr   AttrMask
}

type Color struct {
	R, G, B uint8
}

func (c *Color) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var hex string
	err := unmarshal(&hex)
	if err != nil {
		return err
	}
	clr, err := HexColor(hex)
	if err != nil {
		return err
	}
	*c = clr
	return nil
}

func (c Color) Hex() string {
	return "#" + fmt.Sprintf("%02x", c.R) + fmt.Sprintf("%02x", c.G) + fmt.Sprintf("%02x", c.B)
}

func HexColor(s string) (c Color, err error) {
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}

	return c, err
}

func MustHexColor(s string) Color {
	c, err := HexColor(s)
	if err != nil {
		panic(err)
	}
	return c
}

type AttrMask int

const (
	AttrBold = 1 << iota
	AttrBlink
	AttrReverse
	AttrUnderline
	AttrDim
	AttrItalic
	AttrStrikethrough
	AttrHidden
	AttrNone AttrMask = 0
)

func Attr(s string) (AttrMask, error) {
	switch s {
	case "bold":
		return AttrBold, nil
	case "blink":
		return AttrBlink, nil
	case "reverse":
		return AttrReverse, nil
	case "underline":
		return AttrUnderline, nil
	case "dim":
		return AttrDim, nil
	case "italic":
		return AttrItalic, nil
	case "strikethrough":
		return AttrStrikethrough, nil
	case "hidden":
		return AttrHidden, nil
	}
	return AttrNone, fmt.Errorf("invalid attribute: %s", s)
}

func (a *AttrMask) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var attrs []string
	err := unmarshal(&attrs)
	if err != nil {
		return err
	}

	for _, s := range attrs {
		attr, err := Attr(s)
		if err != nil {
			return err
		}
		*a |= attr
	}
	return nil
}
