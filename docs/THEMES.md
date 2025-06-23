# BubbleTable Themes Guide

BubbleTable includes a comprehensive theme system that allows you to customize the visual appearance of your tables. This guide covers built-in themes, custom theme creation, and advanced styling techniques.

## Table of Contents

- [Built-in Themes](#built-in-themes)
- [Using Themes](#using-themes)
- [Custom Theme Creation](#custom-theme-creation)
- [Theme Elements](#theme-elements)
- [Advanced Styling](#advanced-styling)
- [Accessibility](#accessibility)

## Built-in Themes

BubbleTable includes several professionally designed themes:

### Default Theme

Clean, minimal design suitable for any terminal environment.

- **Use case**: Professional applications, documentation
- **Colors**: Monochrome with subtle accents
- **Style**: Clean, minimal, readable

### Dracula Theme

Dark theme with purple and pink accents, inspired by the popular Dracula color scheme.

- **Use case**: Dark mode applications, creative projects
- **Colors**: Dark purple background, pink highlights
- **Style**: Modern, vibrant, eye-catching

### Monokai Theme

Popular dark theme with vibrant colors, inspired by the Monokai editor theme.

- **Use case**: Development tools, code-related applications
- **Colors**: Dark background with bright green, yellow, and orange accents
- **Style**: High contrast, developer-friendly

### GitHub Theme

Clean theme inspired by GitHub's interface design.

- **Use case**: Git-related tools, documentation viewers
- **Colors**: Light background with blue accents
- **Style**: Clean, professional, familiar

### Terminal Theme

High-contrast theme designed for maximum accessibility.

- **Use case**: Accessibility-focused applications, system administration
- **Colors**: High contrast black and white
- **Style**: Maximum readability, accessibility-first

### Solarized Dark Theme

Based on the popular Solarized color scheme (dark variant).

- **Use case**: Long-form reading, reduced eye strain
- **Colors**: Warm dark background with balanced accent colors
- **Style**: Eye-friendly, scientifically designed

### Solarized Light Theme

Light variant of the Solarized color scheme.

- **Use case**: Bright environments, documentation
- **Colors**: Warm light background with balanced accent colors
- **Style**: Eye-friendly, scientifically designed

## Using Themes

### Basic Theme Application

```go
import (
    "github.com/anurag-roy/bubbletable/components"
    "github.com/anurag-roy/bubbletable/renderer"
)

// Apply theme during table creation
model := components.NewTable(data).
    WithTheme(renderer.DraculaTheme)

// Or apply to renderer directly
r := renderer.New().WithTheme(renderer.MonokaiTheme)
```

### Dynamic Theme Switching

```go
// Create table with initial theme
model := components.NewTable(data).WithTheme(renderer.DefaultTheme)

// Switch themes dynamically
model = model.WithTheme(renderer.DraculaTheme)
model = model.WithTheme(renderer.GithubTheme)
```

### Theme Selection by Name

```go
// Get theme by name
theme := renderer.GetTheme("dracula")
if theme.Name == "" {
    theme = renderer.DefaultTheme // fallback
}

model := components.NewTable(data).WithTheme(theme)
```

### List Available Themes

```go
themes := renderer.ListThemes()
for _, themeName := range themes {
    fmt.Printf("Available theme: %s\n", themeName)
}
```

## Custom Theme Creation

### Creating a Theme from Scratch

```go
import "github.com/charmbracelet/lipgloss"

customTheme := renderer.Theme{
    Name: "Corporate Blue",
    Styles: map[string]lipgloss.Style{
        "header": lipgloss.NewStyle().
            Background(lipgloss.Color("#1E3A8A")).
            Foreground(lipgloss.Color("#FFFFFF")).
            Bold(true).
            Padding(0, 1).
            Align(lipgloss.Center),

        "row": lipgloss.NewStyle().
            Foreground(lipgloss.Color("#1F2937")).
            Padding(0, 1),

        "selected": lipgloss.NewStyle().
            Background(lipgloss.Color("#3B82F6")).
            Foreground(lipgloss.Color("#FFFFFF")).
            Bold(true).
            Padding(0, 1),

        "cell": lipgloss.NewStyle().
            Padding(0, 1),

        "border": lipgloss.NewStyle().
            Foreground(lipgloss.Color("#9CA3AF")),

        "search": lipgloss.NewStyle().
            Background(lipgloss.Color("#FEF3C7")).
            Foreground(lipgloss.Color("#92400E")).
            Bold(true).
            Padding(0, 1),

        "help": lipgloss.NewStyle().
            Foreground(lipgloss.Color("#6B7280")).
            Italic(true),
    },
}
```

### Customizing Existing Themes

```go
// Customize an existing theme
corporateTheme := renderer.CustomizeTheme(
    renderer.DefaultTheme,
    "Corporate",
    map[string]lipgloss.Style{
        "header": lipgloss.NewStyle().
            Background(lipgloss.Color("#1E3A8A")).
            Foreground(lipgloss.Color("#FFFFFF")).
            Bold(true),
        "selected": lipgloss.NewStyle().
            Background(lipgloss.Color("#3B82F6")).
            Foreground(lipgloss.Color("#FFFFFF")),
    },
)
```

### Theme Inheritance

```go
// Start with a base theme and override specific elements
nightTheme := renderer.CustomizeTheme(
    renderer.DraculaTheme,
    "Midnight",
    map[string]lipgloss.Style{
        "header": renderer.DraculaTheme.Styles["header"].
            Background(lipgloss.Color("#000000")),
        "selected": renderer.DraculaTheme.Styles["selected"].
            Background(lipgloss.Color("#1a1a1a")),
    },
)
```

## Theme Elements

Each theme consists of styled elements for different table components:

### Header (`"header"`)

Styling for column headers.

```go
"header": lipgloss.NewStyle().
    Background(lipgloss.Color("#4F46E5")).
    Foreground(lipgloss.Color("#FFFFFF")).
    Bold(true).
    Padding(0, 1).
    Align(lipgloss.Center)
```

### Row (`"row"`)

Default styling for table rows.

```go
"row": lipgloss.NewStyle().
    Foreground(lipgloss.Color("#374151")).
    Padding(0, 1)
```

### Selected (`"selected"`)

Styling for the currently selected row.

```go
"selected": lipgloss.NewStyle().
    Background(lipgloss.Color("#3B82F6")).
    Foreground(lipgloss.Color("#FFFFFF")).
    Bold(true).
    Padding(0, 1)
```

### Cell (`"cell"`)

Individual cell styling (applied to all cells).

```go
"cell": lipgloss.NewStyle().
    Padding(0, 1).
    Align(lipgloss.Left)
```

### Border (`"border"`)

Table borders and separators.

```go
"border": lipgloss.NewStyle().
    Foreground(lipgloss.Color("#D1D5DB"))
```

### Search (`"search"`)

Search mode indicator and input styling.

```go
"search": lipgloss.NewStyle().
    Background(lipgloss.Color("#FEF3C7")).
    Foreground(lipgloss.Color("#92400E")).
    Bold(true).
    Padding(0, 1)
```

### Help (`"help"`)

Help text and keyboard shortcuts.

```go
"help": lipgloss.NewStyle().
    Foreground(lipgloss.Color("#6B7280")).
    Italic(true)
```

### Scrollbar (`"scrollbar"`)

Pagination indicators and scroll elements.

```go
"scrollbar": lipgloss.NewStyle().
    Foreground(lipgloss.Color("#9CA3AF"))
```

## Advanced Styling

### Conditional Styling

```go
// Style based on data values
priorityRenderer := func(value interface{}, selected bool) string {
    priority := value.(string)
    style := lipgloss.NewStyle()

    switch priority {
    case "High":
        style = style.Foreground(lipgloss.Color("#EF4444"))
    case "Medium":
        style = style.Foreground(lipgloss.Color("#F59E0B"))
    case "Low":
        style = style.Foreground(lipgloss.Color("#10B981"))
    }

    if selected {
        style = style.Background(lipgloss.Color("#1E40AF"))
    }

    return style.Render(priority)
}
```

### Gradient Effects

```go
// Create gradient-like effects with color transitions
gradientHeader := lipgloss.NewStyle().
    Background(lipgloss.AdaptiveColor{
        Light: "#E0E7FF",
        Dark:  "#1E1B4B",
    }).
    Foreground(lipgloss.AdaptiveColor{
        Light: "#3730A3",
        Dark:  "#A5B4FC",
    })
```

### Responsive Styling

```go
// Adapt styles based on terminal capabilities
adaptiveTheme := renderer.Theme{
    Name: "Adaptive",
    Styles: map[string]lipgloss.Style{
        "header": lipgloss.NewStyle().
            Background(lipgloss.AdaptiveColor{
                Light: "#1E40AF",
                Dark:  "#3B82F6",
            }).
            Foreground(lipgloss.AdaptiveColor{
                Light: "#FFFFFF",
                Dark:  "#000000",
            }),
    },
}
```

### Animation and Effects

```go
// Add subtle effects (use sparingly in TUI)
shimmerStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#A855F7")).
    Bold(true).
    Italic(true)
```

## Color Schemes

### Popular Color Palettes

#### Material Design

```go
materialTheme := renderer.Theme{
    Name: "Material",
    Styles: map[string]lipgloss.Style{
        "header": lipgloss.NewStyle().
            Background(lipgloss.Color("#2196F3")).
            Foreground(lipgloss.Color("#FFFFFF")),
        "selected": lipgloss.NewStyle().
            Background(lipgloss.Color("#FFC107")).
            Foreground(lipgloss.Color("#000000")),
    },
}
```

#### Nord

```go
nordTheme := renderer.Theme{
    Name: "Nord",
    Styles: map[string]lipgloss.Style{
        "header": lipgloss.NewStyle().
            Background(lipgloss.Color("#5E81AC")).
            Foreground(lipgloss.Color("#ECEFF4")),
        "selected": lipgloss.NewStyle().
            Background(lipgloss.Color("#88C0D0")).
            Foreground(lipgloss.Color("#2E3440")),
    },
}
```

#### Gruvbox

```go
gruvboxTheme := renderer.Theme{
    Name: "Gruvbox",
    Styles: map[string]lipgloss.Style{
        "header": lipgloss.NewStyle().
            Background(lipgloss.Color("#458588")).
            Foreground(lipgloss.Color("#FBF1C7")),
        "selected": lipgloss.NewStyle().
            Background(lipgloss.Color("#B8BB26")).
            Foreground(lipgloss.Color("#1D2021")),
    },
}
```

## Accessibility

### High Contrast Theme

```go
highContrastTheme := renderer.Theme{
    Name: "High Contrast",
    Styles: map[string]lipgloss.Style{
        "header": lipgloss.NewStyle().
            Background(lipgloss.Color("#000000")).
            Foreground(lipgloss.Color("#FFFFFF")).
            Bold(true),
        "row": lipgloss.NewStyle().
            Foreground(lipgloss.Color("#000000")),
        "selected": lipgloss.NewStyle().
            Background(lipgloss.Color("#FFFFFF")).
            Foreground(lipgloss.Color("#000000")).
            Bold(true),
        "border": lipgloss.NewStyle().
            Foreground(lipgloss.Color("#000000")),
    },
}
```

### Color Blind Friendly Theme

```go
colorBlindFriendlyTheme := renderer.Theme{
    Name: "ColorBlind Friendly",
    Styles: map[string]lipgloss.Style{
        "header": lipgloss.NewStyle().
            Background(lipgloss.Color("#0173B2")).  // Blue
            Foreground(lipgloss.Color("#FFFFFF")),
        "selected": lipgloss.NewStyle().
            Background(lipgloss.Color("#CC78BC")).  // Pink
            Foreground(lipgloss.Color("#000000")),
        // Use patterns and shapes in addition to colors
    },
}
```

### Testing Accessibility

```go
// Test theme with different terminal capabilities
func testThemeAccessibility(theme renderer.Theme) {
    // Test with 16 colors
    lipgloss.SetColorProfile(termenv.TrueColor)

    // Test with 256 colors
    lipgloss.SetColorProfile(termenv.ANSI256)

    // Test with basic colors only
    lipgloss.SetColorProfile(termenv.ANSI)
}
```

## Best Practices

### Theme Design Guidelines

1. **Consistency**: Maintain visual consistency across all elements
2. **Contrast**: Ensure sufficient contrast for readability
3. **Accessibility**: Consider color-blind users and screen readers
4. **Terminal Compatibility**: Test with different terminal emulators
5. **Performance**: Avoid overly complex styles that slow rendering

### Theme Naming

```go
// Good theme names
"Corporate Blue"
"Dark Mode"
"High Contrast"
"Solarized Light"

// Avoid generic names
"Theme1"
"MyTheme"
"Custom"
```

### Color Selection

```go
// Use semantic color names
backgroundColor := lipgloss.Color("#1E40AF")  // Blue-800
textColor := lipgloss.Color("#F9FAFB")        // Gray-50

// Or use adaptive colors for terminal compatibility
adaptiveColor := lipgloss.AdaptiveColor{
    Light: "#1E40AF",
    Dark:  "#3B82F6",
}
```

For more examples and advanced techniques, see the `examples/custom_theme/` directory in the repository.
