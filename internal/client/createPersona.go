package client

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpcf/meh/internal/ollama"
)

const (
	statusNormal state = iota
	stateDone
)

type CreatePersonaModel struct {
	config *Config
	lg     *lipgloss.Renderer
	styles *Theme
	form   *huh.Form
	width  int
	height int
	done   bool
}

func NewCreatePersonaModel(c *Config) CreatePersonaModel {
	m := CreatePersonaModel{
		width:  maxWidth,
		config: c,
	}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewTheme(m.lg)
	var (
		url string
	)
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Persona name").
				Validate(func(str string) error {
					if _, ok := FindPersona(*m.config, str); ok {
						return errors.New("That persona already exists.")
					}
					if str == "" {
						return errors.New("Name cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Key("url").
				Value(&url).
				Title("API URL").
				Validate(func(str string) error {
					api := ollama.NewAPI(str, "", "")
					if !api.Verify() {
						return errors.New("Could not connect to API")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Key("model").
				Title("Model").
				OptionsFunc(func() []huh.Option[string] {
					if url == "" {
						return []huh.Option[string]{}
					}
					a := ollama.NewAPI(url, "", "")
					m := a.Models()
					return huh.NewOptions(m...)
				}, &url),
		),
		huh.NewGroup(
			huh.NewText().
				Key("prompt").
				Title("System Prompt").
				Placeholder("Optional"),
			huh.NewConfirm().
				Key("default").
				Title("Set as default persona?").
				Affirmative("Yes").
				Negative("No"),
			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).
		WithWidth(50).
		WithHeight(m.height - 10).
		WithShowHelp(false).
		WithShowErrors(false).
		WithStrictGroupBactracking(false)

	return m
}

func (m CreatePersonaModel) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), tea.WindowSize())
}

func (m CreatePersonaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base().GetHorizontalFrameSize()
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, BackToMain
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if !m.done && m.form.State == huh.StateCompleted {
		persona := Persona{}
		persona.Name = m.form.GetString("name")
		persona.APIURL = m.form.GetString("url")
		persona.Model = m.form.GetString("model")
		persona.SystemPrompt = m.form.GetString("prompt")
		m.config.AddPersona(persona, m.form.GetBool("default"))
		cmds = append(cmds, BackToMain, SetPersonaCmd(persona))
		m.done = true
	}

	return m, tea.Batch(cmds...)
}

func (m CreatePersonaModel) View() string {
	s := m.styles
	switch m.form.State {
	case huh.StateCompleted:
		return ""
	default:

		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		p := Persona{
			Name:         m.form.GetString("name"),
			APIURL:       m.form.GetString("url"),
			Model:        m.form.GetString("model"),
			SystemPrompt: m.form.GetString("prompt"),
		}
		status := CreateStatusBar(s, p, m.width-lipgloss.Width(form), m.Height()-8, "Current Persona")

		errors := m.form.Errors()
		header := appBoundaryView(&m, "Persona Creator")
		if len(errors) > 0 {
			header = appErrorBoundaryView(&m, m.errorView())
		}

		body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)
		footer := appBoundaryView(&m, m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = appErrorBoundaryView(&m, "")
		}

		return s.Base().Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m CreatePersonaModel) Height() int {
	return m.height
}
func (m CreatePersonaModel) Width() int {
	return m.width
}

func (m *CreatePersonaModel) SetHeight(height int) {
	m.height = height
}
func (m *CreatePersonaModel) SetWidth(width int) {
	m.width = width
}

func (m CreatePersonaModel) Styles() *Theme {
	return m.styles
}

func (m CreatePersonaModel) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}
