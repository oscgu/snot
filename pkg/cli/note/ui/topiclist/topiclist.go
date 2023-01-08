package topiclist

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ta "github.com/oscgu/snot/pkg/cli/note/ui/textarea"
	"github.com/oscgu/snot/pkg/cli/snotdb"
)

const listHeight = 10

type sessionState uint

const (
	topcicView sessionState = iota
	titleView
	noteView

	width int = 25
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)

	focusedModelStyle = lipgloss.NewStyle().
				Width(width).
				Height(5).
				Align(lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	modelStyle = lipgloss.NewStyle().
			Width(width).
			Height(5).
			Align(lipgloss.Center).
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
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	topicList list.Model
	titleList list.Model
	text      ta.Model
	state     sessionState
	selTopic  string
	selTitle  string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == titleView {
				m.state = topcicView
			} else {
				return m, tea.Quit
			}
		case "enter":
			if m.state == topcicView {
				m.state = titleView

				i, ok := m.topicList.SelectedItem().(item)
				if ok {
					m.selTopic = string(i)
					var titles []string
					snotdb.Db.Table("notes").Where("topic = ?", string(i)).Select("title").Find(&titles)

					m.titleList = makeTitleList(titles)
				}
			}
		}
		switch m.state {
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

func (m model) View() string {
	var s string

	if m.state == topcicView {
		s += focusedModelStyle.Render(fmt.Sprintf("%4s", m.topicList.View()))
	} else {
		if m.state == titleView {
			if item(m.selTitle) != m.titleList.SelectedItem().(item) {
				i, ok := m.titleList.SelectedItem().(item)
				if ok {
					m.selTitle = string(i)

					var content []string
					snotdb.Db.Table("notes").
						Select("content").
						Where("topic = ?", m.selTopic).
						Where("title = ?", m.selTitle).
						Find(&content)

					taC := textarea.New()
					taC.Placeholder = strings.Join(content, " ")
					m.text = ta.CreateModel(taC, m.selTitle, ta.View)
				}
			}
		}

		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", m.topicList.View())), focusedModelStyle.Render(m.titleList.View()), noteStyle.Render(m.text.View()))
	}

	return s
}

func List(items []string) {
	m := &model{topicList: makeTopicList(items)}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

func convertToItems(items []string) []list.Item {
	var itemList []list.Item

	for _, i := range items {
		itemList = append(itemList, item(i))
	}

	return itemList
}

func makeTopicList(items []string) list.Model {
	itemList := convertToItems(items)
	const defaultWidth = 80
	l := list.New(itemList, itemDelegate{}, defaultWidth, listHeight)
	l.SetWidth(defaultWidth)
	l.Title = "Topics"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return l
}

func makeTitleList(items []string) list.Model {
	const defaultWidth = 80
	titleList := convertToItems(items)
	titleListUi := list.New(titleList, itemDelegate{}, defaultWidth, listHeight)
	titleListUi.SetWidth(defaultWidth)
	titleListUi.DisableQuitKeybindings()
	titleListUi.SetFilteringEnabled(false)
	titleListUi.SetShowHelp(false)
	titleListUi.Title = "Titles"

	return titleListUi
}
