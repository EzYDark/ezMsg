package layouts

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type MarginOpts struct {
	// Set individual margins
	Top, Bottom, Left, Right unit.Dp

	// Or set all margins the same
	All unit.Dp
}

func Margin(opts MarginOpts, child layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.All != 0 {
			inset := layout.Inset{
				Top:    opts.All,
				Bottom: opts.All,
				Left:   opts.All,
				Right:  opts.All,
			}
			return inset.Layout(gtx, child)
		} else {
			inset := layout.Inset{
				Top:    opts.Top,
				Bottom: opts.Bottom,
				Left:   opts.Left,
				Right:  opts.Right,
			}
			return inset.Layout(gtx, child)
		}
	}
}
