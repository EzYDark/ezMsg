package chat

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/x/richtext"
	. "github.com/ezydark/ezMsg/ezio"
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/ezydark/ezMsg/libs/gui/widgets"
)

var chatTitleState richtext.InteractiveText
var backButtonState richtext.InteractiveText
var chatListState layout.List

var inputBoxState widget.Editor

func Chat(gtx layout.Context) {
	activeChat := gui.AppState.LoggedUser.Chats[0]

	var msgWidgets []layout.ListElement
	for i, msg := range activeChat.Messages {
		msgWidgets = append(msgWidgets, ListChild(widgets.ChatMsgBubble(msg, i)))
	}

	BackgroundBox(
		Rect(RectOpts{Color: DarkBackground.NRGBA()}),
		Margin(&MarginOpts{All: 20},
			FlexBox(&FlexBoxOpts{Axis: Vertical},
				// Header
				FlexChild(&FlexChildOpts{H: 70},
					// Rect(RectOpts{Color: LightRed.NRGBA()}),
					FlexBox(&FlexBoxOpts{Axis: Horizontal},
						// Back button to Overview
						FlexChild(&FlexChildOpts{Weight: 1},
							widgets.BackButton(&backButtonState),
						),
						// Chat title
						FlexChild(&FlexChildOpts{Weight: 1},
							Align(&StackBoxOpts{Alignment: N},
								Margin(&MarginOpts{Top: 5}, // TEMP fix of bad alignment with the Back button
									Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &chatTitleState},
										func() richtext.SpanStyle {
											var title string

											for _, member := range gui.AppState.LoggedUser.Chats[0].Members {
												// Find the other member
												if member.Username != gui.AppState.LoggedUser.Username {
													title = member.Username
												}
											}

											return TextSpan(SpanStyle{
												Font:    gui.Fonts[1].Font,
												Size:    20,
												Color:   White.NRGBA(),
												Content: title,
											})
										}(),
									),
								),
							),
						),
						// Chat options button
						FlexChild(&FlexChildOpts{Weight: 1}),
					),
				),
				// Chat
				FlexChild(&FlexChildOpts{Weight: 1},
					ListBox(ListOpts{
						ListPtr:     &chatListState,
						Axis:        Vertical,
						ScrollToEnd: true,
					},
						msgWidgets...,
					),
				),
				// Footer (Input box)
				FlexChild(&FlexChildOpts{H: 70},
					BackgroundBox(
						Rect(RectOpts{Color: color.NRGBA{R: uint8(193), G: uint8(193), B: uint8(193), A: uint8(5)}}),
						FlexBox(&FlexBoxOpts{Axis: Horizontal, Spacing: SpaceBetween},
							// Message input box
							FlexChild(nil,
								Margin(&MarginOpts{All: 6},
									Input(InputOpts{
										EditorPtr: &inputBoxState,
										ThemePtr:  gui.MyTheme,
										Hint:      "Enter your message here...",
									}),
								),
							),
							// Add button for attachments
							FlexChild(nil,
								// TODO: Add hover effect on the Add button
								Margin(&MarginOpts{Top: 1, Right: 8},
									widgets.AddButton(&backButtonState),
								),
							),
						),
					),
				),
			),
		),
	)(gtx)
}
