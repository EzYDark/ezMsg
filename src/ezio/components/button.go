package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/rs/zerolog/log"
)

type ButtonOpts struct {
	ButtonPtr *widget.Clickable // [!] Must live across frames
	ThemePtr  *material.Theme
	Text      string
}

func Button(opts ButtonOpts) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ButtonPtr == nil {
			log.Fatal().Msg("opts.ButtonPtr must be initialized!")
		} else if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.Theme must be initialized!")
		}
		btn := material.Button(opts.ThemePtr, opts.ButtonPtr, opts.Text)
		btn.CornerRadius = 0
		return btn.Layout(gtx)
	}
}
