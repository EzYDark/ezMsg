package overview

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/richtext"
	. "github.com/ezydark/ezMsg/ezio"
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/ezydark/ezMsg/libs/gui/widgets"
)

var text_state richtext.InteractiveText
var list_state layout.List
var chat_items_clickable [10]widget.Clickable

func Overview(gtx layout.Context) {
	FlexBox(FlexBoxOpts{Axis: Vertical},
		// App header
		FlexChild(&FlexChildOpts{H: 70},
			Rect(RectOpts{Color: DarkBackground}),
			Margin(&MarginOpts{All: 20},
				FlexBox(FlexBoxOpts{Axis: Horizontal},
					FlexChild(&FlexChildOpts{Weight: 1},
						Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &text_state},
							TextSpan(SpanStyle{
								Content: "Ez",
								Font:    gui.Fonts[1].Font,
								Color:   White,
								Size:    unit.Sp(36),
							}),
							TextSpan(SpanStyle{
								Content: "Msg",
								Font:    gui.Fonts[1].Font,
								Color:   White,
								Size:    unit.Sp(36),
							}),
						),
					),
					FlexChild(&FlexChildOpts{Weight: 1},
						DirectionBox(&DirectionBoxOpts{Direction: SE},
							Rect(RectOpts{Color: LightGreen, W: 24, H: 24}),
						),
					),
				),
			),
		),
		// App main content
		FlexChild(&FlexChildOpts{Weight: 1},
			Rect(RectOpts{Color: DarkBackground}),
			Margin(&MarginOpts{Top: 20, Bottom: 20},
				ListBox(ListOpts{ListPtr: &list_state, Axis: Vertical},
					ListChild(
						widgets.ListItemChat(&chat_items_clickable[0]),
					),
					ListChild(
						widgets.ListItemChat2(gui.DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						widgets.ListItemChat(&chat_items_clickable[1]),
					),
					ListChild(
						widgets.ListItemChat2(gui.DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						widgets.ListItemChat(&chat_items_clickable[2]),
					),
					ListChild(
						widgets.ListItemChat2(gui.DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						widgets.ListItemChat(&chat_items_clickable[3]),
					),
					ListChild(
						widgets.ListItemChat2(gui.DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						widgets.ListItemChat(&chat_items_clickable[4]),
					),
					ListChild(
						widgets.ListItemChat2(gui.DBPtr.RegisteredUsers[0]),
					),
					ListChild(
						widgets.ListItemChat(&chat_items_clickable[5]),
					),
					ListChild(
						widgets.ListItemChat2(gui.DBPtr.RegisteredUsers[0]),
					),
				),
			),
		),
		// // App footer
		// FlexChild(&FlexChildOpts{H: 70},
		// 	Rect(RectOpts{Color: LightRed}),
		// ),
	)(gtx)
}
