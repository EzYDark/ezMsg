// File: ezio/layouts/flex/flex.go
package flex

import "gioui.org/layout"

type FlexBoxOpts struct {
	Axis      layout.Axis
	Spacing   layout.Spacing
	Alignment layout.Alignment
}

func FlexBox(opts *FlexBoxOpts, children ...FlexChildOpts) layout.Widget {
	if opts == nil {
		opts = &FlexBoxOpts{}
	}

	return func(gtx layout.Context) layout.Dimensions {
		var kids []layout.FlexChild
		for _, ch := range children {
			ch_minWidth := gtx.Dp(ch.MinW)
			ch_minHeight := gtx.Dp(ch.MinH)
			ch_Width := gtx.Dp(ch.W)
			ch_Height := gtx.Dp(ch.H)

			composite_widget := func(gtx layout.Context) layout.Dimensions {
				c := gtx.Constraints
				if ch.MinW > 0 && c.Min.X < ch_minWidth {
					c.Min.X = ch_minWidth
				}
				if ch.MinH > 0 && c.Min.Y < ch_minHeight {
					c.Min.Y = ch_minHeight
				}
				if ch.W > 0 {
					c.Min.X = ch_Width
					c.Max.X = ch_Width
				}
				if ch.H > 0 {
					c.Min.Y = ch_Height
					c.Max.Y = ch_Height
				}
				gtx.Constraints = c
				var d layout.Dimensions
				for _, childW := range ch.widgets {
					d = childW(gtx)
				}
				return d
			}
			switch {
			case ch.W > 0 || ch.H > 0:
				kids = append(kids, layout.Rigid(composite_widget))
			case ch.Weight > 0:
				kids = append(kids, layout.Flexed(ch.Weight, composite_widget))
			default: // Child has Weight: 0 and no explicit W/H. Should be Rigid.
				kids = append(kids, layout.Rigid(composite_widget))
			}
		}
		return layout.Flex{Axis: opts.Axis, Spacing: opts.Spacing, Alignment: opts.Alignment}.Layout(gtx, kids...)
	}
}
