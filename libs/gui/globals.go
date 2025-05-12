package gui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
)

const (
	Horizontal layout.Axis = iota
	Vertical
)

var Fonts = gofont.Collection()
var MyTheme = func() *material.Theme {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(Fonts))
	return th
}()
