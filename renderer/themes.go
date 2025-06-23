package renderer

import "github.com/charmbracelet/lipgloss"

// Theme represents a complete styling theme for tables
type Theme struct {
	Name        string
	Header      lipgloss.Style
	Cell        lipgloss.Style
	SelectedRow lipgloss.Style
	Border      lipgloss.Style
	Status      lipgloss.Style
	Search      lipgloss.Style
}

// Predefined themes
var (
	// DefaultTheme is a clean, professional look
	DefaultTheme = Theme{
		Name: "Default",
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282A36")).
			Background(lipgloss.Color("#C4A9F4")).
			Bold(true).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().
			Padding(0, 1),
		SelectedRow: lipgloss.NewStyle().
			Background(lipgloss.Color("#44475A")).
			Foreground(lipgloss.Color("#F8F8F2")).
			Padding(0, 1),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#874BFD")),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4")).
			Italic(true),
		Search: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB86C")).
			Background(lipgloss.Color("#282A36")).
			Padding(0, 1),
	}

	// DraculaTheme is based on the popular Dracula color scheme
	DraculaTheme = Theme{
		Name: "Dracula",
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282A36")).
			Background(lipgloss.Color("#BD93F9")).
			Bold(true).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Padding(0, 1),
		SelectedRow: lipgloss.NewStyle().
			Background(lipgloss.Color("#6272A4")).
			Foreground(lipgloss.Color("#F8F8F2")).
			Bold(true).
			Padding(0, 1),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#BD93F9")),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4")).
			Italic(true),
		Search: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Background(lipgloss.Color("#44475A")).
			Padding(0, 1),
	}

	// MonokaiTheme is inspired by the Monokai color scheme
	MonokaiTheme = Theme{
		Name: "Monokai",
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Background(lipgloss.Color("#E6DB74")).
			Bold(true).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Padding(0, 1),
		SelectedRow: lipgloss.NewStyle().
			Background(lipgloss.Color("#75715E")).
			Foreground(lipgloss.Color("#F8F8F2")).
			Padding(0, 1),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#A6E22E")),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#75715E")).
			Italic(true),
		Search: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FD971F")).
			Background(lipgloss.Color("#49483E")).
			Padding(0, 1),
	}

	// GithubTheme is inspired by GitHub's interface
	GithubTheme = Theme{
		Name: "GitHub",
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#24292e")).
			Background(lipgloss.Color("#f6f8fa")).
			Bold(true).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#24292e")).
			Padding(0, 1),
		SelectedRow: lipgloss.NewStyle().
			Background(lipgloss.Color("#0366d6")).
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 1),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#e1e4e8")),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#586069")).
			Italic(true),
		Search: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0366d6")).
			Background(lipgloss.Color("#f1f8ff")).
			Padding(0, 1),
	}

	// TerminalTheme is a minimalist black and white theme
	TerminalTheme = Theme{
		Name: "Terminal",
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#ffffff")).
			Bold(true).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 1),
		SelectedRow: lipgloss.NewStyle().
			Background(lipgloss.Color("#ffffff")).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#ffffff")),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Italic(true),
		Search: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#ffffff")).
			Padding(0, 1),
	}

	// SolarizedDarkTheme is based on the Solarized Dark color scheme
	SolarizedDarkTheme = Theme{
		Name: "Solarized Dark",
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#002b36")).
			Background(lipgloss.Color("#268bd2")).
			Bold(true).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#839496")).
			Padding(0, 1),
		SelectedRow: lipgloss.NewStyle().
			Background(lipgloss.Color("#073642")).
			Foreground(lipgloss.Color("#93a1a1")).
			Padding(0, 1),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#586e75")),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#586e75")).
			Italic(true),
		Search: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#b58900")).
			Background(lipgloss.Color("#073642")).
			Padding(0, 1),
	}

	// SolarizedLightTheme is based on the Solarized Light color scheme
	SolarizedLightTheme = Theme{
		Name: "Solarized Light",
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fdf6e3")).
			Background(lipgloss.Color("#268bd2")).
			Bold(true).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#657b83")).
			Padding(0, 1),
		SelectedRow: lipgloss.NewStyle().
			Background(lipgloss.Color("#eee8d5")).
			Foreground(lipgloss.Color("#586e75")).
			Padding(0, 1),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#93a1a1")),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#93a1a1")).
			Italic(true),
		Search: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#b58900")).
			Background(lipgloss.Color("#eee8d5")).
			Padding(0, 1),
	}
)

// GetAllThemes returns a slice of all available themes
func GetAllThemes() []Theme {
	return []Theme{
		DefaultTheme,
		DraculaTheme,
		MonokaiTheme,
		GithubTheme,
		TerminalTheme,
		SolarizedDarkTheme,
		SolarizedLightTheme,
	}
}

// GetThemeByName returns a theme by name
func GetThemeByName(name string) (Theme, bool) {
	allThemes := GetAllThemes()
	for i := range allThemes {
		if allThemes[i].Name == name {
			return allThemes[i], true
		}
	}
	return Theme{}, false
}

// CustomizeTheme creates a new theme based on an existing theme with customizations
func CustomizeTheme(base *Theme, name string, customizations map[string]lipgloss.Style) Theme {
	theme := Theme{
		Name:        name,
		Header:      base.Header,
		Cell:        base.Cell,
		SelectedRow: base.SelectedRow,
		Status:      base.Status,
		Search:      base.Search,
	}

	// Apply customizations
	for field := range customizations {
		style := customizations[field]
		switch field {
		case "Header":
			theme.Header = style
		case "Cell":
			theme.Cell = style
		case "SelectedRow":
			theme.SelectedRow = style
		case "Status":
			theme.Status = style
		case "Search":
			theme.Search = style
		}
	}

	return theme
}
