package layouts

import "gioui.org/layout"

type DirectionBoxOpts struct {
	Direction layout.Direction
}

func DirectionBox(opts *DirectionBoxOpts, children ...layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		var alignDirection layout.Direction
		if opts != nil {
			alignDirection = opts.Direction
		}
		stackOpts := StackBoxOpts{
			Alignment: alignDirection,
		}

		content := StackedChild(children...)

		return StackBox(stackOpts, content)(gtx)
	}
}
