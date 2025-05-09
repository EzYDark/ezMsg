package components

import (
	"gioui.org/layout"
	"gioui.org/x/component"
	"github.com/rs/zerolog/log"
)

type ContextMenuOpts struct {
	ContextArea *component.ContextArea // [!] Must live across frames
}

func ContextMenu(opts ContextMenuOpts, contextualWidget layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ContextArea == nil {
			log.Fatal().Msg("ContextMenuOpts.ContextArea must be initialized and provided as a pointer!")
			return layout.Dimensions{}
		}
		return opts.ContextArea.Layout(gtx, contextualWidget)
	}
}
