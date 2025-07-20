package main

import (
	"karandorie/pkg/calendar"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

// Help is the help component
var helpComp = help.New()

// Main model
type model struct {
	calendar.CalendarModel
	help   help.Model
}

func newModel() model {
	// Initialize the calendar
	calendarModel := calendar.InitialModel()

	// Initialize help
	helpComp = help.New()

	// Initialize the model
	model := model{
		CalendarModel: calendarModel,
		help:          helpComp,
	}
	return model
}

func main() {
	model := newModel()
	// Run the program
	if _, err := tea.NewProgram(model).Run(); err != nil {
		panic(err)
	}
}
