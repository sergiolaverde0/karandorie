package calendar

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

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
	currentView  string
	monthName    string
	daysInMonth  int
	days         [][]bool
	selectedDay  time.Time
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
		monthName:    t.Month().String(),
		selectedDay:  t,
		days:         make([][]bool, 7),
		keyMap:       keyMap,
	}
	return model
}

// Update handles messages
func (m CalendarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	initialDate := m.selectedDay
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.nextMonth):
			m.NextMonth()
		case key.Matches(msg, m.keyMap.prevMonth):
			m.PreviousMonth()
		case key.Matches(msg, m.keyMap.nextDay):
			m.selectedDay = m.selectedDay.AddDate(0, 0, 1)
		case key.Matches(msg, m.keyMap.prevDay):
			m.selectedDay = m.selectedDay.AddDate(0, 0, -1)
		case key.Matches(msg, m.keyMap.nextWeek):
			m.selectedDay = m.selectedDay.AddDate(0, 0, 7)
		case key.Matches(msg, m.keyMap.prevWeek):
			m.selectedDay = m.selectedDay.AddDate(0, 0, -7)
		case key.Matches(msg, m.keyMap.quit):
			return m, tea.Quit
		}
	}
	m.currentMonth = m.selectedDay.Month()
	m.currentYear = m.selectedDay.Year()
	m.updateDays()
	return m, nil
}

func (m CalendarModel) Init() tea.Cmd {
	m.updateDays()
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

// updateDays builds the calendar grid
func (model *CalendarModel) updateDays() {
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
		model.days[0][day] = false // blank days before first day
	}
	for i := range model.daysInMonth {
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
	sb = append(sb, "----------------------------------\n"...)
	sb = append(sb, "Sun Mon Tue Wed Thu Fri Sat\n"...)
	sb = append(sb, "----------------------------------\n"...)

	// Print the days
	dayPrinted := 1
	for i := range m.days {
		for j := range m.days[i] {
			if m.days[i][j] {
				// Highlight selected day
				if dayPrinted == m.selectedDay.Day() {
					sb = append(sb, fmt.Sprintf("\x1b[1;32m%3d\x1b[0m ", dayPrinted)...)
				} else {
					sb = append(sb, fmt.Sprintf("%3d ", dayPrinted)...)
				}
				dayPrinted += 1
			} else {
				sb = append(sb, "    "...)
			}
		}
		sb = append(sb, "\n"...)
	}

	return string(sb)
}
