package components

import (
	"gioui.org/layout"
	"gioui.org/x/richtext"
	"github.com/ezydark/ezMsg/app/db"
	. "github.com/ezydark/ezMsg/ezio"
	"github.com/ezydark/ezMsg/libs/gui"
)

var chat_name richtext.InteractiveText
var chat_msg richtext.InteractiveText
var chat_time richtext.InteractiveText

func ListItemChat(loggedUser db.User) layout.Widget {
	return FlexBox(FlexBoxOpts{},
		FlexChild(&FlexChildOpts{H: 76},
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
														Color:   Gray,
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
															Color:   Gray,
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
	)
}

func ListItemChat2(loggedUser db.User) layout.Widget {
	return FlexBox(FlexBoxOpts{},
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
	)
}
