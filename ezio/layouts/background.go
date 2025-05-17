package layouts

import (
	"gioui.org/layout"
)

type BackgroundBoxOpts struct {
	Alignment layout.Direction
}

func BackgroundBox(opts BackgroundBoxOpts, bgWidget layout.Widget, fgWidgets ...layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		backgroundAsExpandedChild := ExpandedChild(bgWidget)
		foregroundAsStackedChild := StackedChild(fgWidgets...)

		stackBehaviorOptions := StackBoxOpts{
			Alignment: opts.Alignment,
		}

		return StackBox(stackBehaviorOptions, backgroundAsExpandedChild, foregroundAsStackedChild)(gtx)
	}
}
