package widgets

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/x/richtext"
	"github.com/ezydark/ezMsg/app/db"
	. "github.com/ezydark/ezMsg/ezio"
	"github.com/ezydark/ezMsg/libs/gui"
)

// func ChatMessageOld(message db.Message) layout.Widget {
// 	var localMessageTextState richtext.InteractiveText
// 	var localSenderTextState richtext.InteractiveText

// 	return func(gtx layout.Context) layout.Dimensions {
// 		isCurrentUser := message.User.Username == gui.AppState.LoggedUser.Username
// 		// messageAlignment := W
// 		bubbleColor := color.NRGBA{R: 0x4A, G: 0x4A, B: 0x4A, A: 0xFF} // Other user's message bubble
// 		textColor := White

// 		if isCurrentUser {
// 			// messageAlignment = E
// 			bubbleColor = color.NRGBA{R: 0x00, G: 0x7A, B: 0xFF, A: 0xFF} // Current user's message bubble
// 		}

// 		messageContent := FlexBox(FlexBoxOpts{Axis: Vertical},
// 			FlexChild(nil,
// 				Margin(&MarginOpts{Bottom: unit.Dp(2), Left: unit.Dp(8), Right: unit.Dp(8), Top: unit.Dp(5)},
// 					Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &localSenderTextState},
// 						TextSpan(SpanStyle{
// 							Content: message.User.Username,
// 							Font:    gui.Fonts[1].Font,
// 							Size:    unit.Sp(14),
// 							Color:   LightGray.NRGBA(),
// 						}),
// 					),
// 				),
// 			),
// 			FlexChild(nil,
// 				Margin(&MarginOpts{Bottom: unit.Dp(5), Left: unit.Dp(8), Right: unit.Dp(8)},
// 					Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &localMessageTextState},
// 						TextSpan(SpanStyle{
// 							Content: message.Message,
// 							Font:    gui.Fonts[0].Font,
// 							Size:    unit.Sp(16),
// 							Color:   textColor.NRGBA(),
// 						}),
// 					),
// 				),
// 			),
// 		)

// 		return FlexBox(FlexBoxOpts{Axis: Horizontal},
// 			FlexChild(&FlexChildOpts{Weight: 1},
// 				FlexBox(FlexBoxOpts{Axis: Vertical},
// 					FlexChild(&FlexChildOpts{MaxW: 50},
// 						Rect(RectOpts{Color: bubbleColor}),
// 						messageContent,
// 					),
// 				),
// 			),
// 		)(gtx)
// 	}
// }

func ChatMessage(message db.Message, i int) layout.Widget {
	var localMessageTextState richtext.InteractiveText
	var localSenderTextState richtext.InteractiveText

	return func(gtx layout.Context) layout.Dimensions {
		isCurrentUser := message.User.Username == gui.AppState.LoggedUser.Username
		var isRight = layout.W
		bubbleColor := color.NRGBA{R: 0x4A, G: 0x4A, B: 0x4A, A: 0xFF} // Other user's message bubble
		textColor := White

		if isCurrentUser {
			isRight = layout.E
			bubbleColor = color.NRGBA{R: 0x00, G: 0x7A, B: 0xFF, A: 0xFF} // Current user's message bubble
		}

		messageContent := FlexBox(FlexBoxOpts{Axis: Vertical},
			FlexChild(nil,
				Margin(&MarginOpts{Bottom: unit.Dp(2), Left: unit.Dp(8), Right: unit.Dp(8), Top: unit.Dp(5)},
					Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &localSenderTextState},
						TextSpan(SpanStyle{
							Content: message.User.Username,
							Font:    gui.Fonts[1].Font,
							Size:    unit.Sp(14),
							Color:   LightGray.NRGBA(),
						}),
					),
				),
			),
			FlexChild(nil,
				Margin(&MarginOpts{Bottom: unit.Dp(5), Left: unit.Dp(8), Right: unit.Dp(8)},
					Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &localMessageTextState},
						TextSpan(SpanStyle{
							Content: message.Message,
							Font:    gui.Fonts[0].Font,
							Size:    unit.Sp(16),
							Color:   textColor.NRGBA(),
						}),
					),
				),
			),
		)

		return FlexBox(FlexBoxOpts{Axis: Horizontal},
			FlexChild(&FlexChildOpts{H: 50},
				FlexBox(FlexBoxOpts{Axis: Horizontal},
					FlexChild(&FlexChildOpts{Weight: 1},
						// Rect(RectOpts{Color: colors.GetLightNRGBA(i)}),
						DirectionBox(&DirectionBoxOpts{Direction: isRight},
							FlexBox(FlexBoxOpts{Axis: Horizontal},
								FlexChild(nil,
									BackgroundBox(BackgroundBoxOpts{},
										Rect(RectOpts{Color: bubbleColor}),
										messageContent,
									),
								),
							),
						),
					),
				),
			),
		)(gtx)
	}
}
