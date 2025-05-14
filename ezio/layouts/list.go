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

func ListBox(opts ListOpts, items ...layout.ListElement) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ListPtr == nil {
			log.Fatal().Msg("opts.List must be initialized!")
		}
		opts.ListPtr.Axis = opts.Axis
		opts.ListPtr.Alignment = opts.Alignment
		opts.ListPtr.ScrollToEnd = opts.ScrollToEnd

		length := opts.Length
		if length == 0 {
			length = len(items)
		}

		// Create a function that selects the appropriate item function based on index
		itemFunc := func(gtx layout.Context, index int) layout.Dimensions {
			if index < 0 || index >= len(items) {
				// Out-Of-Bounds: Nothing to layout
				return layout.Dimensions{}
			}
			return items[index](gtx, 0) // Call the appropriate item function
		}

		return opts.ListPtr.Layout(gtx, length, itemFunc)
	}
}

func ListChild(children ...layout.Widget) layout.ListElement {
	return func(gtx layout.Context, index int) layout.Dimensions {
		if index < 0 || index >= len(children) {
			// Out-Of-Bounds: Nothing to layout
			return layout.Dimensions{}
		}
		return children[index](gtx)
	}
}
