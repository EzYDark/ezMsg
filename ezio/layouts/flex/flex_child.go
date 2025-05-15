package flex

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type FlexChildOpts struct {
	Weight     float32
	MaxW, MaxH unit.Dp
	MinW, MinH unit.Dp
	W, H       unit.Dp // Static size
	widgets    []layout.Widget
}

func FlexChild(opts *FlexChildOpts, widgets ...layout.Widget) FlexChildOpts {
	if opts == nil {
		opts = &FlexChildOpts{
			widgets: widgets,
		}
	} else {
		opts.widgets = widgets
	}
	return *opts
}
