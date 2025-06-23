// Package renderer provides table rendering functionality and theme system for BubbleTable.
//
// This package handles the visual presentation of table data, including styling,
// themes, and layout calculations. It bridges the headless table core with
// terminal-based rendering using lipgloss for styling.
//
// # Core Types
//
//   - Renderer: Main rendering engine that converts table data to styled strings
//   - Theme: Defines the visual appearance with colors and styles
//   - StyleSet: Collection of lipgloss styles for different table elements
//
// # Key Features
//
//   - Multiple built-in themes with consistent styling
//   - Automatic column width distribution and text truncation
//   - Sort indicators and selection highlighting
//   - Responsive layout that adapts to terminal size
//   - Custom theme creation and modification
//   - Efficient rendering optimized for terminal output
//
// # Built-in Themes
//
// BubbleTable includes several professionally designed themes:
//
//   - DefaultTheme: Clean, minimal design suitable for any terminal
//   - DraculaTheme: Dark theme with purple and pink accents
//   - MonokaiTheme: Popular dark theme with vibrant colors
//   - GithubTheme: Clean theme inspired by GitHub's interface
//   - TerminalTheme: High-contrast theme for accessibility
//   - SolarizedDarkTheme: Based on the popular Solarized color scheme
//   - SolarizedLightTheme: Light variant of the Solarized theme
//
// # Usage Examples
//
// Basic rendering with default theme:
//
//	r := renderer.New().WithTheme(renderer.DefaultTheme)
//	output := r.Render(table, width, height)
//
// Creating a custom theme:
//
//	customTheme := renderer.CustomizeTheme(
//	    renderer.DefaultTheme,
//	    "My Theme",
//	    map[string]lipgloss.Style{
//	        "header": lipgloss.NewStyle().
//	            Background(lipgloss.Color("#FF6B6B")).
//	            Foreground(lipgloss.Color("#FFFFFF")).
//	            Bold(true),
//	        "selected": lipgloss.NewStyle().
//	            Background(lipgloss.Color("#4ECDC4")).
//	            Foreground(lipgloss.Color("#000000")),
//	    },
//	)
//
// Applying themes dynamically:
//
//	r.WithTheme(renderer.DraculaTheme)
//	darkOutput := r.Render(table, width, height)
//
//	r.WithTheme(renderer.GithubTheme)
//	lightOutput := r.Render(table, width, height)
//
// # Theme System
//
// Themes consist of styled elements for different table components:
//
//   - header: Column headers styling
//   - row: Regular row styling
//   - selected: Currently selected row highlighting
//   - cell: Individual cell styling
//   - border: Table borders and separators
//   - scrollbar: Pagination indicators
//   - search: Search mode styling
//   - help: Help text styling
//
// Each theme maintains consistency across all elements while providing
// unique visual identity.
//
// # Custom Theme Creation
//
// Create themes by defining styles for each element:
//
//	myTheme := renderer.Theme{
//	    Name: "Corporate",
//	    Styles: map[string]lipgloss.Style{
//	        "header": lipgloss.NewStyle().
//	            Background(lipgloss.Color("#1E3A8A")).
//	            Foreground(lipgloss.Color("#FFFFFF")).
//	            Bold(true).
//	            Padding(0, 1),
//	        "selected": lipgloss.NewStyle().
//	            Background(lipgloss.Color("#3B82F6")).
//	            Foreground(lipgloss.Color("#FFFFFF")),
//	        // ... other styles
//	    },
//	}
//
// Or customize existing themes:
//
//	corporateTheme := renderer.CustomizeTheme(
//	    renderer.DefaultTheme,
//	    "Corporate",
//	    map[string]lipgloss.Style{
//	        "header": lipgloss.NewStyle().Background(lipgloss.Color("#1E3A8A")),
//	    },
//	)
//
// # Layout and Sizing
//
// The renderer automatically handles:
//   - Column width distribution based on content and constraints
//   - Text truncation with ellipsis for overflow
//   - Responsive layout that adapts to terminal size changes
//   - Proper alignment for different data types
//   - Border and separator rendering
//
// # Performance
//
// The renderer is optimized for terminal rendering:
//   - Efficient string building to minimize allocations
//   - Lazy evaluation of styles to reduce computation
//   - Caching of layout calculations for consistent performance
//   - Minimal terminal escape sequence generation
//
// # Accessibility
//
// Themes are designed with accessibility in mind:
//   - High contrast options available
//   - ANSI color fallbacks for limited terminals
//   - Clear visual hierarchy and indicators
//   - Support for screen readers through semantic structure
package renderer
