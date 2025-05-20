package layouts

import (
	"gioui.org/layout"
)

type BackgroundBoxOpts struct {
	Alignment layout.Direction
	// Add other options like spacing if needed in the future
}

func BackgroundBox(opts BackgroundBoxOpts, bgWidget layout.Widget, fgWidgets layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx, bgWidget, fgWidgets)
	}
}
