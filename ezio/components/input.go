package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/rs/zerolog/log"
)

type InputOpts struct {
	EditorPtr  *widget.Editor // [!] Must live across frames
	ThemePtr   *material.Theme
	Hint       string
	W, H       unit.Dp
	MinW, MinH unit.Dp
}

func Input(opts InputOpts) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.EditorPtr == nil {
			log.Fatal().Msg("opts.EditorPtr must be initialized!")
		} else if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.Theme must be initialized!")
		}

		ed := material.Editor(opts.ThemePtr, opts.EditorPtr, opts.Hint)
		return ed.Layout(gtx)
	}
}
