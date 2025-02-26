package client

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme interface {
	Name() string
	Background() lipgloss.Color
	BackgroundOver() lipgloss.Color
	SubText() lipgloss.Color
	Accent() lipgloss.Color
	Foreground() lipgloss.Color
	Error() lipgloss.Color
	Success() lipgloss.Color
}

type DefaultTheme struct{}

func (d DefaultTheme) Name() string                   { return "default" }
func (d DefaultTheme) Background() lipgloss.Color     { return lipgloss.Color("#000000") }
func (d DefaultTheme) BackgroundOver() lipgloss.Color { return lipgloss.Color("#333333") }
func (d DefaultTheme) SubText() lipgloss.Color        { return lipgloss.Color("#999999") }
func (d DefaultTheme) Accent() lipgloss.Color         { return lipgloss.Color("#5A56E0") }
func (d DefaultTheme) Foreground() lipgloss.Color     { return lipgloss.Color("#FFFFFF") }
func (d DefaultTheme) Error() lipgloss.Color          { return lipgloss.Color("#FE5F86") }
func (d DefaultTheme) Success() lipgloss.Color        { return lipgloss.Color("#02BA84") }

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

func NewStyles(lg *lipgloss.Renderer, theme Theme) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1).
		Background(theme.Background())
	s.HeaderText = lg.NewStyle().
		Foreground(theme.Accent()).
		Background(theme.Background()).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Background(theme.Background()).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Accent()).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Background(theme.Background()).
		Foreground(theme.Success()).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Background(theme.Background()).
		Foreground(theme.SubText())
	s.ErrorHeaderText = s.HeaderText.
		Background(theme.Background()).
		Foreground(theme.Error())
	s.Help = lg.NewStyle().
		Background(theme.Background()).
		Foreground(lipgloss.Color("240"))

	return &s
}

func appBoundaryView(m StylableModel, text string) string {
	return lipgloss.PlaceHorizontal(
		m.Width(),
		lipgloss.Left,
		m.Styles().HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(m.Styles().HeaderText.GetForeground()),
		lipgloss.WithWhitespaceBackground(m.Styles().HeaderText.GetBackground()),
	)
}

func appErrorBoundaryView(m StylableModel, text string) string {
	return lipgloss.PlaceHorizontal(
		m.Width(),
		lipgloss.Left,
		m.Styles().ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(m.Styles().ErrorHeaderText.GetForeground()),
		lipgloss.WithWhitespaceBackground(m.Styles().ErrorHeaderText.GetBackground()),
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
