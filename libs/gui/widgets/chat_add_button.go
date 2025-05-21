package widgets

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/x/richtext"
	. "github.com/ezydark/ezMsg/ezio"
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/ezydark/ezMsg/libs/gui/pages"
	"github.com/rs/zerolog/log"
)

func AddButton(buttonState *richtext.InteractiveText) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		for {
			_, event, ok := buttonState.Update(gtx)
			if !ok {
				break
			}
			switch event.Type {
			case richtext.Click:
				log.Debug().Msgf("BackButton in Chat Clicked: %s", event.ClickData.Kind)
				gui.AppState.CurrentPage = pages.OverviewPage
				gtx.Execute(op.InvalidateCmd{})
			}
		}

		return Text(TextOpts{ThemePtr: gui.MyTheme, TextState: buttonState},
			TextSpan(SpanStyle{
				Font:        gui.FontsNerd[0].Font,
				Size:        30,
				Color:       Gray.NRGBA(),
				Interactive: true,
				Content:     "󰐕", // Plus icon from Nerd font
			}),
		)(gtx)
	}
}
