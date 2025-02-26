package themes

import (
	catppuccin "github.com/catppuccin/go"
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

type CatppuccinTheme struct {
	flavour catppuccin.Flavour
}

func (c CatppuccinTheme) Name() string               { return c.flavour.Name() }
func (c CatppuccinTheme) Background() lipgloss.Color { return lipgloss.Color(c.flavour.Base().Hex) }
func (c CatppuccinTheme) BackgroundOver() lipgloss.Color {
	return lipgloss.Color(c.flavour.Surface0().Hex)
}
func (c CatppuccinTheme) SubText() lipgloss.Color    { return lipgloss.Color(c.flavour.Subtext1().Hex) }
func (c CatppuccinTheme) Accent() lipgloss.Color     { return lipgloss.Color(c.flavour.Mauve().Hex) }
func (c CatppuccinTheme) Foreground() lipgloss.Color { return lipgloss.Color(c.flavour.Text().Hex) }
func (c CatppuccinTheme) Error() lipgloss.Color      { return lipgloss.Color(c.flavour.Red().Hex) }
func (c CatppuccinTheme) Success() lipgloss.Color    { return lipgloss.Color(c.flavour.Green().Hex) }

var (
	Mocha     = CatppuccinTheme{flavour: catppuccin.Mocha}
	Frappe    = CatppuccinTheme{flavour: catppuccin.Frappe}
	Macchiato = CatppuccinTheme{flavour: catppuccin.Macchiato}
	Latte     = CatppuccinTheme{flavour: catppuccin.Latte}
)
