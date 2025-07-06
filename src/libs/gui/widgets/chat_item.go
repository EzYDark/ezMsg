package widgets

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"gioui.org/widget"
	"gioui.org/x/richtext"
	"github.com/ezydark/ezMsg/src/_test/app/db"
	. "github.com/ezydark/ezMsg/src/ezio"
	"github.com/ezydark/ezMsg/src/libs/gui"
	"github.com/ezydark/ezMsg/src/libs/gui/pages"
	"github.com/rs/zerolog/log"
)

var chat_name richtext.InteractiveText
var chat_msg richtext.InteractiveText
var chat_time richtext.InteractiveText

func ListItemChat(clickable *widget.Clickable) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if clickable.Clicked(gtx) {
			log.Debug().Msg("Switched to chat page")
			gui.AppState.CurrentPage = pages.ChatPage
			gtx.Execute(op.InvalidateCmd{})
		}

		return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if clickable.Hovered() {
				hover_back := color.NRGBA{R: 255, G: 255, B: 255, A: 2}
				paint.FillShape(gtx.Ops, hover_back, clip.Rect{Max: gtx.Constraints.Max}.Op())
			}

			return chatItemContent(gtx)(gtx)
		})
	}
}

func chatItemContent(gtx layout.Context) layout.Widget {
	return Margin(&MarginOpts{Left: 20, Right: 20},
		FlexBox(&FlexBoxOpts{},
			FlexChild(&FlexChildOpts{H: 76},
				FlexBox(&FlexBoxOpts{Axis: Horizontal},
					// Profile picture space
					FlexChild(&FlexChildOpts{W: 54},
						StackBox(&StackBoxOpts{Alignment: Center},
							StackedChild(
								func() layout.Widget {
									var pictureURL string

									for _, member := range gui.AppState.LoggedUser.Chats[0].Members {
										// Find the other member
										if member.Username != gui.AppState.LoggedUser.Username {
											pictureURL = member.ProfilePictureURL
										}
									}

									return Circle(CircleOpts{
										R:      27,
										Color:  LightRed.NRGBA(),
										ImgURL: pictureURL,
									})
								}(),
							),
						),
					),
					// Chat information space
					FlexChild(&FlexChildOpts{Weight: 8.5},
						Margin(&MarginOpts{Left: 16},
							StackBox(&StackBoxOpts{Alignment: W},
								StackedChild(
									FlexBox(&FlexBoxOpts{Axis: Vertical},
										// Chat header (Friend`s name)
										FlexChild(nil,
											Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &chat_name},
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
														Size:    18,
														Color:   White.NRGBA(),
														Content: title,
													})
												}(),
											),
										),
										// Last chat message
										FlexChild(nil,
											FlexBox(&FlexBoxOpts{Axis: Horizontal},
												FlexChild(nil,
													Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &chat_msg},
														TextSpan(SpanStyle{
															Font:    gui.Fonts[0].Font,
															Size:    16,
															Color:   Gray.NRGBA(),
															Content: gui.AppState.LoggedUser.Chats[0].Messages[len(gui.AppState.LoggedUser.Chats[0].Messages)-1].Message,
														}),
													),
												),
												FlexChild(nil,
													Margin(&MarginOpts{Left: 16},
														Text(TextOpts{ThemePtr: gui.MyTheme, TextState: &chat_time},
															TextSpan(SpanStyle{
																Font:    gui.Fonts[0].Font,
																Size:    16,
																Color:   Gray.NRGBA(),
																Content: gui.AppState.LoggedUser.Chats[0].Messages[len(gui.AppState.LoggedUser.Chats[0].Messages)-1].Timestamp.Format("15:04"),
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

func ListItemChat2(loggedUser *db.User) layout.Widget {
	return Margin(&MarginOpts{Left: 20, Right: 20},
		FlexBox(&FlexBoxOpts{},
			FlexChild(&FlexChildOpts{H: 70},
				FlexBox(&FlexBoxOpts{Axis: Horizontal},
					FlexChild(&FlexChildOpts{Weight: 1},
						Rect(RectOpts{Color: Blue.NRGBA()}),
					),
					FlexChild(&FlexChildOpts{Weight: 1},
						Rect(RectOpts{Color: Orange.NRGBA()}),
					),
				),
			),
		),
	)
}
