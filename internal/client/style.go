package client

import (
	"github.com/charmbracelet/lipgloss"
)

type StylableModel interface {
	Width() int
	SetWidth(int)
	SetHeight(int)
	Styles() *Styles
}

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

func appBoundaryView(m StylableModel, text string) string {
	return lipgloss.PlaceHorizontal(
		m.Width(),
		lipgloss.Left,
		m.Styles().HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func appErrorBoundaryView(m StylableModel, text string) string {
	return lipgloss.PlaceHorizontal(
		m.Width(),
		lipgloss.Left,
		m.Styles().ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func UpdateWidth(m StylableModel, width int) {
	m.SetWidth(min(width, maxWidth) - m.Styles().Base.GetHorizontalFrameSize())
}

func UpdateHeight(m StylableModel, height int) {
	m.SetHeight(min(height, maxHeight) - m.Styles().Base.GetVerticalFrameSize())
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
