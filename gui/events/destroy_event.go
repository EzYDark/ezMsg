package gui

import (
	"os"

	gio_app "gioui.org/app"
	"github.com/ezydark/ezlog/log"
)

// HandleDestroyEvent processes window close events
// Separated to improve error handling and make cleanup more explicit
func HandleDestroyEvent(event gio_app.DestroyEvent) error {
	if event.Err != nil {
		log.Error().Msgf("Error closing application:\n%v", event.Err)
		return event.Err
	}
	log.Info().Msg("Application closed.")
	os.Exit(0)
	return nil
}
