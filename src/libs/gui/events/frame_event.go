package gui

import (
	gio_app "gioui.org/app"
	"gioui.org/op"
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/ezydark/ezMsg/libs/gui/pages"
	"github.com/ezydark/ezMsg/libs/gui/pages/chat"
	"github.com/ezydark/ezMsg/libs/gui/pages/overview"
)

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
