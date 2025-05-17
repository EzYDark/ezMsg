package chat

import (
	"gioui.org/layout"
	. "github.com/ezydark/ezMsg/ezio"
)

// var DBPtr = db.InitDB();

func Chat(gtx layout.Context) {
	FlexBox(FlexBoxOpts{Axis: Vertical},
		FlexChild(&FlexChildOpts{H: 70},
			Rect(RectOpts{Color: LightPink}),
		),
		FlexChild(nil,
			Rect(RectOpts{Color: LightPurple}),
		),
	)(gtx)
}
