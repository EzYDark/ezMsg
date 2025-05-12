package gui

import (
	gio_app "gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/x/component"
	. "github.com/ezydark/ezMsg/ezio"
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/ezydark/ezMsg/libs/gui/pages"
	"github.com/rs/zerolog/log"
)

// HandleFrameEvent coordinates frame rendering and input processing
// Splitting this out provides a clear separation between event handling and rendering
func HandleFrameEvent(event gio_app.FrameEvent) {
	// Operations stack stores drawing operations between frames
	// Gio uses a retained mode API where UI operations are recorded then executed
	var ops op.Ops

	// Create a layout context for this frame with fresh operations list
	// Context carries layout constraints and input state
	gtx := gio_app.NewContext(&ops, event)

	// TestEzMsg(gtx)
	pages.Overview(gtx)

	// Submit operations to be rendered
	// This is where the actual drawing instructions are executed
	event.Frame(gtx.Ops)
}

var myBtn widget.Clickable
var myInput widget.Editor
var myContextArea = component.ContextArea{Activation: pointer.ButtonSecondary}

// var myList layout.List

func TestEzMsg(gtx layout.Context) {
	for myBtn.Clicked(gtx) {
		log.Debug().Msg("Clicked button!")
	}

	Flex(FlexOpts{Alignment: Vertical},
		FlexChild(
			Rect(RectOpts{Color: LightRed}),
			Flex(FlexOpts{Alignment: Vertical},
				FlexChildStatic(FlexChildStaticOpts{W: 150, H: 80},
					Rect(RectOpts{Color: LightBlue}),
					ContextMenu(ContextMenuOpts{ContextArea: &myContextArea},
						Rect(RectOpts{W: 100, H: 100, Color: LightPurple}),
					),
				),
				FlexChildStatic(FlexChildStaticOpts{W: 200, H: 75},
					Rect(RectOpts{Color: LightOrange}),
				),
			),
		),
		FlexChildStatic(FlexChildStaticOpts{H: 50},
			Flex(FlexOpts{Alignment: Horizontal},
				FlexChild(
					Rect(RectOpts{Color: LightPink}),
					Margin(MarginOpts{All: gtx.Metric.PxToDp(8)},
						Input(InputOpts{
							EditorPtr: &myInput,
							ThemePtr:  gui.MyTheme,
							Hint:      "Type new message...",
						}),
					),
				),
				FlexChildStatic(FlexChildStaticOpts{W: 100},
					Button(ButtonOpts{
						ButtonPtr: &myBtn,
						Text:      "Click Me!",
						ThemePtr:  gui.MyTheme,
					}),
				),
			),
		),
	)(gtx)

	// ContextMenu(ContextMenuOpts{ContextArea: &myContextArea},
	// 	Rect(RectOpts{W: 100, H: 100, Color: LightPurple}),
	// )(gtx)
}
