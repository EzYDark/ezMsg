package flex

import "gioui.org/layout"

func FlexChild(widgets ...layout.Widget) FlexChildOpts {
	return FlexChildOpts{Weight: 1, Widgets: widgets}
}

func FlexChildWeight(wt float32, widgets ...layout.Widget) FlexChildOpts {
	return FlexChildOpts{Weight: wt, Widgets: widgets}
}

func FlexChildStatic(opts FlexChildStaticOpts, widgets ...layout.Widget) FlexChildOpts {
	return FlexChildOpts{
		W:       opts.W,
		H:       opts.H,
		MinW:    opts.MinW,
		MinH:    opts.MinH,
		Widgets: widgets,
	}
}
