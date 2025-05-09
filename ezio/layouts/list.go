package layouts

import (
	"gioui.org/layout"
	"github.com/rs/zerolog/log"
)

type ListOpts struct {
	ListPtr     *layout.List // [!] Must live across frames
	Axis        layout.Axis
	ScrollToEnd bool
	Alignment   layout.Alignment
	Length      int // Number of items
}

func List(opts ListOpts, item func(gtx layout.Context, index int) layout.Dimensions) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ListPtr == nil {
			log.Fatal().Msg("opts.List must be initialized!")
		}
		opts.ListPtr.Axis = opts.Axis
		opts.ListPtr.Alignment = opts.Alignment
		opts.ListPtr.ScrollToEnd = opts.ScrollToEnd
		return opts.ListPtr.Layout(gtx, opts.Length, item)
	}
}

func ListChild(children ...layout.Widget) func(gtx layout.Context, index int) layout.Dimensions {
	return func(gtx layout.Context, index int) layout.Dimensions {
		if index < 0 || index >= len(children) {
			// Out-Of-Bounds: Nothing to layout
			return layout.Dimensions{}
		}
		return children[index](gtx)
	}
}
