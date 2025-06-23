# BubbleTable Examples

This document provides practical examples for using BubbleTable in various scenarios.

## Basic Usage

### Simple Table with Struct Data

```go
package main

import (
    "log"
    "github.com/anurag-roy/bubbletable/components"
    tea "github.com/charmbracelet/bubbletea"
)

type Person struct {
    ID   int    `table:"ID,sortable,width:5"`
    Name string `table:"Name,sortable,width:20"`
    Age  int    `table:"Age,sortable,width:10"`
}

func main() {
    people := []Person{
        {1, "Alice Johnson", 30},
        {2, "Bob Smith", 25},
        {3, "Carol Brown", 35},
    }

    model := components.NewTable(people).
        WithPageSize(10).
        WithSorting(true).
        WithSearch(true)

    program := tea.NewProgram(model, tea.WithAltScreen())
    if _, err := program.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Custom Columns with Map Data

```go
columns := []table.Column{
    *table.NewColumn("name", "Employee Name").WithWidth(25),
    *table.NewColumn("department", "Department").WithWidth(15),
    *table.NewColumn("salary", "Salary").
        WithFormatter(table.CurrencyFormatter).
        WithWidth(12),
    *table.NewColumn("active", "Status").
        WithFormatter(table.BooleanFormatter).
        WithWidth(10),
}

data := []map[string]interface{}{
    {"name": "Alice", "department": "Engineering", "salary": 75000.0, "active": true},
    {"name": "Bob", "department": "Marketing", "salary": 65000.0, "active": true},
    {"name": "Carol", "department": "Sales", "salary": 70000.0, "active": false},
}

model := components.NewTableWithColumns(data, columns)
```

## Advanced Examples

### Table with Callbacks

```go
model := components.NewTable(employees).
    WithOnSelect(func(row table.Row) {
        fmt.Printf("Selected employee: %s\n", row["name"])
    }).
    WithOnSort(func(column string, ascending bool) {
        fmt.Printf("Sorted by %s (ascending: %t)\n", column, ascending)
    }).
    WithOnSearch(func(query string) {
        fmt.Printf("Searching for: %s\n", query)
    })
```

### Custom Theme Usage

```go
customTheme := renderer.CustomizeTheme(
    renderer.DefaultTheme,
    "Corporate",
    map[string]lipgloss.Style{
        "header": lipgloss.NewStyle().
            Background(lipgloss.Color("#1E40AF")).
            Foreground(lipgloss.Color("#FFFFFF")).
            Bold(true),
    },
)

model := components.NewTable(data).WithTheme(customTheme)
```

### Headless Usage

```go
// Create headless table for custom rendering
tbl := table.New().
    WithData(employees).
    WithPageSize(10)

// Custom render loop
for {
    page := tbl.GetPage()
    columns := tbl.GetColumns()

    // Custom rendering logic here
    for _, row := range page {
        fmt.Printf("Employee: %s\n", row["name"])
    }

    // Handle navigation
    if userInput == "next" {
        tbl.NextPage()
    }
}
```

## Integration Examples

### Embed in Larger Application

```go
type AppModel struct {
    table components.Model
    // other components
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
    return m, cmd
}

func (m AppModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        "My Application Header",
        m.table.View(),
        "Footer",
    )
}
```

For complete examples, see the `examples/` directory in the repository.
