package widgets

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/x/richtext"
	"github.com/ezydark/ezMsg/_test/app/db"
	"github.com/ezydark/ezMsg/gui"
	. "github.com/ezydark/ezio"
)

func ChatMsgBubble(message *db.Message, index int) layout.Widget {
	isLoggedUser := message.Sender.Username == gui.AppState.LoggedUser.Username

	var messageTextState richtext.InteractiveText
	var senderTextStae richtext.InteractiveText

	var msgAlignment layout.Alignment
	var msgDirection layout.Direction

	var bubbleColor color.NRGBA

	if isLoggedUser {
		// Current user's message bubble
		msgAlignment = layout.End
		msgDirection = layout.E
		bubbleColor = color.NRGBA{R: 0x00, G: 0x7A, B: 0xFF, A: 0xFF}
	} else {
		// Other user's message bubble
		msgAlignment = layout.Start
		msgDirection = layout.W
		bubbleColor = color.NRGBA{R: 0x4A, G: 0x4A, B: 0x4A, A: 0xFF}
	}

	return func(gtx layout.Context) layout.Dimensions {
		return Margin(&MarginOpts{Bottom: 8},
			FlexBox(&FlexBoxOpts{Axis: Vertical, Alignment: msgAlignment},
				FlexChild(nil,
					FlexBox(&FlexBoxOpts{Axis: Vertical, Alignment: End},
						// Sender of the message
						FlexChild(nil,
							Margin(&MarginOpts{Bottom: 4},
								Align(&AlignOpts{Alignment: msgDirection},
									Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &senderTextStae},
										TextSpan(SpanStyle{
											Content: message.Sender.Username,
											Font:    gui.Fonts[1].Font,
											Size:    unit.Sp(14),
											Color:   LightGray.NRGBA(),
										}),
									),
								),
							),
						),
						// The text message
						FlexChild(nil,
							BackgroundStackBox(&StackBoxOpts{Alignment: msgDirection},
								Rect(RectOpts{Color: bubbleColor}),
								Margin(&MarginOpts{Top: 6, Bottom: 6, Left: 12, Right: 12},
									func() layout.Widget {
										if len(message.Files) > 0 {
											return FlexBox(&FlexBoxOpts{Axis: Vertical},
												FlexChild(nil,
													Margin(&MarginOpts{Bottom: 6},
														Rect(RectOpts{
															Color:  Gray.NRGBA(),
															W:      200,
															H:      200,
															ImgURL: message.Files[0].URL,
														}),
													),
												),
												FlexChild(nil,
													Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &messageTextState},
														TextSpan(SpanStyle{
															Content: message.Message,
															Font:    gui.Fonts[0].Font,
															Size:    unit.Sp(16),
															Color:   White.NRGBA(),
														}),
													),
												),
											)
										} else {
											return Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &messageTextState},
												TextSpan(SpanStyle{
													Content: message.Message,
													Font:    gui.Fonts[0].Font,
													Size:    unit.Sp(16),
													Color:   White.NRGBA(),
												}),
											)
										}
									}(),
								),
							),
						),
					),
				),
			),
		)(gtx)
	}
}
