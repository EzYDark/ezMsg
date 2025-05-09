package flex

import "gioui.org/layout"

func Flex(opts FlexOpts, children ...FlexChildOpts) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		var kids []layout.FlexChild
		for _, ch := range children {
			composite_widget := func(gtx layout.Context) layout.Dimensions {
				c := gtx.Constraints
				if ch.MinW > 0 && c.Min.X < ch.MinW {
					c.Min.X = ch.MinW
				}
				if ch.MinH > 0 && c.Min.Y < ch.MinH {
					c.Min.Y = ch.MinH
				}
				if ch.W > 0 {
					c.Min.X = ch.W
					c.Max.X = ch.W
				}
				if ch.H > 0 {
					c.Min.Y = ch.H
					c.Max.Y = ch.H
				}
				gtx.Constraints = c
				var d layout.Dimensions
				for _, childW := range ch.Widgets {
					d = childW(gtx)
				}
				return d
			}
			switch {
			case ch.W > 0 || ch.H > 0:
				kids = append(kids, layout.Rigid(composite_widget))
			case ch.Weight > 0:
				kids = append(kids, layout.Flexed(ch.Weight, composite_widget))
			default:
				kids = append(kids, layout.Flexed(1, composite_widget))
			}
		}
		return layout.Flex{Axis: opts.Alignment}.Layout(gtx, kids...)
	}
}
