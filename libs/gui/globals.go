package gui

import (
	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/ezydark/ezMsg/libs/gui/fonts"
)

var Fonts []text.FontFace
var MyTheme = func() *material.Theme {
	callunaRegularBytes := fonts.LoadFontFile("./libs/gui/fonts/CallunaSans-Regular.otf")
	Fonts = fonts.AppendFont(Fonts,
		font.Font{
			Typeface: "CallunaSans",
		},
		callunaRegularBytes)

	callunaBoldBytes := fonts.LoadFontFile("./libs/gui/fonts/CallunaSans-Bold.otf")
	Fonts = fonts.AppendFont(Fonts,
		font.Font{
			Typeface: "CallunaSans",
			Weight:   font.Bold,
		},
		callunaBoldBytes)

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(Fonts))
	return th
}()

// var Fonts = gofont.Collection()
// var MyTheme = func() *material.Theme {
// 	th := material.NewTheme()
// 	th.Shaper = text.NewShaper(text.WithCollection(Fonts))
// 	return th
// }()
