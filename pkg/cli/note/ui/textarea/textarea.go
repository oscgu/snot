package textarea

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaMode uint

const (
	lightBlue  = lipgloss.Color("#3490dc")
	lightGreen = lipgloss.Color("#38c172")
	lightRed   = lipgloss.Color("#e3342f")
	darkPurple = lipgloss.Color("#9561e2")
	purpBlue   = lipgloss.Color("#6574cd")
	orange     = lipgloss.Color("#f6993f")
	pink       = lipgloss.Color("#f66d9b")

	Edit TaMode = iota
	View
)

var (
	titleStyle = lipgloss.NewStyle().Foreground(lightBlue).Bold(true)
	dateStyle  = lipgloss.NewStyle().Foreground(pink)

	cancelStyle = lipgloss.NewStyle().Foreground(lightRed).Bold(true)
	saveStyle   = lipgloss.NewStyle().Foreground(lightGreen).Bold(true)
)

type Model struct {
	textarea  textarea.Model
	err       error
	title     string
	date      time.Time
	mode      TaMode
	cancelled bool
}

func CreateModel(ta textarea.Model, title string, mode TaMode) Model {
	initModel := Model{
		textarea: ta,
		title:    title,
		err:      nil,
		mode:     mode,
		date:     time.Now(),
	}

	return initModel
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			m.cancelled = true
			return m, tea.Quit
		case tea.KeyCtrlS:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	termWidth := m.textarea.Width() + 6

	if m.mode == Edit {
		return fmt.Sprintf(
			"\n%s\n\n%s\n\n%s\n\n%s",
			dateStyle.Width(termWidth).Align(lipgloss.Right).Render(m.date.Format("2006-01-02")),
			titleStyle.Width(termWidth).Align(lipgloss.Center).Render(m.title),
			m.textarea.View(),
			lipgloss.NewStyle().Width(termWidth).Align(lipgloss.Center).Render("ctrl+"+saveStyle.Render("s")+"("+saveStyle.Render("ave")+")"+"  "+"ctrl+"+cancelStyle.Render("c")+"("+cancelStyle.Render("ancel")+")"),
		) + "\n\n"
	}

	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s\n\n",
		dateStyle.Width(termWidth).Align(lipgloss.Right).Render(m.date.Format("2006-01-02")),
		titleStyle.Width(termWidth).Align(lipgloss.Center).Render(m.title),
		m.textarea.View(),
	) + "\n\n"
}

func TextArea(title string) (string, bool, time.Time) {
	ta := textarea.New()
	ta.CharLimit = 200
	ta.ShowLineNumbers = true
	ta.FocusedStyle.CursorLine.Background(lipgloss.NoColor{})
	ta.FocusedStyle.Prompt.Foreground(purpBlue)
	ta.Focus()

	m := CreateModel(ta, title, Edit)
	p := tea.NewProgram(&m)

	if err := p.Start(); err != nil {
		fmt.Println(err)
	}

	return m.textarea.Value(), m.cancelled, m.date
}
