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
		w := gtx.Dp(opts.W)
		h := gtx.Dp(opts.H)
		min_w := gtx.Dp(opts.MinW)
		min_h := gtx.Dp(opts.MinH)

		if opts.EditorPtr == nil {
			log.Fatal().Msg("opts.EditorPtr must be initialized!")
		} else if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.Theme must be initialized!")
		}

		c := gtx.Constraints
		if opts.MinW > 0 && c.Min.X < min_w {
			c.Min.X = min_w
		}
		if opts.MinH > 0 && c.Min.Y < min_h {
			c.Min.Y = min_h
		}
		if opts.W > 0 {
			c.Min.X = w
			c.Max.X = w
		}
		if opts.H > 0 {
			c.Min.Y = h
			c.Max.Y = h
		}
		gtx.Constraints = c

		ed := material.Editor(opts.ThemePtr, opts.EditorPtr, opts.Hint)
		return ed.Layout(gtx)
	}
}
