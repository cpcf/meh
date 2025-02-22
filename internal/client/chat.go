package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

func ChatTui(api API) {
	p := tea.NewProgram(initialModel(api))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type newLlmMsg struct{}
type contLlmMsg struct{}

type (
	errMsg error
)

type model struct {
	api          API
	viewport     viewport.Model
	messages     []string
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	results      chan string
	waitingOnLlm bool
	err          error
}

func initialModel(api API) model {
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

	return model{
		api:          api,
		textarea:     ta,
		messages:     []string{},
		viewport:     vp,
		senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		waitingOnLlm: false,
		err:          nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, tea.Quit
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

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil

	// Add a new LLM message to the history
	case newLlmMsg:
		m.messages = append(m.messages, m.senderStyle.Render("Assistant: ")+m.textarea.Value())
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
				return contLlmMsg{}
			}
		}
		m.waitingOnLlm = false
		return m, nil
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
