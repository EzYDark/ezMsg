package chat

import (
	"gioui.org/layout"
	"gioui.org/x/richtext"
	. "github.com/ezydark/ezMsg/ezio"
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/ezydark/ezMsg/libs/gui/widgets"
)

var chatTitleState richtext.InteractiveText
var backButtonState richtext.InteractiveText

func Chat(gtx layout.Context) {
	BackgroundBox(BackgroundBoxOpts{},
		Rect(RectOpts{Color: DarkBackground}),
		Margin(&MarginOpts{All: 20},
			FlexBox(FlexBoxOpts{Axis: Vertical},
				// Header
				FlexChild(&FlexChildOpts{H: 70},
					FlexBox(FlexBoxOpts{Axis: Horizontal},
						// Back button to Overview
						FlexChild(&FlexChildOpts{Weight: 1},
							DirectionBox(&DirectionBoxOpts{Direction: W},
								widgets.BackButton(&backButtonState),
							),
						),
						// Chat title
						FlexChild(&FlexChildOpts{Weight: 1},
							DirectionBox(&DirectionBoxOpts{Direction: Center},
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
											Color:   White,
											Content: title,
										})
									}(),
								),
							),
						),
						// Chat options button
						FlexChild(&FlexChildOpts{Weight: 1}),
					),
				),
				// Chat
				FlexChild(&FlexChildOpts{Weight: 1},
					Rect(RectOpts{Color: LightPurple}),
				),
				// Footer (Input box)
				FlexChild(&FlexChildOpts{H: 70},
					Rect(RectOpts{Color: LightOrange}),
				),
			),
		),
	)(gtx)
}
