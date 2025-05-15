package overview

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/x/richtext"
	"github.com/ezydark/ezMsg/app/db"
	. "github.com/ezydark/ezMsg/ezio"
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/ezydark/ezMsg/libs/gui/components"
)

var text_state richtext.InteractiveText
var list_state layout.List

var DBPtr = db.InitDB()

func Overview(gtx layout.Context) {
	DarkBackground := color.NRGBA{R: uint8(37), G: uint8(35), B: uint8(49), A: uint8(255)}

	FlexBox(FlexBoxOpts{Axis: Vertical},
		// App header
		FlexChild(&FlexChildOpts{H: 70},
			Rect(RectOpts{Color: DarkBackground}),
			Margin(&MarginOpts{All: 20},
				FlexBox(FlexBoxOpts{Axis: Horizontal, Spacing: SpaceBetween, Alignment: Middle},
					FlexChild(nil,
						DirectionBox(&DirectionBoxOpts{Direction: W},
							Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &text_state},
								TextSpan(SpanStyle{
									Content:     "Ez",
									Font:        gui.Fonts[1].Font,
									Color:       White,
									Interactive: true,
									Size:        unit.Sp(36),
								}),
								TextSpan(SpanStyle{
									Content:     "Msg",
									Font:        gui.Fonts[1].Font,
									Color:       White,
									Interactive: true,
									Size:        unit.Sp(36),
								}),
							),
						),
					),
					FlexChild(nil,
						Rect(RectOpts{Color: LightGreen, W: 12, H: 12}),
					),
				),
			),
		),
		// App main content
		FlexChild(&FlexChildOpts{Weight: 1},
			Rect(RectOpts{Color: DarkBackground}),
			Margin(&MarginOpts{All: 20},
				ListBox(ListOpts{ListPtr: &list_state, Axis: Vertical},
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
				),
			),
		),
		// App footer
		FlexChild(&FlexChildOpts{H: 70},
			Rect(RectOpts{Color: LightRed}),
		),
	)(gtx)
}
