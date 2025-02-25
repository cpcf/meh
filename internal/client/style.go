package client

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type StylableModel interface {
	Width() int
	SetWidth(int)
	SetHeight(int)
	Styles() *Theme
}

var (
	peach      = lipgloss.AdaptiveColor{Light: "#ef9f76", Dark: "#ef9f76"}
	mauve      = lipgloss.AdaptiveColor{Light: "#ca9ee6", Dark: "#ca9ee6"}
	pink       = lipgloss.AdaptiveColor{Light: "#f4b8e4", Dark: "#f4b8e4"}
	red        = lipgloss.AdaptiveColor{Light: "#e78284", Dark: "#e78284"}
	lavendar   = lipgloss.AdaptiveColor{Light: "#babbf1", Dark: "#babbf1"}
	blue       = lipgloss.AdaptiveColor{Light: "#8caaee", Dark: "#8caaee"}
	green      = lipgloss.AdaptiveColor{Light: "#a6d189", Dark: "#a6d189"}
	background = lipgloss.AdaptiveColor{Dark: "#303446", Light: "#303446"}
	border     = lipgloss.AdaptiveColor{Dark: "#D7DBDF", Light: "#D7DBDF"}
	body       = lipgloss.AdaptiveColor{Dark: "#c6d0f5", Light: "#c6d0f5"}
	accent     = lipgloss.AdaptiveColor{Dark: "#11181C", Light: "#11181C"}
	helptext   = lipgloss.AdaptiveColor{Dark: "#838ba7", Light: "#838ba7"}
	surface    = lipgloss.AdaptiveColor{Dark: "#414559", Light: "#414559"}
	subtext    = lipgloss.AdaptiveColor{Dark: "#a5adce", Light: "#a5adce"}
)

type Theme struct {
	body,
	accent,
	background,
	highlight,
	border,
	highlightBg,
	subtext,
	error lipgloss.TerminalColor
	base,
	header,
	status,
	statusHeader,
	help,
	errorHeader lipgloss.Style
}

func NewTheme(lg *lipgloss.Renderer) *Theme {
	base := lg.NewStyle().
		Padding(1, 4, 0, 1).
		Foreground(body).
		Background(background)
	headerText := lg.NewStyle().
		Foreground(mauve).
		Bold(true).
		Padding(0, 1, 0, 2)
	status := lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mauve).
		PaddingLeft(1).
		MarginTop(1)
	statusHeader := lg.NewStyle().
		Foreground(green).
		Bold(true)
	errorHeaderText := headerText.
		Foreground(red)
	help := lg.NewStyle().
		Foreground(helptext)
	theme := Theme{
		base:         base,
		body:         body,
		accent:       accent,
		highlightBg:  surface,
		background:   background,
		highlight:    pink,
		error:        red,
		header:       headerText,
		status:       status,
		statusHeader: statusHeader,
		help:         help,
		border:       border,
		subtext:      subtext,
		errorHeader:  errorHeaderText,
	}

	return &theme
}

func (b Theme) Body() lipgloss.TerminalColor {
	return b.body
}

func (b Theme) Highlight() lipgloss.TerminalColor {
	return b.highlight
}

func (b Theme) Background() lipgloss.TerminalColor {
	return b.background
}

func (b Theme) Accent() lipgloss.TerminalColor {
	return b.accent
}

func (b Theme) Base() lipgloss.Style {
	return b.base
}

func (b Theme) Status() lipgloss.Style {
	return b.status
}

func (b Theme) StatusHeader() lipgloss.Style {
	return b.statusHeader
}

func (b Theme) TextBody() lipgloss.Style {
	return b.Base().Foreground(b.body)
}

func (b Theme) TextAccent() lipgloss.Style {
	return b.Base().Foreground(b.accent)
}

func (b Theme) TextHighlight() lipgloss.Style {
	return b.Base().Foreground(b.highlight)
}

func (b Theme) TextError() lipgloss.Style {
	return b.Base().Foreground(b.error)
}

func (b Theme) Border() lipgloss.TerminalColor {
	return b.border
}

func (b Theme) Error() lipgloss.TerminalColor {
	return b.error
}

func (b Theme) ErrorHeader() lipgloss.Style {
	return b.errorHeader
}

func (b Theme) HeaderText() lipgloss.Style {
	return b.header
}

func (b Theme) Help() lipgloss.Style {
	return b.help
}

func appBoundaryView(m StylableModel, text string) string {
	return lipgloss.PlaceHorizontal(
		m.Width(),
		lipgloss.Left,
		m.Styles().HeaderText().Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(mauve),
	)
}

func appErrorBoundaryView(m StylableModel, text string) string {
	return lipgloss.PlaceHorizontal(
		m.Width(),
		lipgloss.Left,
		m.Styles().ErrorHeader().Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func UpdateWidth(m StylableModel, width int) {
	m.SetWidth(min(width, maxWidth) - m.Styles().Base().GetHorizontalFrameSize())
}

func UpdateHeight(m StylableModel, height int) {
	m.SetHeight(min(height, maxHeight) - m.Styles().Base().GetVerticalFrameSize())
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func NewListDelegate(t Theme) list.DefaultDelegate {
	// Create a new default delegate
	d := list.NewDefaultDelegate()

	// Change colors
	h := t.highlight
	n := t.body
	s := t.subtext
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(h).
		BorderLeftForeground(h)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(h).
		BorderLeftForeground(h)
	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Foreground(n).
		BorderLeftForeground(n)
	d.Styles.NormalDesc = d.Styles.NormalDesc.
		Foreground(s).
		BorderLeftForeground(s)

	return d
}

func NewListTheme(t Theme) list.Styles {
	s := list.DefaultStyles()
	s.TitleBar = s.TitleBar.
		Background(t.accent).
		Foreground(t.accent)
	s.Title = s.Title.
		Background(t.accent).
		Foreground(t.accent)
	s.StatusBar = s.StatusBar.
		Background(t.accent).
		Foreground(t.accent)
	return s
}
