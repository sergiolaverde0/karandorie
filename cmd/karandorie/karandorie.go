package main

import (
	"karandorie/pkg/calendar"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Define keybindings for the calendar
type keyMap struct {
	nextMonth key.Binding
	prevMonth key.Binding
	selectDay key.Binding
	quit      key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.nextMonth, k.prevMonth, k.selectDay, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.nextMonth, k.prevMonth, k.selectDay, k.quit},
	}
}

// Help is the help component
var helpComp = help.New()

// Main model
type model struct {
	calendar.CalendarModel
	help   help.Model
	keymap keyMap
}

func newModel() model {
	// Initialize the calendar
	calendarModel := calendar.InitialModel()

	// Initialize help
	helpComp = help.New()

	// Define keybindings
	keyMap := keyMap{
		nextMonth: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next month"),
		),
		prevMonth: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "previous month"),
		),
		selectDay: key.NewBinding(
			key.WithKeys("s", " "),
			key.WithHelp("s / space", "select day"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q / esc", "quit"),
		),
	}

	// Initialize the model
	model := model{
		CalendarModel: calendarModel,
		help:          helpComp,
		keymap:        keyMap,
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
