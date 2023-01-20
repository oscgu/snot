package viewer

import (
	"fmt"
	"io"
	"strings"

	log "github.com/oscgu/snot/internal/log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	data "github.com/oscgu/snot/pkg/cli/dataproviders"
	editor "github.com/oscgu/snot/pkg/cli/note/ui/editor"
	theme "github.com/oscgu/snot/pkg/cli/note/ui/theme"
)

type viewState uint

const (
	listHeight = 15
	boxWidth   = 20
	boxHeight  = 18
	dot        = "•"
)

const (
	topcicView viewState = iota
	titleView
	noteView
)

var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2)
	itemStyle  = lipgloss.NewStyle().MarginLeft(2)

	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(theme.Purple)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(1).PaddingBottom(1)

	focusedModelStyle = lipgloss.NewStyle().
				Width(boxWidth).
				Height(boxHeight).
				Align(lipgloss.Left).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(theme.Blue)

	modelStyle = lipgloss.NewStyle().
			Width(boxWidth).
			Height(boxHeight).
			Align(lipgloss.Left).
			BorderStyle(lipgloss.HiddenBorder())

	noteStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render(dot + " " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type viewerModel struct {
	topicList    list.Model
	titleList    list.Model
	editor       editor.EditorModel
	view         viewState
	selTopic     string
	selTitle     string
	dataProvider data.DataProvider
}

func (m viewerModel) Init() tea.Cmd {
	return nil
}

func (m viewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.view == titleView {
				m.view = topcicView
			} else {
				return m, tea.Quit
			}
		case "enter":
			if m.view == topcicView {
				m.view = titleView

				i, ok := m.topicList.SelectedItem().(item)
				if ok {
					m.selTopic = string(i)

					titles := m.dataProvider.GetTitles(string(i))
					list := newBaseList(titles)
					list.DisableQuitKeybindings()
					list.SetFilteringEnabled(false)
					list.SetShowHelp(false)
					list.Title = "Titles"

					m.titleList = list
				}
			}
		}
		switch m.view {
		case topcicView:
			m.topicList, cmd = m.topicList.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.titleList, cmd = m.titleList.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m viewerModel) View() string {
	var s strings.Builder

	if m.view == topcicView {
		s.WriteString(focusedModelStyle.Render(fmt.Sprintf("%4s", m.topicList.View())))
	} else {
		if m.view == titleView {
			if item(m.selTitle) != m.titleList.SelectedItem().(item) {
				i, ok := m.titleList.SelectedItem().(item)
				if ok {
					m.selTitle = string(i)

					note, _ := m.dataProvider.GetNote(m.selTopic, m.selTitle)

					textArea := textarea.New()
					textArea.SetValue(note.Content)
					m.editor = editor.CreateModel(textArea, m.selTitle, editor.View, note.Created)
				}
			}
		}

		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", m.topicList.View())), focusedModelStyle.Render(m.titleList.View()), noteStyle.Render(m.editor.View())))
	}

	return s.String()
}

func Create() {
	provider := data.GetProvider()
	topics := provider.GetTopics()

	list := newBaseList(topics)
	list.Title = "Topics"
	list.Styles.Title = titleStyle
	list.Styles.PaginationStyle = paginationStyle
	list.Styles.HelpStyle = helpStyle
	list.DisableQuitKeybindings()

	m := &viewerModel{topicList: list, dataProvider: provider}

	if err := tea.NewProgram(m).Start(); err != nil {
		log.Fatal(err)
	}

}

func convertToItems(items []string) []list.Item {
	var itemList []list.Item

	for _, i := range items {
		itemList = append(itemList, item(i))
	}

	return itemList
}

func newBaseList(entries []string) list.Model {
	const defaultWidth = 20
	items := convertToItems(entries)
	itemList := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	itemList.SetFilteringEnabled(false)
	itemList.Styles.PaginationStyle = paginationStyle

	return itemList
}
