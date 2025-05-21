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

func Margin(opts *MarginOpts, children ...layout.Widget) layout.Widget {
	if opts == nil {
		opts = &MarginOpts{}
	}

	return func(gtx layout.Context) layout.Dimensions {
		// Create a composite widget from all children
		composite_widget := func(gtx layout.Context) layout.Dimensions {
			var d layout.Dimensions
			for _, w := range children {
				d = w(gtx)
			}
			return d
		}

		// Apply margin to the composite widget
		if opts.All != 0 {
			inset := layout.Inset{
				Top:    opts.All,
				Bottom: opts.All,
				Left:   opts.All,
				Right:  opts.All,
			}
			return inset.Layout(gtx, composite_widget)
		} else {
			inset := layout.Inset{
				Top:    opts.Top,
				Bottom: opts.Bottom,
				Left:   opts.Left,
				Right:  opts.Right,
			}
			return inset.Layout(gtx, composite_widget)
		}
	}
}
