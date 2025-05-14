package layouts

import "gioui.org/layout"

type DirectionBoxOpts struct {
	Direction layout.Direction
}

func DirectionBox(opts *DirectionBoxOpts, children ...layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts == nil {
			opts = &DirectionBoxOpts{}
		}

		composite_widget := func(gtx layout.Context) layout.Dimensions {
			var d layout.Dimensions
			for _, w := range children {
				d = w(gtx)
			}
			return d
		}
		return layout.Direction(opts.Direction).Layout(gtx, composite_widget)
	}
}
