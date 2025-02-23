package client

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Persona represents a user-defined persona that overrides API/model settings and may include a system prompt.
type Persona struct {
	Name         string `yaml:"name"`
	APIURL       string `yaml:"api_url"`
	Model        string `yaml:"model"`
	SystemPrompt string `yaml:"system_prompt,omitempty"`
}

func (r Persona) String() string {
	return fmt.Sprintf("Persona{Name: %s, APIURL: %s, Model: %s, SystemPrompt: %s}", r.Name, r.APIURL, r.Model, r.SystemPrompt)
}

// FindPersona searches for a persona by name.
func FindPersona(conf Config, personaName string) (Persona, bool) {
	for _, r := range conf.Personas {
		if r.Name == personaName {
			return r, true
		}
	}
	return Persona{}, false
}

type item struct {
	persona Persona
}

const maxListWidth = 50
const listVerticalOffset = 7

func (i item) Title() string       { return i.persona.Name }
func (i item) Description() string { return i.persona.APIURL + " - " + i.persona.Model }
func (i item) FilterValue() string { return i.persona.Name }

type SelectPersonaModel struct {
	list           list.Model
	width          int
	height         int
	styles         *Styles
	lg             *lipgloss.Renderer
	currentPersona Persona
	delegate       list.DefaultDelegate
}

func NewPersonaListModel(personas []Persona, currentPersona Persona) SelectPersonaModel {
	items := make([]list.Item, len(personas))
	for i, persona := range personas {
		items[i] = item{persona: persona}
	}
	d := list.NewDefaultDelegate()
	lg := lipgloss.DefaultRenderer()
	s := NewStyles(lg)

	l := list.New(items, d, 0, 0)
	l.SetShowHelp(false)
	l.Title = "Personas"
	return SelectPersonaModel{list: l, lg: lg, styles: s, currentPersona: currentPersona}
}

func (m SelectPersonaModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m SelectPersonaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, func() tea.Msg { return switchMsg(mainState) }
		case "enter":
			item, ok := m.list.SelectedItem().(item)
			if !ok {
				return m, nil
			}
			m.currentPersona = item.persona
			return m, func() tea.Msg { return personaMsg(item.persona) }
		}
	case tea.WindowSizeMsg:
		UpdateWidth(&m, msg.Width)
		h, v := m.styles.Base.GetFrameSize()
		m.height = msg.Height - v
		m.list.SetSize(min(msg.Width-h, maxListWidth), msg.Height-v-listVerticalOffset)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SelectPersonaModel) View() string {
	s := m.styles
	// List (left side)
	v := m.list.View()
	list := s.Base.Render(v)

	var status string
	{
		name := fmt.Sprintf("%s%s\n", s.Highlight.Render("Name:  "), m.currentPersona.Name)
		url := fmt.Sprintf("%s%s\n", s.Highlight.Render("URL:   "), m.currentPersona.APIURL)
		model := fmt.Sprintf("%s%s\n", s.Highlight.Render("Model: "), m.currentPersona.Model)
		prompt := fmt.Sprintf("%s\n%s\n", s.Highlight.Render("Prompt:"), strings.Split(m.currentPersona.SystemPrompt, "\n")[0])
		h, v := m.styles.Base.GetFrameSize()
		statusMarginLeft := m.width - statusWidth - (h+1)*2 - trueWidth(m.list) - s.Status.GetMarginRight()
		status = s.Status.
			Height(m.list.Height() - v).
			Width(statusWidth).
			MarginLeft(statusMarginLeft).
			Render(s.StatusHeader.Render("Current Persona") + "\n" +
				name +
				url +
				model +
				prompt)

		header := appBoundaryView(&m, "select a persona")
		body := lipgloss.JoinHorizontal(lipgloss.Top, list, status)
		footer := appBoundaryView(&m, m.list.Help.ShortHelpView(m.list.ShortHelp()))

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m SelectPersonaModel) Height() int {
	return m.height
}
func (m SelectPersonaModel) Width() int {
	return m.width
}

func (m *SelectPersonaModel) SetHeight(height int) {
	m.height = height
}
func (m *SelectPersonaModel) SetWidth(width int) {
	m.width = width
}

func (m SelectPersonaModel) Styles() *Styles {
	return m.styles
}

// Gets the width of the widest item, or max width
func trueWidth(l list.Model) int {
	maxw := l.Width()
	w := 12
	for _, i := range l.Items() {
		i, ok := i.(item)
		if !ok {
			continue
		}
		if w < lipgloss.Width(i.Description()) {
			w = lipgloss.Width(i.Description())
		}
	}
	return min(maxw, w)
}
