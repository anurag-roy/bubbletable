// Package bubbletable provides a flexible, powerful, and beautiful table component library
// for Bubble Tea TUI applications.
//
// BubbleTable is inspired by TanStack React Table's API design and provides both a headless
// core for maximum flexibility and high-level components for quick implementation.
//
// # Features
//
//   - Headless Core: Full control over rendering with a headless table engine
//   - Beautiful Themes: Multiple built-in themes (Dracula, Monokai, GitHub, Solarized, etc.)
//   - Configurable Key Bindings: Default, Vim, and Emacs key binding presets
//   - Struct Tag Support: Automatic column inference from struct tags
//   - Real-time Search: Built-in search with customizable filters
//   - Smart Sorting: Multi-state sorting (unsorted → asc → desc → unsorted)
//   - Efficient Pagination: Handle large datasets with intelligent pagination
//   - Type-Safe: Full type safety with generics support
//   - Custom Formatters: Built-in formatters for currency, dates, percentages, etc.
//   - Custom Renderers: Add icons, colors, and custom styling to cells
//   - Performance Optimized: Efficient rendering and memory usage
//
// # Quick Start
//
// Basic usage with struct tags:
//
//	type Employee struct {
//	    ID         int     `table:"ID,sortable,width:5"`
//	    Name       string  `table:"Name,sortable,width:20"`
//	    Department string  `table:"Department,sortable,width:15"`
//	    Salary     float64 `table:"Salary,sortable,width:12,format:currency"`
//	    Active     bool    `table:"Active,sortable,width:8"`
//	}
//
//	func main() {
//	    employees := []Employee{
//	        {1, "Alice Johnson", "Engineering", 75000.0, true},
//	        {2, "Bob Smith", "Marketing", 65000.0, true},
//	    }
//
//	    tableModel := components.NewTable(employees).
//	        WithPageSize(10).
//	        WithSorting(true).
//	        WithSearch(true).
//	        WithTheme(renderer.DraculaTheme)
//
//	    p := tea.NewProgram(tableModel, tea.WithAltScreen())
//	    if _, err := p.Run(); err != nil {
//	        log.Fatal(err)
//	    }
//	}
//
// # Architecture
//
// BubbleTable follows a layered architecture:
//
//   - table: Core headless table functionality (data, sorting, filtering, pagination)
//   - renderer: Presentation layer with themes and styling
//   - components: Bubble Tea component integration
//
// This separation allows you to use the headless core for custom implementations
// or the high-level components for quick setup.
//
// # Package Structure
//
//   - table: Core table data structures and operations
//   - renderer: Table rendering and theme system
//   - components: Bubble Tea model and key bindings
//
// For detailed examples, see the examples/ directory in the repository.
package bubbletable
