# BubbleTable API Reference

This document provides detailed API reference for all BubbleTable packages.

## Table of Contents

- [Core Package (`table`)](#core-package-table)
- [Renderer Package (`renderer`)](#renderer-package-renderer)
- [Components Package (`components`)](#components-package-components)
- [Examples](#examples)

## Core Package (`table`)

The `table` package provides the headless core functionality.

### Types

#### `Table`

The main table data structure that manages rows, columns, and state.

```go
type Table struct {
    // Contains filtered or unexported fields
}
```

**Methods:**

```go
func New() *Table
func (t *Table) WithData(data interface{}) *Table
func (t *Table) WithColumns(columns []Column) *Table
func (t *Table) WithPageSize(size int) *Table
func (t *Table) SetSortColumn(column string, ascending bool)
func (t *Table) SetFilter(query string)
func (t *Table) NextPage() bool
func (t *Table) PreviousPage() bool
func (t *Table) FirstPage()
func (t *Table) LastPage()
func (t *Table) GetPage() []Row
func (t *Table) GetColumns() []Column
func (t *Table) GetCurrentPage() int
func (t *Table) GetTotalPages() int
func (t *Table) GetTotalRows() int
func (t *Table) GetSelectedRowIndex() int
func (t *Table) SetSelectedRowIndex(index int)
```

#### `Column`

Defines column metadata, formatting, and behavior.

```go
type Column struct {
    Key        string
    Title      string
    Type       ColumnType
    Width      int
    Sortable   bool
    Searchable bool
    Formatter  FormatterFunc
    Renderer   RendererFunc
}
```

**Methods:**

```go
func NewColumn(key, title string) *Column
func (c *Column) WithType(columnType ColumnType) *Column
func (c *Column) WithWidth(width int) *Column
func (c *Column) WithSortable(sortable bool) *Column
func (c *Column) WithSearchable(searchable bool) *Column
func (c *Column) WithFormatter(formatter FormatterFunc) *Column
func (c *Column) WithRenderer(renderer RendererFunc) *Column
```

#### `Row`

Represents a single row of data.

```go
type Row map[string]interface{}
```

#### `ColumnType`

Enumeration of supported data types.

```go
type ColumnType int

const (
    String ColumnType = iota
    Int
    Float
    Bool
    Date
    Time
)
```

### Formatters

#### Built-in Formatters

```go
var (
    CurrencyFormatter         FormatterFunc
    PercentFormatter         FormatterFunc
    DateFormatter           FormatterFunc
    TimeFormatter           FormatterFunc
    BooleanFormatter        FormatterFunc
    NumberWithCommasFormatter FormatterFunc
)
```

#### Custom Formatters

```go
// TruncateFormatter creates a formatter that truncates text to specified length
func TruncateFormatter(maxLength int) FormatterFunc

// PrefixFormatter creates a formatter that adds a prefix to values
func PrefixFormatter(prefix string) FormatterFunc

// SuffixFormatter creates a formatter that adds a suffix to values
func SuffixFormatter(suffix string) FormatterFunc
```

#### Function Types

```go
type FormatterFunc func(value interface{}) string
type RendererFunc func(value interface{}, selected bool) string
```

### Struct Tags

The `table` package supports comprehensive struct tag configuration:

```go
type Example struct {
    ID          int     `table:"ID,sortable,width:5"`
    Name        string  `table:"Full Name,sortable,searchable,width:25"`
    Price       float64 `table:"Price,sortable,width:12,format:currency"`
    InStock     bool    `table:"Available,sortable,width:10"`
    Description string  `table:"Description,!sortable,!searchable,width:30"`
}
```

**Tag Options:**

- `sortable` / `!sortable`: Enable/disable sorting
- `searchable` / `!searchable`: Enable/disable search filtering
- `width:N`: Set column width in characters
- `format:type`: Apply built-in formatters (currency, date, percent, etc.)

## Renderer Package (`renderer`)

The `renderer` package handles visual presentation and themes.

### Types

#### `Renderer`

Main rendering engine that converts table data to styled strings.

```go
type Renderer struct {
    // Contains filtered or unexported fields
}
```

**Methods:**

```go
func New() *Renderer
func (r *Renderer) WithTheme(theme Theme) *Renderer
func (r *Renderer) Render(table *table.Table, width, height int) string
func (r *Renderer) RenderHeader(columns []table.Column, width int) string
func (r *Renderer) RenderRow(row table.Row, columns []table.Column, selected bool, width int) string
```

#### `Theme`

Defines the visual appearance with colors and styles.

```go
type Theme struct {
    Name   string
    Styles map[string]lipgloss.Style
}
```

### Built-in Themes

```go
var (
    DefaultTheme        Theme
    DraculaTheme        Theme
    MonokaiTheme        Theme
    GithubTheme         Theme
    TerminalTheme       Theme
    SolarizedDarkTheme  Theme
    SolarizedLightTheme Theme
)
```

### Theme Functions

```go
// CustomizeTheme creates a new theme by customizing an existing one
func CustomizeTheme(baseTheme Theme, name string, overrides map[string]lipgloss.Style) Theme

// GetTheme returns a theme by name, or DefaultTheme if not found
func GetTheme(name string) Theme

// ListThemes returns all available theme names
func ListThemes() []string
```

## Components Package (`components`)

The `components` package provides Bubble Tea integration.

### Types

#### `Model`

Main Bubble Tea model that implements the table interface.

```go
type Model struct {
    // Contains filtered or unexported fields
}
```

**Methods:**

```go
func NewTable(data interface{}) Model
func NewTableWithColumns(data interface{}, columns []table.Column) Model
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
```

#### Builder Methods

```go
func (m Model) WithTitle(title string) Model
func (m Model) WithPageSize(size int) Model
func (m Model) WithSorting(enabled bool) Model
func (m Model) WithSearch(enabled bool) Model
func (m Model) WithTheme(theme renderer.Theme) Model
func (m Model) WithKeyBindings(keyBindings KeyBindings) Model
func (m Model) WithOnSelect(callback func(row table.Row)) Model
func (m Model) WithOnSort(callback func(column string, ascending bool)) Model
func (m Model) WithOnSearch(callback func(query string)) Model
func (m Model) WithOnRefresh(callback func()) Model
func (m Model) WithOnPageChange(callback func(page int)) Model
```

#### `KeyBindings`

Configurable key binding system.

```go
type KeyBindings struct {
    Up         []string
    Down       []string
    Left       []string
    Right      []string
    FirstPage  []string
    LastPage   []string
    Search     []string
    Sort       map[int][]string
    PageSizeUp []string
    PageSizeDown []string
    Help       []string
    Quit       []string
    Select     []string
}
```

### Key Binding Presets

```go
func DefaultKeyBindings() KeyBindings
func VimKeyBindings() KeyBindings
func EmacsKeyBindings() KeyBindings
```

### State Management

```go
func (m Model) GetSelectedRow() table.Row
func (m Model) GetCurrentPage() int
func (m Model) GetTotalPages() int
func (m Model) GetSearchQuery() string
func (m Model) IsSearchMode() bool
func (m Model) GetSortColumn() string
func (m Model) GetSortDirection() bool
```

## Examples

### Basic Usage

```go
package main

import (
    "log"
    "github.com/anurag-roy/bubbletable/components"
    tea "github.com/charmbracelet/bubbletea"
)

type Employee struct {
    ID     int    `table:"ID,sortable,width:5"`
    Name   string `table:"Name,sortable,width:20"`
    Salary float64 `table:"Salary,sortable,width:12,format:currency"`
}

func main() {
    employees := []Employee{
        {1, "Alice Johnson", 75000},
        {2, "Bob Smith", 65000},
    }

    model := components.NewTable(employees).
        WithPageSize(10).
        WithSorting(true).
        WithSearch(true)

    program := tea.NewProgram(model, tea.WithAltScreen())
    if _, err := program.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Advanced Usage with Callbacks

```go
model := components.NewTable(data).
    WithOnSelect(func(row table.Row) {
        fmt.Printf("Selected row: %v\n", row)
    }).
    WithOnSort(func(column string, ascending bool) {
        fmt.Printf("Sorted by %s (ascending: %t)\n", column, ascending)
    }).
    WithOnSearch(func(query string) {
        fmt.Printf("Searching for: %s\n", query)
    })
```

### Headless Usage

```go
package main

import (
    "fmt"
    "github.com/anurag-roy/bubbletable/table"
    "github.com/anurag-roy/bubbletable/renderer"
)

func main() {
    // Create headless table
    tbl := table.New().WithData(myData).WithPageSize(20)

    // Create renderer
    r := renderer.New().WithTheme(renderer.DraculaTheme)

    // Custom rendering loop
    for !done {
        output := r.Render(tbl, termWidth, termHeight)
        fmt.Print(output)

        // Handle input and update table state
        // tbl.NextPage(), tbl.SetSortColumn(), etc.
    }
}
```

### Custom Themes

```go
customTheme := renderer.CustomizeTheme(
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

model := components.NewTable(data).WithTheme(customTheme)
```

### Custom Formatters and Renderers

```go
// Custom formatter
statusFormatter := func(value interface{}) string {
    if active, ok := value.(bool); ok && active {
        return "✅ Active"
    }
    return "❌ Inactive"
}

// Custom renderer with colors
priorityRenderer := func(value interface{}, selected bool) string {
    priority := value.(string)
    style := lipgloss.NewStyle()

    switch priority {
    case "High":
        style = style.Foreground(lipgloss.Color("#FF0000"))
    case "Medium":
        style = style.Foreground(lipgloss.Color("#FFA500"))
    case "Low":
        style = style.Foreground(lipgloss.Color("#00FF00"))
    }

    if selected {
        style = style.Background(lipgloss.Color("#0000FF"))
    }

    return style.Render(priority)
}

columns := []table.Column{
    *table.NewColumn("status", "Status").WithFormatter(statusFormatter),
    *table.NewColumn("priority", "Priority").WithRenderer(priorityRenderer),
}
```

For more examples, see the `examples/` directory in the repository.
