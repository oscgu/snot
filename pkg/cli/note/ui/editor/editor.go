package editor

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EditorMode uint

const (
	lightBlue  = lipgloss.Color("#3490dc")
	lightGreen = lipgloss.Color("#38c172")
	lightRed   = lipgloss.Color("#e3342f")
	darkPurple = lipgloss.Color("#9561e2")
	purpBlue   = lipgloss.Color("#6574cd")
	orange     = lipgloss.Color("#f6993f")
	pink       = lipgloss.Color("#f66d9b")
)

const (
	Edit EditorMode = iota
	View
)

var (
	titleStyle = lipgloss.NewStyle().Foreground(lightBlue).Bold(true)
	dateStyle  = lipgloss.NewStyle().Foreground(pink)

	cancelStyle = lipgloss.NewStyle().Foreground(lightRed).Bold(true)
	saveStyle   = lipgloss.NewStyle().Foreground(lightGreen).Bold(true)
)

type EditorModel struct {
	textarea  textarea.Model
	err       error
	title     string
	date      time.Time
	mode      EditorMode
	cancelled bool
}

func CreateModel(ta textarea.Model, title string, mode EditorMode) EditorModel {
	initModel := EditorModel{
		textarea: ta,
		title:    title,
		err:      nil,
		mode:     mode,
		date:     time.Now(),
	}

	return initModel
}

func (m EditorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m EditorModel) View() string {
	termWidth := m.textarea.Width() + 6

	str := strings.Builder{}
	str.WriteString("\n" + dateStyle.Width(termWidth).Align(lipgloss.Right).Render(m.date.Format("2006-01-02")) + "\n\n")
	str.WriteString(titleStyle.Width(termWidth).Align(lipgloss.Center).Render(m.title) + "\n\n")
	str.WriteString(m.textarea.View() + "\n\n")

	if m.mode == Edit {
		str.WriteString(lipgloss.NewStyle().Width(termWidth).Align(lipgloss.Center).Render("ctrl+"+saveStyle.Render("s")+"("+saveStyle.Render("ave")+")"+"  "+"ctrl+"+cancelStyle.Render("c")+"("+cancelStyle.Render("ancel")+")") + "\n\n")
	}

	return str.String()
}

func Create(title string) (string, bool, time.Time) {
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
