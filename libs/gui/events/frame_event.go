package gui

import (

	// For debugging output

	gio_app "gioui.org/app"
	"gioui.org/op"
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/ezydark/ezMsg/libs/gui/pages"
	"github.com/ezydark/ezMsg/libs/gui/pages/chat"
	"github.com/ezydark/ezMsg/libs/gui/pages/overview"
)

// HandleFrameEvent remains the same
func HandleFrameEvent(event gio_app.FrameEvent) {
	var ops op.Ops
	gtx := gio_app.NewContext(&ops, event)

	switch gui.AppState.CurrentPage {
	case pages.OverviewPage:
		overview.Overview(gtx)
	case pages.SettingsPage:
		// TODO
	case pages.ChatPage:
		chat.Chat(gtx)
	}

	event.Frame(gtx.Ops)
}
