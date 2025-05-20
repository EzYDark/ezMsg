package layouts

import (
	"gioui.org/layout"
)

func BackgroundBox(bgWidget layout.Widget, fgWidgets layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx, bgWidget, fgWidgets)
	}
}

// Background widget`s size based on the foreground widget`s size
func BackgroundStackBox(opts StackBoxOpts, bgWidget layout.Widget, fgWidgets layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return StackBox(opts,
			ExpandedChild(bgWidget),
			StackedChild(fgWidgets),
		)(gtx)
	}
}
