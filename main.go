package main

import (
	"flag"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ezydark/ezMsg/src/libs"
	"github.com/ezydark/ezMsg/src/libs/client"
	tui "github.com/ezydark/ezMsg/src/libs/debug_tui"
	gui "github.com/ezydark/ezMsg/src/libs/gui/events"
	server "github.com/ezydark/ezMsg/src/libs/server/comm"
	"github.com/rs/zerolog/log"
)

const (
	WindowWidth  = 450
	WindowHeight = 900
	WindowTitle  = "ezChat"
)

func main() {
	serverFlag := flag.Bool("server", false, "Run in server mode")
	clientFlag := flag.Bool("client", false, "Run in client mode")
	debugFlag := flag.Bool("debug", false, "Run in debug TUI mode")
	flag.Parse()

	if *debugFlag {
		app := tui.NewApp()
		libs.InitWithWriter(app.LogWriter())
		log.Info().Msgf("Starting %s in debug TUI mode...", WindowTitle)
		if err := app.Run(); err != nil {
			log.Fatal().Msgf("Error running debug TUI:\n%v", err)
		}
	}

	libs.InitWithWriter(os.Stderr)
	log.Info().Msgf("Starting %s...", WindowTitle)

	if *serverFlag {
		server.InitServer()
	}

	if *clientFlag {
		err := client.InitClient()
		if err != nil {
			log.Fatal().Msgf("Error running test client:\n%v", err)
		}
	}

	if !*serverFlag && !*clientFlag {
		// Start UI in a separate goroutine because app.Main() blocks
		// This allows clean separation between initialization and the event loop
		go func() {
			// Create a new window without showing it yet
			window := new(app.Window)

			window.Option(app.Title(WindowTitle))
			window.Option(app.Size(unit.Dp(WindowWidth), unit.Dp(WindowHeight)))

			err := run_gui(window)
			if err != nil {
				log.Fatal().Msgf("Error running %s:\n%v", WindowTitle, err)
			}
			os.Exit(0)
		}()

		// app.Main() handles platform-specific event loop initialization
		// It blocks until all windows are closed and must be the last call in main
		app.Main()
	}
}

// run_gui initializes and runs the main GUI event loop
// This is the core of a Gio application where events are processed
func run_gui(window *app.Window) error {
	// Event loop continuously processes window events
	for {
		// Window.Event() blocks until next event occurs
		event := window.Event()
		switch e := event.(type) {
		case app.DestroyEvent:
			// Handle window closing events
			return gui.HandleDestroyEvent(e)
		case app.FrameEvent:
			// Process frame update events - this is where UI rendering happens
			gui.HandleFrameEvent(e)
		}
	}
}
