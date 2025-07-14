package main

import (
	"flag"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	gui "github.com/ezydark/ezMsg/gui/events"
	"github.com/ezydark/ezdebug/tui"
	"github.com/ezydark/ezlog"
	"github.com/ezydark/ezlog/log"
)

const (
	WindowWidth  = 450
	WindowHeight = 750
	WindowTitle  = "ezChat"
)

func main() {
	serverFlag := flag.Bool("server", false, "Run in server mode")
	clientFlag := flag.Bool("client", false, "Run in client mode")
	debugFlag := flag.Bool("debug", false, "Run in debug TUI mode")
	flag.Parse()

	if *debugFlag {
		debugTUI := tui.Init()

		ezlog.New().WithWriter(debugTUI.GetLogWriter()).WithTviewCompat().Build()

		log.Info().Msgf("Starting %s in debug TUI mode...", WindowTitle)

		features := []tui.Feature{
			{
				Name:           "Run Server",
				Description:    "Run Server for other clients to connect",
				StartOnStartup: *serverFlag,
				OnStart: func(self *tui.Feature) {
					self.Enabled = true
					log.Info().Msg("Starting Server...")
					// server.InitServer()
				},
				OnStop: func(self *tui.Feature) {
					self.Enabled = false
					log.Info().Msg("Stopping Server...")
					// server.CloseServer()
				},
			},
			{
				Name:           "Run Client",
				Description:    "Run Client to connect to the already running Server",
				StartOnStartup: *clientFlag,
				OnStart: func(self *tui.Feature) {
					self.Enabled = true
					log.Info().Msg("Starting Client...")
					// err := client.InitClient()
					// if err != nil {
					// 	log.Fatal().Msgf("Error running test client:\n%v", err)
					// }
				},
				OnStop: func(self *tui.Feature) {
					self.Enabled = false
					log.Info().Msg("Stopping Client...")
					// client.CloseClient()
				},
			},
		}
		tui.GetFeatureList().Set(features)

		// Start DebugTUI application (Blocking call)
		if err := debugTUI.Start(); err != nil {
			log.Fatal().Msgf("Error running ezDebugTUI example:\n%v", err)
		}
	} else {
		ezlog.New().Build()
		log.Info().Msgf("Starting %s...", WindowTitle)

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
