package calendar

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Date represents a date
type Date struct {
	Year  int
	Month int
	Day   int
}

// Define keybindings for the calendar
type keyMap struct {
	nextMonth key.Binding
	prevMonth key.Binding
	nextDay   key.Binding
	prevDay   key.Binding
	nextWeek  key.Binding
	prevWeek  key.Binding
	selectDay key.Binding
	quit      key.Binding
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.nextMonth, k.prevMonth, k.selectDay, k.quit},
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.nextMonth, k.prevMonth, k.selectDay, k.quit}
}

// CalendarModel holds the state of the calendar
type CalendarModel struct {
	currentMonth time.Month
	currentYear  int
	selectedDate *Date
	currentView  string
	monthName    string
	daysInMonth  int
	days         [][]bool
	currentDay   int
	keyMap       keyMap
}

// InitialModel creates a new calendar model for the current month
func InitialModel() CalendarModel {
	t := time.Now()

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
		nextDay: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("l", "next day"),
		),
		prevDay: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("h", "previous day"),
		),
		nextWeek: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j", "next week"),
		),
		prevWeek: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("j", "previous week"),
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

	model := CalendarModel{
		currentMonth: t.Month(),
		currentYear:  t.Year(),
		selectedDate: &Date{
			Year:  t.Year(),
			Month: int(t.Month()),
			Day:   t.Day(),
		},
		currentView: "month",
		monthName:   t.Month().String(),
		currentDay:  t.Day(),
		days:        make([][]bool, 6),
		keyMap:      keyMap,
	}
	// Get the first day of the month
	firstDay := time.Date(model.currentYear, model.currentMonth, 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(model.currentYear, model.currentMonth+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

	model.daysInMonth = lastDay.Day()
	for i := range model.days {
		model.days[i] = make([]bool, 7)
	}
	// Find the first day of the month (0=Sunday, 6=Saturday)
	dayOfWeek := int(firstDay.Weekday())
	for day := range dayOfWeek {
		model.days[0][day] = true // blank days before first day
	}
	for i := 0; i < model.daysInMonth-1; i++ {
		// Calculate the row and column
		row := (i + dayOfWeek) / 7
		col := (i + dayOfWeek) % 7
		model.days[row][col] = true
	}
	return model
}

// Update handles messages
func (m CalendarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.nextMonth):
			m.NextMonth()
		case key.Matches(msg, m.keyMap.prevMonth):
			m.PreviousMonth()
		case key.Matches(msg, m.keyMap.nextDay):
			m.currentDay += 1
		case key.Matches(msg, m.keyMap.prevDay):
			m.currentDay -= 1
		case key.Matches(msg, m.keyMap.nextWeek):
			m.currentDay += 7
		case key.Matches(msg, m.keyMap.prevWeek):
			m.currentDay -= 7
		case key.Matches(msg, m.keyMap.selectDay):
			m.SelectDate()
		case key.Matches(msg, m.keyMap.quit):
			return m, tea.Quit
		}
		return m, nil
	}
	return m, nil
}

func (m CalendarModel) Init() tea.Cmd {
	return nil
}

// NextMonth moves to the next month
func (m *CalendarModel) NextMonth() {
	nextMonth := time.Date(m.currentYear, m.currentMonth+1, 1, 0, 0, 0, 0, time.UTC)
	m.currentMonth = nextMonth.Month()
	m.currentYear = nextMonth.Year()
	m.monthName = nextMonth.Month().String()
	m.updateDays()
}

// PreviousMonth moves to the previous month
func (m *CalendarModel) PreviousMonth() {
	prevMonth := time.Date(m.currentYear, m.currentMonth-1, 1, 0, 0, 0, 0, time.UTC)
	m.currentMonth = prevMonth.Month()
	m.currentYear = prevMonth.Year()
	m.monthName = prevMonth.Month().String()
	m.updateDays()
}

// SelectDate selects the current day
func (m *CalendarModel) SelectDate() {
	m.selectedDate.Day = m.currentDay
}

// updateDays builds the calendar grid
func (model *CalendarModel) updateDays() {
	firstDay := time.Date(model.currentYear, model.currentMonth, 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(model.currentYear, model.currentMonth+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

	model.daysInMonth = lastDay.Day()
	for i := range model.days {
		model.days[i] = make([]bool, 7)
	}
	// Find the first day of the month (0=Sunday, 6=Saturday)
	dayOfWeek := int(firstDay.Weekday())
	for day := range dayOfWeek {
		model.days[0][day] = true // blank days before first day
	}
	for i := 0; i < model.daysInMonth-1; i++ {
		// Calculate the row and column
		row := (i + dayOfWeek) / 7
		col := (i + dayOfWeek) % 7
		model.days[row][col] = true
	}
}

// View draws the calendar
func (m CalendarModel) View() string {
	var sb []byte

	// Title
	sb = append(sb, fmt.Sprintf("Calendar - %s %d\n", m.monthName, m.currentYear)...)
	sb = append(sb, "----------------------------------------\n"...)
	sb = append(sb, "Sun Mon Tue Wed Thu Fri Sat\n"...)
	sb = append(sb, "----------------------------------------\n"...)

	// Print the days
	for i := range m.days {
		for j := range m.days[i] {
			if m.days[i][j] {
				// Highlight selected day
				if i*7+j+1 == m.currentDay {
					sb = append(sb, fmt.Sprintf("\x1b[1;32m%d\x1b[0m ", i*7+j+1)...)
				} else {
					sb = append(sb, fmt.Sprintf("%d ", i*7+j+1)...)
				}
			} else {
				sb = append(sb, ' ')
			}
		}
		sb = append(sb, "\n"...)
	}

	return string(sb)
}
