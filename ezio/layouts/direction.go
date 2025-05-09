package layouts

import "gioui.org/layout"

const (
	NW layout.Direction = iota
	N
	NE
	E
	SE
	S
	SW
	W
	Center
)

type DirectionOpts struct {
	Direction layout.Direction
}

// Direction wraps one or more widgets in the given layout.Direction.
// It runs all the children in order inside the directional context.
func Direction(opts DirectionOpts, children ...layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// build one composite_widget widget out of all children
		composite_widget := func(gtx layout.Context) layout.Dimensions {
			var d layout.Dimensions
			for _, w := range children {
				d = w(gtx)
			}
			return d
		}
		// invoke Gio's Direction.Layout
		return layout.Direction(opts.Direction).Layout(gtx, composite_widget)
	}
}
