package flex

import "gioui.org/layout"

type FlexChildOpts struct {
	Weight     float32
	MaxW, MaxH int
	MinW, MinH int
	W, H       int // Static size
	widgets    []layout.Widget
}

func FlexChild(opts *FlexChildOpts, widgets ...layout.Widget) FlexChildOpts {
	if opts == nil {
		opts = &FlexChildOpts{
			Weight:  1,
			widgets: widgets,
		}
	}

	opts.widgets = widgets
	return *opts
}
