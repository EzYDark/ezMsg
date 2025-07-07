package libs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var initialized bool

// Init initializes the logger with the default stdout writer.
func Init() error {
	return InitWithWriter(os.Stdout)
}

// InitWithWriter initializes the logger with a custom writer.
func InitWithWriter(writer io.Writer) error {
	if initialized {
		return errors.New("logger is already initialized")
	}

	// Configure global settings
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = "15:04:05.000"

	// Configure custom console writer
	consoleOutput := zerolog.ConsoleWriter{
		Out:        writer, // Use the provided writer
		TimeFormat: "15:04:05.000",
		NoColor:    false,
	}

	consoleOutput.FormatLevel = func(i any) string {
		levelStr := strings.ToUpper(fmt.Sprintf("%s", i))

		switch levelStr {
		case "DEBUG":
			return color.New(color.FgBlue).Sprintf("{%s}", levelStr)
		case "INFO":
			return color.New(color.FgGreen).Sprintf("{%s}", levelStr)
		case "WARN":
			return color.New(color.FgYellow).Sprintf("{%s}", levelStr)
		case "ERROR":
			return color.New(color.FgRed).Sprintf("{%s}", levelStr)
		case "FATAL":
			return color.New(color.FgRed, color.Bold).Sprintf("{%s}", levelStr)
		default:
			return color.New(color.FgWhite).Sprintf("{%s}", levelStr)
		}
	}

	consoleOutput.FormatMessage = func(i any) string {
		return fmt.Sprintf("%s", i)
	}

	consoleOutput.FormatFieldName = func(i any) string {
		return fmt.Sprintf("%s=", i)
	}

	consoleOutput.FormatFieldValue = func(i any) string {
		return fmt.Sprintf("%s", i)
	}

	// Create a new logger instance
	newLogger := zerolog.New(consoleOutput).With().Timestamp().Logger()

	// Set the global logger
	log.Logger = newLogger

	initialized = true
	return nil
}
