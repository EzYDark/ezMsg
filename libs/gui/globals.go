package gui

import (
	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/ezydark/ezMsg/libs/gui/fonts"
)

var Fonts []text.FontFace
var MyTheme = func() *material.Theme {
	// URLs for the fonts
	regularFontURL := "https://pub-d415488005c546cab60d54ee8a268200.r2.dev/CallunaSans-Regular.otf"
	boldFontURL := "https://pub-d415488005c546cab60d54ee8a268200.r2.dev/CallunaSans-Bold.otf"
	italicFontURL := "https://pub-d415488005c546cab60d54ee8a268200.r2.dev/CallunaSans-Italic.otf"
	boldItalicFontURL := "https://pub-d415488005c546cab60d54ee8a268200.r2.dev/CallunaSansBoldItalic.otf"

	// Why: Load the regular font from the specified URL.
	callunaRegularBytes := fonts.LoadFontURL(regularFontURL)
	Fonts = fonts.AppendFont(Fonts,
		font.Font{
			Typeface: "CallunaSans",
		},
		callunaRegularBytes)

	// Why: Load the bold font from the specified URL.
	callunaBoldBytes := fonts.LoadFontURL(boldFontURL)
	Fonts = fonts.AppendFont(Fonts,
		font.Font{
			Typeface: "CallunaSans",
			Weight:   font.Bold,
		},
		callunaBoldBytes)

	// Why: Load the italic font from the specified URL.
	callunaItalicBytes := fonts.LoadFontURL(italicFontURL)
	Fonts = fonts.AppendFont(Fonts, font.Font{Typeface: "CallunaSans", Style: font.Italic}, callunaItalicBytes)

	// Why: Load the bold italic font from the specified URL.
	callunaBoldItalicBytes := fonts.LoadFontURL(boldItalicFontURL)
	Fonts = fonts.AppendFont(Fonts,
		font.Font{
			Typeface: "CallunaSans",
			Weight:   font.Bold,
			Style:    font.Italic,
		}, callunaBoldItalicBytes)

	th := material.NewTheme()
	// Why: Initialize the text shaper with the loaded font collection.
	th.Shaper = text.NewShaper(text.WithCollection(Fonts))
	return th
}()

// var Fonts = gofont.Collection()
// var MyTheme = func() *material.Theme {
// 	th := material.NewTheme()
// 	th.Shaper = text.NewShaper(text.WithCollection(Fonts))
// 	return th
// }()
