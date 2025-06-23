package main

import (
	"fmt"
	"os"

	"github.com/anurag-roy/bubbletable/components"
	"github.com/anurag-roy/bubbletable/renderer"
	"github.com/anurag-roy/bubbletable/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Product represents a sample product struct
type Product struct {
	SKU       string  `table:"SKU,sortable,width:10"`
	Name      string  `table:"Product Name,sortable,width:25"`
	Category  string  `table:"Category,sortable,width:15"`
	Price     float64 `table:"Price,sortable,width:10,format:currency"`
	Stock     int     `table:"Stock,sortable,width:8"`
	Available bool    `table:"Available,sortable,width:10"`
}

func main() {
	// Sample product data
	products := []Product{
		{"SKU-001", "Wireless Headphones", "Electronics", 99.99, 50, true},
		{"SKU-002", "Smart Watch", "Electronics", 249.99, 25, true},
		{"SKU-003", "Laptop Stand", "Accessories", 35.99, 100, true},
		{"SKU-004", "USB-C Cable", "Accessories", 12.99, 200, true},
		{"SKU-005", "Mechanical Keyboard", "Electronics", 129.99, 0, false},
		{"SKU-006", "Gaming Mouse", "Electronics", 69.99, 75, true},
		{"SKU-007", "Monitor", "Electronics", 299.99, 15, true},
		{"SKU-008", "Desk Lamp", "Furniture", 45.99, 30, true},
		{"SKU-009", "Phone Case", "Accessories", 19.99, 150, true},
		{"SKU-010", "Tablet Stand", "Accessories", 28.99, 80, true},
		{"SKU-011", "Webcam", "Electronics", 89.99, 40, true},
		{"SKU-012", "Microphone", "Electronics", 159.99, 20, true},
		{"SKU-013", "Speaker Set", "Electronics", 199.99, 35, true},
		{"SKU-014", "Power Bank", "Electronics", 39.99, 60, true},
		{"SKU-015", "Charging Pad", "Electronics", 29.99, 90, true},
	}

	// Create custom columns with renderers
	columns := []table.Column{
		*table.NewColumn("SKU", "SKU").
			WithType(table.String).
			WithWidth(10).
			WithSortable(true),
		*table.NewColumn("Name", "Product Name").
			WithType(table.String).
			WithWidth(25).
			WithSortable(true).
			WithRenderer(func(val interface{}, selected bool) string {
				return fmt.Sprintf("📦 %s", val)
			}),
		*table.NewColumn("Category", "Category").
			WithType(table.String).
			WithWidth(15).
			WithSortable(true),
		*table.NewColumn("Price", "Price").
			WithType(table.Float).
			WithWidth(10).
			WithSortable(true).
			WithFormatter(table.CurrencyFormatter),
		*table.NewColumn("Stock", "Stock").
			WithType(table.Integer).
			WithWidth(8).
			WithSortable(true).
			WithRenderer(func(val interface{}, selected bool) string {
				stock := val.(int)
				if stock == 0 {
					return "❌ 0"
				} else if stock < 30 {
					return fmt.Sprintf("⚠️ %d", stock)
				}
				return fmt.Sprintf("✅ %d", stock)
			}),
		*table.NewColumn("Available", "Available").
			WithType(table.Boolean).
			WithWidth(10).
			WithSortable(true).
			WithFormatter(table.BooleanFormatter("✅ Yes", "❌ No")),
	}

	// Convert products to map format for custom columns
	data := make([]map[string]interface{}, len(products))
	for i, p := range products {
		data[i] = map[string]interface{}{
			"SKU":       p.SKU,
			"Name":      p.Name,
			"Category":  p.Category,
			"Price":     p.Price,
			"Stock":     p.Stock,
			"Available": p.Available,
		}
	}

	// Create custom theme
	customTheme := renderer.CustomizeTheme(&renderer.DraculaTheme, "Custom Corporate", map[string]lipgloss.Style{
		"Header": lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#1E3A8A")).
			Bold(true).
			Padding(0, 1),
		"SelectedRow": lipgloss.NewStyle().
			Background(lipgloss.Color("#3B82F6")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1),
		"Status": lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Background(lipgloss.Color("#F3F4F6")).
			Italic(true),
	})

	// Create table with custom theme and configuration
	tableModel := components.NewTableWithColumns(data, columns).
		WithPageSize(12).
		WithTheme(&customTheme).
		WithKeyBindings(components.VimKeyBindings()).
		WithSorting(true).
		WithSearch(true).
		WithOnSelect(func(row table.Row) {
			// Handle product selection
			fmt.Printf("Selected product: %s\n", row.Cells[0].Value)
		}).
		WithOnSort(func(columnIndex int, desc bool) {
			direction := "ascending"
			if desc {
				direction = "descending"
			}
			fmt.Printf("Sorted by column %d (%s)\n", columnIndex, direction)
		}).
		WithOnSearch(func(term string) {
			if term != "" {
				fmt.Printf("Searching for: %s\n", term)
			}
		})

	// Create and run the program
	p := tea.NewProgram(tableModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
