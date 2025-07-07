package tui

import (
	"fmt"
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App is the main TUI application.
type App struct {
	*tview.Application
	LogView *tview.TextView
}

// Feature defines a toggleable feature with its state and action.
type Feature struct {
	Name        string
	Description string
	Enabled     bool
	Action      func(enabled bool)
}

// featureState holds the initial configuration of all toggleable features.
var featureState = []Feature{
	{
		Name:        "Enable Verbose Logging",
		Description: "When enabled, logs will include detailed debug information.",
		Enabled:     false,
	},
	{
		Name:        "Mock Server Responses",
		Description: "Simulates server responses without a live connection. Good for UI testing.",
		Enabled:     true,
	},
	{
		Name:        "Performance HUD",
		Description: "Displays a Heads-Up Display with performance metrics like FPS and memory usage.",
		Enabled:     false,
	},
}

// NewApp creates a new TUI application.
func NewApp() *App {
	app := tview.NewApplication()
	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	logView.SetBorder(true).SetTitle("Logs")

	return &App{
		Application: app,
		LogView:     logView,
	}
}

// LogWriter returns an io.Writer that streams log data into the log view.
func (a *App) LogWriter() io.Writer {
	return tview.ANSIWriter(a.LogView)
}

// Run starts the TUI application and builds the layout.
func (a *App) Run() error {
	isNormalMode := false

	featureState[0].Action = func(enabled bool) {
		fmt.Fprintf(a.LogWriter(), "Verbose logging toggled: %v\n", enabled)
	}
	featureState[1].Action = func(enabled bool) {
		fmt.Fprintf(a.LogWriter(), "Mock server responses toggled: %v\n", enabled)
	}
	featureState[2].Action = func(enabled bool) {
		fmt.Fprintf(a.LogWriter(), "Performance HUD toggled: %v\n", enabled)
	}

	table := tview.NewTable().
		SetSelectable(true, false)
	table.SetBorder(true).SetTitle("Features")

	infoBar := tview.NewTextView().SetDynamicColors(true)
	infoBar.SetBorder(true).SetTitle("Info")

	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	updateLayout := func() {
		rootFlex.Clear()
		if isNormalMode {
			mainContent := tview.NewFlex().
				AddItem(table, 0, 1, true).
				AddItem(a.LogView, 0, 2, false)
			rootFlex.AddItem(mainContent, 0, 1, false).
				AddItem(infoBar, 3, 0, false)
			a.SetFocus(table)
		} else {
			rootFlex.AddItem(a.LogView, 0, 1, true)
			a.SetFocus(a.LogView)
		}
	}

	populateTable := func() {
		current_row, _ := table.GetSelection()
		table.Clear()
		for i, feature := range featureState {
			table.SetCell(i, 0,
				tview.NewTableCell(fmt.Sprintf("[yellow](%d)[-]", i+1)).
					SetAlign(tview.AlignCenter))
			checkbox := "[ ]"
			if feature.Enabled {
				checkbox = "[[green]X[-]]"
			}
			table.SetCell(i, 1,
				tview.NewTableCell(fmt.Sprintf(" %s ", checkbox)).
					SetAlign(tview.AlignCenter))
			table.SetCell(i, 2,
				tview.NewTableCell(feature.Name).
					SetExpansion(1))
		}
		table.SetCell(len(featureState), 0, tview.NewTableCell("[yellow](q)[-]").
			SetAlign(tview.AlignCenter))
		table.SetCell(len(featureState), 1, tview.NewTableCell(""))
		table.SetCell(len(featureState), 2, tview.NewTableCell("Quit"))
		table.Select(current_row, 0)
	}

	table.SetSelectionChangedFunc(func(row, column int) {
		if row >= 0 && row < len(featureState) {
			infoBar.SetText(featureState[row].Description)
		} else if row == len(featureState) {
			infoBar.SetText("Press Enter or 'q' to exit the application.")
		}
	})

	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			isNormalMode = !isNormalMode
			updateLayout()
			return nil
		}

		if isNormalMode {
			r := event.Rune()
			featureIndex := -1
			switch r {
			case '1', '+':
				featureIndex = 0
			case '2', 'ě':
				featureIndex = 1
			case '3', 'š':
				featureIndex = 2
			}

			if featureIndex != -1 {
				featureState[featureIndex].Enabled = !featureState[featureIndex].Enabled
				if featureState[featureIndex].Action != nil {
					featureState[featureIndex].Action(featureState[featureIndex].Enabled)
				}
				populateTable()
				return nil
			}

			if r == 'q' {
				a.Stop()
				return nil
			}
		}
		return event
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			row, _ := table.GetSelection()
			if row >= 0 && row < len(featureState) {
				featureIndex := row
				featureState[featureIndex].Enabled = !featureState[featureIndex].Enabled
				if featureState[featureIndex].Action != nil {
					featureState[featureIndex].Action(featureState[featureIndex].Enabled)
				}
				populateTable()
			} else if row == len(featureState) {
				a.Stop()
			}
			return nil
		}
		return event
	})

	populateTable()
	updateLayout()

	fmt.Fprintln(a.LogWriter(), "[yellow]Press 'Esc' to toggle command mode.[-]")

	a.SetRoot(rootFlex, true).EnableMouse(true)

	if err := a.Application.Run(); err != nil {
		return err
	}

	return nil
}
