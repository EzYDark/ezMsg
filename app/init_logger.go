package app

import (
	"fmt"
	"os"

	"github.com/ezydark/ezforce/libs/logger"
	"github.com/fatih/color"
)

// Init custom version of Zerolog logger
func InitLogger() {
	err := logger.Init()
	if err != nil {
		fatal_tag := color.New(color.FgRed, color.Bold).Sprintf("[FATAL]")
		fmt.Println(fatal_tag, "Could not initialize custom logger:\n", err)
		os.Exit(1)
	}
}
