package editor

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	data "github.com/oscgu/snot/pkg/cli/dataproviders"
	theme "github.com/oscgu/snot/pkg/cli/note/ui/theme"
)

type EditorMode uint

const (
	Edit EditorMode = iota
	View
)

var (
	titleStyle = lipgloss.NewStyle().Foreground(theme.Blue).Bold(true)
	dateStyle  = lipgloss.NewStyle().Foreground(theme.Pink)

	cancelStyle = lipgloss.NewStyle().Foreground(theme.Red).Bold(true)
	saveStyle   = lipgloss.NewStyle().Foreground(theme.Green).Bold(true)
)

type EditorModel struct {
	textarea     textarea.Model
	title        string
	date         time.Time
	mode         EditorMode
	cancelled    bool
	dataProvider data.DataProvider
}

func CreateModel(ta textarea.Model, title string, mode EditorMode, date time.Time) EditorModel {
	initModel := EditorModel{
		textarea:     ta,
		title:        title,
		mode:         mode,
		date:         date,
		dataProvider: data.GetProvider(),
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

func Create(topic string, title string, date time.Time) (string, bool, time.Time) {
	ta := textarea.New()
	ta.CharLimit = 200
	ta.ShowLineNumbers = true
	ta.FocusedStyle.CursorLine.Background(lipgloss.NoColor{})
	ta.FocusedStyle.Prompt.Foreground(theme.PurpBlue)
	ta.Focus()

	m := CreateModel(ta, title, Edit, date)
	note, _ := m.dataProvider.GetNote(topic, title)

	ta.SetValue(note.Content)

	p := tea.NewProgram(&m)

	if err := p.Start(); err != nil {
		fmt.Println(err)
	}

	return m.textarea.Value(), m.cancelled, m.date
}
