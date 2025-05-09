package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type RectOpts struct {
	W, H       int
	MinW, MinH int
	Color      color.NRGBA
}

func Rect(opts RectOpts) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		w := opts.W
		if w <= 0 {
			w = gtx.Constraints.Max.X
		}
		h := opts.H
		if h <= 0 {
			h = gtx.Constraints.Max.Y
		}
		if opts.MinW > 0 && w < opts.MinW {
			w = opts.MinW
		}
		if opts.MinH > 0 && h < opts.MinH {
			h = opts.MinH
		}
		size := image.Pt(w, h)
		paint.FillShape(gtx.Ops, opts.Color, clip.Rect{Max: size}.Op())
		return layout.Dimensions{Size: size}
	}
}
