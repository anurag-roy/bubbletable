// Package components provides Bubble Tea integration for BubbleTable.
//
// This package contains the high-level components that integrate the headless table
// core with the Bubble Tea framework, providing a complete TUI table experience
// with minimal setup required.
//
// # Core Types
//
//   - Model: Main Bubble Tea model that implements the table interface
//   - KeyBindings: Configurable key binding system
//   - Builder: Fluent API for table configuration
//
// # Key Features
//
//   - Complete Bubble Tea model implementation
//   - Fluent builder API for easy configuration
//   - Multiple key binding presets (Default, Vim, Emacs)
//   - Custom callback system for user interactions
//   - Built-in search mode with visual feedback
//   - Automatic state management and persistence
//   - Responsive design that adapts to terminal changes
//
// # Usage Examples
//
// Basic table with struct data:
//
//	type Employee struct {
//	    ID     int    `table:"ID,sortable,width:5"`
//	    Name   string `table:"Name,sortable,width:20"`
//	    Salary float64 `table:"Salary,sortable,width:12,format:currency"`
//	}
//
//	employees := []Employee{
//	    {1, "Alice Johnson", 75000},
//	    {2, "Bob Smith", 65000},
//	}
//
//	model := components.NewTable(employees).
//	    WithPageSize(10).
//	    WithSorting(true).
//	    WithSearch(true).
//	    WithTheme(renderer.DraculaTheme)
//
//	program := tea.NewProgram(model, tea.WithAltScreen())
//	if _, err := program.Run(); err != nil {
//	    log.Fatal(err)
//	}
//
// Table with custom columns and callbacks:
//
//	columns := []table.Column{
//	    *table.NewColumn("name", "Employee").WithWidth(25),
//	    *table.NewColumn("salary", "Salary").WithFormatter(table.CurrencyFormatter),
//	}
//
//	data := []map[string]interface{}{
//	    {"name": "Alice", "salary": 75000.0},
//	    {"name": "Bob", "salary": 65000.0},
//	}
//
//	model := components.NewTableWithColumns(data, columns).
//	    WithOnSelect(func(row map[string]interface{}) {
//	        fmt.Printf("Selected: %v\n", row)
//	    }).
//	    WithOnSort(func(column string, ascending bool) {
//	        fmt.Printf("Sorted by %s (asc: %t)\n", column, ascending)
//	    }).
//	    WithKeyBindings(components.VimKeyBindings())
//
// # Builder API
//
// The fluent builder API allows for easy configuration:
//
//	model := components.NewTable(data).
//	    WithTitle("Employee Directory").
//	    WithPageSize(15).
//	    WithSorting(true).
//	    WithSearch(true).
//	    WithTheme(renderer.MonokaiTheme).
//	    WithKeyBindings(components.VimKeyBindings()).
//	    WithOnSelect(handleSelection).
//	    WithOnSort(handleSort).
//	    WithOnSearch(handleSearch).
//	    WithOnRefresh(handleRefresh)
//
// All builder methods return the model for chaining, and most are optional
// with sensible defaults.
//
// # Key Bindings
//
// Three key binding presets are available:
//
// Default key bindings:
//   - ↑/↓: Navigate rows
//   - ←/→: Navigate pages
//   - Home/End: First/last page
//   - 1-9: Sort by column number
//   - /: Enter search mode
//   - +/-: Adjust page size
//   - ?: Toggle help
//   - q/Esc: Quit
//
// Vim key bindings (components.VimKeyBindings()):
//   - j/k: Navigate rows
//   - h/l: Navigate pages
//   - gg/G: First/last page
//   - 1-9: Sort by column number
//   - /: Enter search mode
//   - +/-: Adjust page size
//   - ?: Toggle help
//   - q/Esc: Quit
//
// Emacs key bindings (components.EmacsKeyBindings()):
//   - C-n/C-p: Navigate rows
//   - C-f/C-b: Navigate pages
//   - C-a/C-e: First/last page
//   - 1-9: Sort by column number
//   - C-s: Enter search mode
//   - +/-: Adjust page size
//   - C-h: Toggle help
//   - C-g/Esc: Quit
//
// Custom key bindings can be created by implementing the KeyBindings interface.
//
// # Callbacks
//
// The component system provides callbacks for user interactions:
//
//   - OnSelect: Called when a row is selected (Enter key)
//   - OnSort: Called when column sorting changes
//   - OnSearch: Called when search query changes
//   - OnRefresh: Called when refresh is requested (F5)
//   - OnPageChange: Called when page navigation occurs
//
// Callbacks receive relevant data and can be used to integrate with external
// systems or trigger additional actions.
//
// # State Management
//
// The Model maintains internal state including:
//   - Current page and selected row
//   - Search query and active search mode
//   - Sort column and direction
//   - Theme and key binding configuration
//   - Window dimensions and layout
//
// State is automatically managed and persists across renders and updates.
//
// # Search Mode
//
// Search functionality includes:
//   - Real-time filtering as user types
//   - Case-insensitive search across all searchable columns
//   - Visual indication of search mode
//   - Escape to exit search mode
//   - Automatic pagination adjustment for filtered results
//
// # Integration
//
// The components integrate seamlessly with Bubble Tea applications:
//
//	type AppModel struct {
//	    table components.Model
//	    // ... other components
//	}
//
//	func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	    var cmd tea.Cmd
//	    m.table, cmd = m.table.Update(msg)
//	    // ... handle other updates
//	    return m, cmd
//	}
//
//	func (m AppModel) View() string {
//	    return lipgloss.JoinVertical(
//	        lipgloss.Left,
//	        "My Application",
//	        m.table.View(),
//	        // ... other views
//	    )
//	}
//
// # Performance
//
// The component layer is optimized for interactive use:
//   - Efficient event handling and state updates
//   - Minimal re-rendering on state changes
//   - Responsive keyboard input processing
//   - Smooth pagination and scrolling
//   - Optimized search filtering
package components
