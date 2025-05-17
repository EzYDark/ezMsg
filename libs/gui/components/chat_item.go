package components

import (
	// Added for defining custom colors if needed

	// Added for click logging

	"image/color"

	"gioui.org/layout" // Added for clipping the background shape
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	// Added for painting the background
	"gioui.org/widget" // Added for widget.Clickable
	"gioui.org/x/richtext"
	"github.com/ezydark/ezMsg/app/db"
	. "github.com/ezydark/ezMsg/ezio" // Assuming ezio re-exports colors like White, Gray, etc.
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/ezydark/ezMsg/libs/gui/pages"
	"github.com/rs/zerolog/log"
)

var chat_name richtext.InteractiveText
var chat_msg richtext.InteractiveText
var chat_time richtext.InteractiveText

// ListItemChat now accepts a user and a clickable to manage its state.
func ListItemChat(loggedUser db.User, clickable *widget.Clickable) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if clickable.Clicked(gtx) {
			log.Debug().Msg("Switched to chat page")
			pages.AppState.CurrentPage = pages.ChatPage
			gtx.Execute(op.InvalidateCmd{})
		}

		return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if clickable.Hovered() {
				hover_back := color.NRGBA{R: 255, G: 255, B: 255, A: 2}
				paint.FillShape(gtx.Ops, hover_back, clip.Rect{Max: gtx.Constraints.Max}.Op())
			}

			return chatItemContent(gtx, loggedUser)(gtx)
		})
	}
}

// chatItemContent contains the original drawing logic for the chat item.
func chatItemContent(gtx layout.Context, loggedUser db.User) layout.Widget {
	return Margin(&MarginOpts{Left: 20, Right: 20},
		FlexBox(FlexBoxOpts{},
			FlexChild(&FlexChildOpts{H: 76}, // Ensure this height matches clickable area or vice-versa
				FlexBox(FlexBoxOpts{Axis: Horizontal},
					// Profile picture space
					FlexChild(&FlexChildOpts{W: 54},
						StackBox(StackBoxOpts{Alignment: Center},
							StackedChild(
								Circle(CircleOpts{R: 27, Color: LightRed, ImgURL: loggedUser.ProfilePictureURL}),
							),
						),
					),
					// Chat information space
					FlexChild(&FlexChildOpts{Weight: 8.5},
						Margin(&MarginOpts{Left: 16},
							StackBox(StackBoxOpts{Alignment: W},
								StackedChild(
									FlexBox(FlexBoxOpts{Axis: Vertical},
										// Chat header (Friend`s name)
										FlexChild(nil,
											Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &chat_name},
												TextSpan(SpanStyle{
													Font:    gui.Fonts[1].Font,
													Size:    18,
													Color:   White,
													Content: loggedUser.Chats[0].Members[0].Username,
												}),
											),
										),
										// Last chat message
										FlexChild(nil,
											FlexBox(FlexBoxOpts{Axis: Horizontal},
												FlexChild(nil,
													Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &chat_msg},
														TextSpan(SpanStyle{
															Font:    gui.Fonts[0].Font,
															Size:    16,
															Color:   Gray, // ezio.Gray
															Content: loggedUser.Chats[0].Messages[len(loggedUser.Chats[0].Messages)-1].Message,
														}),
													),
												),
												FlexChild(nil,
													Margin(&MarginOpts{Left: 16},
														Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &chat_time},
															TextSpan(SpanStyle{
																Font:    gui.Fonts[0].Font,
																Size:    16,
																Color:   Gray, // ezio.Gray
																Content: loggedUser.Chats[0].Messages[len(loggedUser.Chats[0].Messages)-1].Timestamp.Format("15:04"),
															}),
														),
													),
												),
											),
										),
									),
								),
							),
						),
					),
				),
			),
		),
	)
}

// ListItemChat2 remains unchanged for now, but you can apply similar logic if needed.
func ListItemChat2(loggedUser db.User) layout.Widget {
	return Margin(&MarginOpts{Left: 20, Right: 20},
		FlexBox(FlexBoxOpts{},
			FlexChild(&FlexChildOpts{H: 70},
				FlexBox(FlexBoxOpts{Axis: Horizontal},
					FlexChild(&FlexChildOpts{Weight: 1},
						Rect(RectOpts{Color: Blue}),
					),
					FlexChild(&FlexChildOpts{Weight: 1},
						Rect(RectOpts{Color: Orange}),
					),
				),
			),
		),
	)
}
