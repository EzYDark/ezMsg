package gui

import (
	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/ezydark/ezMsg/app/db"
	"github.com/ezydark/ezMsg/libs/gui/fonts"
	"github.com/ezydark/ezMsg/libs/gui/pages"
)

type ezAppState struct {
	CurrentPage pages.Page
	LoggedUser  *db.User
}

var DBPtr = db.InitDB()

// Default state on app startup
var AppState = &ezAppState{
	CurrentPage: pages.ChatPage,
	LoggedUser:  &DBPtr.RegisteredUsers[0], // TODO: Change dynamically
}

var Fonts []text.FontFace
var FontsNerd = func() []text.FontFace {
	var nerd_fonts []text.FontFace

	bytes := fonts.LoadFontURL("https://pub-d415488005c546cab60d54ee8a268200.r2.dev/0xProto/0xProtoNerdFont-Regular.ttf")
	nerd_fonts = fonts.AppendFont(nerd_fonts, font.Font{Typeface: "0xProto", Style: font.Regular}, bytes)

	return nerd_fonts
}()

var MyTheme = func() *material.Theme {
	regularFontURL := "https://pub-d415488005c546cab60d54ee8a268200.r2.dev/CallunaSans-Regular.otf"
	boldFontURL := "https://pub-d415488005c546cab60d54ee8a268200.r2.dev/CallunaSans-Bold.otf"
	italicFontURL := "https://pub-d415488005c546cab60d54ee8a268200.r2.dev/CallunaSans-Italic.otf"
	boldItalicFontURL := "https://pub-d415488005c546cab60d54ee8a268200.r2.dev/CallunaSansBoldItalic.otf"

	callunaRegularBytes := fonts.LoadFontURL(regularFontURL)
	Fonts = fonts.AppendFont(Fonts,
		font.Font{
			Typeface: "CallunaSans",
		},
		callunaRegularBytes)

	callunaBoldBytes := fonts.LoadFontURL(boldFontURL)
	Fonts = fonts.AppendFont(Fonts,
		font.Font{
			Typeface: "CallunaSans",
			Weight:   font.Bold,
		},
		callunaBoldBytes)

	callunaItalicBytes := fonts.LoadFontURL(italicFontURL)
	Fonts = fonts.AppendFont(Fonts, font.Font{Typeface: "CallunaSans", Style: font.Italic}, callunaItalicBytes)

	callunaBoldItalicBytes := fonts.LoadFontURL(boldItalicFontURL)
	Fonts = fonts.AppendFont(Fonts,
		font.Font{
			Typeface: "CallunaSans",
			Weight:   font.Bold,
			Style:    font.Italic,
		}, callunaBoldItalicBytes)

	var allFontsCollection []text.FontFace
	allFontsCollection = append(allFontsCollection, Fonts...)     // Add CallunaSans fonts
	allFontsCollection = append(allFontsCollection, FontsNerd...) // Add Nerd fonts

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(allFontsCollection))
	return th
}()
