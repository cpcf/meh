package client

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	mainState state = iota
	chatState
	selectPersonaState
	createPersonaState
)
const maxHeight = 1200
const maxWidth = 400
const statusWidth = 40
const statusMarginOffset = 18 // Longest string on Left side

type MainModel struct {
	currentState       state
	chatModel          ChatModel
	createPersonaModel CreatePersonaModel
	selectPersonaModel SelectPersonaModel
	config             *Config
	persona            Persona
	styles             *Styles
	width              int
	height             int
}

type switchMsg state
type personaMsg Persona

func BackToMain() tea.Msg {
	return switchMsg(mainState)
}

func setPersona(persona Persona) tea.Msg {
	return personaMsg(persona)
}

func SetPersonaCmd(persona Persona) tea.Cmd {
	return func() tea.Msg {
		return personaMsg(persona)
	}
}

func NewMainModel(c *Config, p Persona) MainModel {

	m := MainModel{
		currentState:       mainState,
		config:             c,
		createPersonaModel: NewCreatePersonaModel(c),
		styles:             NewStyles(lipgloss.DefaultRenderer()),
	}
	m.persona = p
	return m
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case personaMsg:
		m.persona = Persona(msg)
	case switchMsg:
		// Reload config when we switch
		conf, err := LoadConfig()

		if err != nil {
			if _, ok := err.(ErrNoConfig); ok {
				conf = &Config{}
				SaveConfig(conf)
			}
		}

		m.config = conf
		m.currentState = state(msg)
		return m, tea.WindowSize()
	}

	switch m.currentState {
	case chatState:
		updatedModel, cmd := m.chatModel.Update(msg)
		m.chatModel = updatedModel.(ChatModel)
		cmds = append(cmds, cmd)
	case createPersonaState:
		updatedModel, cmd := m.createPersonaModel.Update(msg)
		m.createPersonaModel = updatedModel.(CreatePersonaModel)
		cmds = append(cmds, cmd)
	case selectPersonaState:
		updatedModel, cmd := m.selectPersonaModel.Update(msg)
		m.selectPersonaModel = updatedModel.(SelectPersonaModel)
		cmds = append(cmds, cmd)
	case mainState:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				cmds = append(cmds, tea.Quit)
			case "c":
				m.currentState = chatState
				if m.persona != (Persona{}) {
					m.chatModel = NewChatModel(m.persona)
					cmds = append(cmds, m.chatModel.Init())
				}
			case "n":
				m.currentState = createPersonaState
				m.createPersonaModel = NewCreatePersonaModel(m.config)
				cmds = append(cmds, m.createPersonaModel.Init())
			case "r":
				m.currentState = selectPersonaState
				m.selectPersonaModel = NewPersonaListModel(m.config.Personas, m.persona)
				cmds = append(cmds, m.selectPersonaModel.Init())
			}
		case tea.WindowSizeMsg:
			UpdateWidth(&m, msg.Width)
			m.height = msg.Height
		}
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.currentState {
	case chatState:
		return m.chatModel.View()
	case selectPersonaState:
		return m.selectPersonaModel.View()
	case createPersonaState:
		return m.createPersonaModel.View()
	}
	return m.MainMenu()
}

func (m MainModel) MainMenu() string {
	s := m.styles
	// Current Persona (right side)
	status := CreateStatusBar(s, m.persona, m.width-statusMarginOffset, m.height-8, "Current Persona")

	header := appBoundaryView(&m, "meh")
	menu := "Main Menu:\n(c) Chat\n(r) List Personas\n(n) Create Persona\n(q) Quit"
	body := lipgloss.JoinHorizontal(lipgloss.Top, menu, status)

	// TODO: Add help for the main menu
	footer := appBoundaryView(&m, "")

	return s.Base.Render(header + "\n" + body + "\n\n" + footer)
}

func (m MainModel) Width() int {
	return m.width
}

func (m *MainModel) SetHeight(height int) {
	m.height = height
}
func (m *MainModel) SetWidth(width int) {
	m.width = width
}

func (m MainModel) Height() int {
	return m.height
}

func (m MainModel) Styles() *Styles {
	return m.styles
}

func CreateStatusBar(s *Styles, p Persona, width, height int, title string) string {
	var status string
	{
		var (
			name   string
			url    string
			model  string
			prompt string
		)
		if p != (Persona{}) {
			name = s.Highlight.Render("Name:  ")
			url = s.Highlight.Render("URL:   ")
			model = s.Highlight.Render("Model: ")
			prompt = s.Highlight.Render("Prompt:")

			name = fmt.Sprintf("%s%s\n", name, p.Name)
			url = fmt.Sprintf("%s%s\n", url, p.APIURL)
			model = fmt.Sprintf("%s%s\n", model, p.Model)
			prompt = fmt.Sprintf("%s\n%s\n", prompt, strings.Split(p.SystemPrompt, "\n")[0])

		}
		h, v := s.Base.GetFrameSize()
		statusMarginLeft := width - h - statusWidth - s.Status.GetMarginRight()
		status = s.Status.
			Height(height - v).
			Width(statusWidth).
			MarginLeft(statusMarginLeft).
			Render(s.StatusHeader.Render(title) + "\n" +
				name +
				url +
				model +
				prompt)
	}
	return status
}
