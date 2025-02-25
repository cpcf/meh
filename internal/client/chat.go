package client

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpcf/meh/internal/ollama"
)

type API interface {
	Models() []string
	SelectModel(model string)
	Chat(query string, results chan string, flag bool)
	Prompt(query string, results chan string, flag bool)
}

type ChatModel struct {
	api          API
	name         string
	ready        bool
	viewport     viewport.Model
	messages     []string
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	results      chan string
	waitingOnLlm bool
	err          error
}

func (m ChatModel) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.WindowSize())
}

const gap = "\n\n"

type newLlmMsg struct{}
type contLlmMsg struct{}

func NewChatModel(persona Persona) ChatModel {
	ta := textarea.New()
	ta.Placeholder = "Enter message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 750

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Interactive Mode.
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ChatModel{
		api:          ollama.NewAPI(persona.APIURL, persona.Model, persona.SystemPrompt),
		name:         persona.Name,
		ready:        true,
		textarea:     ta,
		messages:     []string{},
		viewport:     vp,
		senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		waitingOnLlm: false,
		err:          nil,
	}
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.ready {
		return m, func() tea.Msg { return switchMsg(mainState) }
	}

	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, func() tea.Msg { return switchMsg(mainState) }
		case tea.KeyEnter:
			if m.waitingOnLlm {
				return m, nil
			}
			message := m.textarea.Value()
			if strings.TrimSpace(message) == "" {
				break
			}

			m.results = make(chan string)
			go m.api.Chat(message, m.results, true)

			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()

			m.waitingOnLlm = true
			return m, func() tea.Msg {
				return newLlmMsg{}
			}
		}

	// Add a new LLM message to the history
	case newLlmMsg:
		m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%s: ", m.name))+m.textarea.Value())
		return m, func() tea.Msg {
			return contLlmMsg{}
		}
	// While results are still being streamed in add them to the latest message in the history
	case contLlmMsg:
		for res := range m.results {
			m.messages[len(m.messages)-1] = m.senderStyle.Render(m.messages[len(m.messages)-1] + res)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()
			return m, func() tea.Msg {
				// we return here so we can render the streamed results
				return contLlmMsg{}
			}
		}
		m.waitingOnLlm = false
		return m, nil
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m ChatModel) View() string {
	if !m.ready {
		return "No persona selected.\nPress any key to return."
	}
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
