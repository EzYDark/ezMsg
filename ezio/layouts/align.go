package layouts

import "gioui.org/layout"

type AlignOpts = StackBoxOpts

func Align(opts *AlignOpts, children ...layout.Widget) layout.Widget {
	if opts == nil {
		opts = &AlignOpts{}
	}

	return func(gtx layout.Context) layout.Dimensions {
		content := StackedChild(children...)

		return StackBox(opts, content)(gtx)
	}
}
