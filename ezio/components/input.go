package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/rs/zerolog/log"
)

type InputOpts struct {
	EditorPtr  *widget.Editor // [!] Must live across frames
	ThemePtr   *material.Theme
	Hint       string
	W, H       int
	MinW, MinH int
}

func Input(opts InputOpts) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.EditorPtr == nil {
			log.Fatal().Msg("opts.EditorPtr must be initialized!")
		} else if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.Theme must be initialized!")
		}

		c := gtx.Constraints
		if opts.MinW > 0 && c.Min.X < opts.MinW {
			c.Min.X = opts.MinW
		}
		if opts.MinH > 0 && c.Min.Y < opts.MinH {
			c.Min.Y = opts.MinH
		}
		if opts.W > 0 {
			c.Min.X = opts.W
			c.Max.X = opts.W
		}
		if opts.H > 0 {
			c.Min.Y = opts.H
			c.Max.Y = opts.H
		}
		gtx.Constraints = c

		ed := material.Editor(opts.ThemePtr, opts.EditorPtr, opts.Hint)
		return ed.Layout(gtx)
	}
}
