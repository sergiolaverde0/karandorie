package calendar

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Date represents a date
type Date struct {
	Year  int
	Month int
	Day   int
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
	selectedDay  int
}

// InitialModel creates a new calendar model for the current month
func InitialModel() CalendarModel {
	t := time.Now()
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
		selectedDay: t.Day(),
		days:        make([][]bool, 6),
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
		switch msg.String() {
		case "n":
			m.NextMonth()
		case "p":
			m.PreviousMonth()
		case "s":
			m.SelectDate()
		case " ":
			m.SelectDate()
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
	m.selectedDate.Day = m.selectedDay
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
				if i*7+j+1 == m.selectedDay {
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
