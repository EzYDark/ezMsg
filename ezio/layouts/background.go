package layouts

import (
	"gioui.org/layout"
)

// Lays out the background first. Then, lays out the foreground, ensuring it's at least as large as the background.
func BackgroundBox(bgWidget layout.Widget, fgWidgets layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx, bgWidget, fgWidgets)
	}
}

// Background (ExpandedChild) takes its minimum size from the Foreground (StackedChild), then expands to fill.
func BackgroundStackBox(opts *StackBoxOpts, bgWidget layout.Widget, fgWidgets layout.Widget) layout.Widget {
	if opts == nil {
		opts = &StackBoxOpts{}
	}

	return func(gtx layout.Context) layout.Dimensions {
		return StackBox(opts,
			ExpandedChild(bgWidget),
			StackedChild(fgWidgets),
		)(gtx)
	}
}
