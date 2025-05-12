package overview

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/x/richtext"
	. "github.com/ezydark/ezMsg/ezio"
	"github.com/ezydark/ezMsg/libs/gui"
)

var text_state richtext.InteractiveText

func Overview(gtx layout.Context) {
	Flex(FlexOpts{Alignment: Vertical},
		FlexChildWeight(1,
			Rect(RectOpts{Color: LightBlue}),
			Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &text_state},
				TextSpan(SpanStyle{
					Content:     "Testing",
					Font:        gui.Fonts[0].Font,
					Color:       Green,
					Interactive: true,
					Size:        unit.Sp(26),
				}, ResponsiveOpts{}),
				TextSpan(SpanStyle{
					Content:     "perfectly",
					Font:        gui.Fonts[0].Font,
					Color:       Blue,
					Interactive: true,
					Size:        unit.Sp(40),
				}, ResponsiveOpts{}),
				TextSpan(SpanStyle{
					Content:     "perfectly",
					Font:        gui.Fonts[0].Font,
					Color:       Red,
					Interactive: true,
					Size:        unit.Sp(30),
				}, ResponsiveOpts{}),
			),
		),
		FlexChildWeight(1,
			Rect(RectOpts{Color: LightGreen}),
		),
		FlexChildWeight(7,
			Rect(RectOpts{Color: LightOrange}),
		),
		FlexChildWeight(1,
			Rect(RectOpts{Color: LightRed}),
		),
	)(gtx)
}
