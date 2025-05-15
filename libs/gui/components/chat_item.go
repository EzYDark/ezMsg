package components

import (
	"gioui.org/layout"
	"github.com/ezydark/ezMsg/app/db"
	. "github.com/ezydark/ezMsg/ezio"
)

func ListItemChat(loggedUser db.User) layout.Widget {
	return FlexBox(FlexBoxOpts{},
		FlexChild(&FlexChildOpts{H: 70},
			FlexBox(FlexBoxOpts{Axis: Horizontal},
				// Profile picture space
				FlexChild(&FlexChildOpts{W: 100},
					Rect(RectOpts{Color: Red}),
					StackBox(StackOpts{Alignment: Center},
						StackedChild(Circle(CircleOpts{R: 25, Color: LightRed})),
					),
				),
				// Chat information space
				FlexChild(&FlexChildOpts{Weight: 8.5},
					Rect(RectOpts{Color: Green}),
					FlexBox(FlexBoxOpts{Axis: Vertical},
						// Chat header (Friend`s name)
						FlexChild(nil,
							Rect(RectOpts{Color: Yellow}),
						),
						// Last chat message
						FlexChild(nil,
							Rect(RectOpts{Color: Purple}),
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
				FlexChild(nil,
					Rect(RectOpts{Color: Blue}),
				),
				FlexChild(nil,
					Rect(RectOpts{Color: Orange}),
				),
			),
		),
	)
}
