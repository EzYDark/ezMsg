package layouts

import "gioui.org/layout"

type StackBoxOpts = layout.Stack

func StackBox(opts *StackBoxOpts, children ...layout.StackChild) layout.Widget {
	if opts == nil {
		opts = &StackBoxOpts{}
	}

	return func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{Alignment: opts.Alignment}.Layout(gtx, children...)
	}
}

func StackedChild(widgets ...layout.Widget) layout.StackChild {
	if len(widgets) == 0 {
		return layout.Stacked(func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{} })
	}
	if len(widgets) > 1 {
		compositeWidget := func(gtx layout.Context) layout.Dimensions {
			dims := layout.Dimensions{}
			for _, w := range widgets {

				d := w(gtx)
				if d.Size.X > dims.Size.X {
					dims.Size.X = d.Size.X
				}
				if d.Size.Y > dims.Size.Y {
					dims.Size.Y = d.Size.Y
				}
				if d.Baseline > dims.Baseline {
					dims.Baseline = d.Baseline
				}
			}
			return dims
		}
		return layout.Stacked(compositeWidget)
	}
	return layout.Stacked(widgets[0])
}

func ExpandedChild(widgets ...layout.Widget) layout.StackChild {
	if len(widgets) == 0 {
		return layout.Expanded(func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{} })
	}
	if len(widgets) > 1 {
		compositeWidget := func(gtx layout.Context) layout.Dimensions {
			dims := layout.Dimensions{}
			for _, w := range widgets {
				d := w(gtx)
				if d.Size.X > dims.Size.X {
					dims.Size.X = d.Size.X
				}
				if d.Size.Y > dims.Size.Y {
					dims.Size.Y = d.Size.Y
				}
				if d.Baseline > dims.Baseline {
					dims.Baseline = d.Baseline
				}
			}
			return dims
		}
		return layout.Expanded(compositeWidget)
	}
	return layout.Expanded(widgets[0])
}
