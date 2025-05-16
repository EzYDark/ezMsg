package overview

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/richtext"
	"github.com/ezydark/ezMsg/app/db"
	. "github.com/ezydark/ezMsg/ezio"
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/ezydark/ezMsg/libs/gui/components"
)

var text_state richtext.InteractiveText
var list_state layout.List
var chat_items_clickable [10]widget.Clickable

var DBPtr = db.InitDB()

func Overview(gtx layout.Context) {
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
						components.ListItemChat(DBPtr.RegisteredUsers[0], &chat_items_clickable[0]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0], &chat_items_clickable[1]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0], &chat_items_clickable[2]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0], &chat_items_clickable[3]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0], &chat_items_clickable[4]),
					),
					ListChild(
						components.ListItemChat2(DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						components.ListItemChat(DBPtr.RegisteredUsers[0], &chat_items_clickable[5]),
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
