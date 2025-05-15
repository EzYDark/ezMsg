package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type RectOpts struct {
	W, H       unit.Dp
	MinW, MinH unit.Dp
	Color      color.NRGBA
}

func Rect(opts RectOpts) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		w := gtx.Dp(opts.W)
		h := gtx.Dp(opts.H)
		min_w := gtx.Dp(opts.MinW)
		min_h := gtx.Dp(opts.MinH)

		if w <= 0 {
			w = gtx.Constraints.Max.X
		}
		if h <= 0 {
			h = gtx.Constraints.Max.Y
		}

		if opts.MinW > 0 && w < min_w {
			w = min_w
		}
		if opts.MinH > 0 && h < min_h {
			h = min_h
		}

		size := image.Pt(w, h)
		paint.FillShape(gtx.Ops, opts.Color, clip.Rect{Max: size}.Op())
		return layout.Dimensions{Size: size}
	}
}
