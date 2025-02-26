package client2

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ()

type HeaderModel struct {
	timer             timer.Model
	spinner           spinner.Model
	tabs              []string
	index             int
	appWidth          int
	appHeight         int
	width             int
	height            int
	count             int
	modelStyle        lipgloss.Style
	focusedModelStyle lipgloss.Style
}

func NewHeader(appWidth, appHeight int, tabs ...string) HeaderModel {
	m := HeaderModel{appWidth: appWidth, appHeight: appHeight, tabs: tabs}
	b := lipgloss.RoundedBorder()
	b.Top, b.TopLeft, b.TopRight, b.Right, b.Left = "", "", "", "", ""
	m.modelStyle = lipgloss.NewStyle().
		Width(m.appWidth/len(tabs)).
		Height(2).
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(b)
	m.focusedModelStyle = lipgloss.NewStyle().
		Width(m.appWidth/len(tabs)).
		Height(2).
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(b).
		Foreground(lipgloss.Color("69"))

	return m
}

func (m HeaderModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m HeaderModel) Update(msg tea.Msg) (HeaderModel, tea.Cmd) {
	prevIndex := m.index
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "shift+tab":
			m.index--
		case "tab":
			m.index++
		case tea.KeyF1.String():
			m.index = 0
		case tea.KeyF2.String():
			m.index = 1
		case tea.KeyF3.String():
			m.index = 2
		case tea.KeyF4.String():
			m.index = 3
		case tea.KeyF5.String():
			m.index = 4
		case tea.KeyF6.String():
			m.index = 5
		case tea.KeyF7.String():
			m.index = 6
		case tea.KeyF8.String():
			m.index = 7
		case tea.KeyF9.String():
			m.index = 8
		case tea.KeyF10.String():
			m.index = 9
		case tea.KeyF11.String():
			m.index = 10
		case tea.KeyF12.String():
			m.index = 11
		}
		if m.index < 0 {
			m.index = prevIndex
		} else if m.index >= len(m.tabs) {
			m.index = prevIndex
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m HeaderModel) View() string {
	var s string
	for i, tab := range m.tabs {
		if i == m.index {
			s = lipgloss.JoinHorizontal(lipgloss.Top, s, m.focusedModelStyle.Render(fmt.Sprintf("[f%d]\n%s", i+1, tab)))
		} else {
			s = lipgloss.JoinHorizontal(lipgloss.Top, s, m.modelStyle.Render(fmt.Sprintf("[f%d]\n%s", i+1, tab)))
		}
	}
	return CenterContent(m.appWidth, m.appHeight, m.width, lipgloss.Height(s), 0, 0, s)
}
