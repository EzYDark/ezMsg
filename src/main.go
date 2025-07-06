package main

import (
	"os"

	gio_app "gioui.org/app"
	"gioui.org/unit"
	"github.com/ezydark/ezMsg/src/libs"
	"github.com/ezydark/ezMsg/src/libs/client"
	gui "github.com/ezydark/ezMsg/src/libs/gui/events"
	"github.com/rs/zerolog/log"
)

const (
	WindowWidth  = 450
	WindowHeight = 900
	WindowTitle  = "ezChat"
)

func main() {
	libs.InitLogger()
	log.Info().Msgf("Starting %s...", WindowTitle)

	err := client.InitClient()
	if err != nil {
		log.Fatal().Msgf("Error running test client:\n%v", err)
	}
	return

	// Start UI in a separate goroutine because gio_app.Main() blocks
	// This allows clean separation between initialization and the event loop
	go func() {
		// Create a new window without showing it yet
		window := new(gio_app.Window)

		window.Option(gio_app.Title(WindowTitle))
		window.Option(gio_app.Size(unit.Dp(WindowWidth), unit.Dp(WindowHeight)))

		err := run_gui(window)
		if err != nil {
			log.Fatal().Msgf("Error running %s:\n%v", WindowTitle, err)
		}
		os.Exit(0)
	}()

	// gio_app.Main() handles platform-specific event loop initialization
	// It blocks until all windows are closed and must be the last call in main
	gio_app.Main()
}

// run_gui initializes and runs the main GUI event loop
// This is the core of a Gio application where events are processed
func run_gui(window *gio_app.Window) error {
	// Event loop continuously processes window events
	for {
		// Window.Event() blocks until next event occurs
		event := window.Event()
		switch e := event.(type) {
		case gio_app.DestroyEvent:
			// Handle window closing events
			return gui.HandleDestroyEvent(e)
		case gio_app.FrameEvent:
			// Process frame update events - this is where UI rendering happens
			gui.HandleFrameEvent(e)
		}
	}
}
