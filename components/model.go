package components

import (
	"fmt"
	"strings"

	"github.com/anurag-roy/bubbletable/renderer"
	"github.com/anurag-roy/bubbletable/table"
	tea "github.com/charmbracelet/bubbletea"
)

// TableModel represents the Bubble Tea model for the table component
type TableModel struct {
	table         *table.Table
	filteredTable *table.Table
	renderer      *renderer.TableRenderer

	// State
	ready       bool
	currentPage int
	selectedRow int
	searchMode  bool
	searchTerm  string

	// Configuration
	keyBindings KeyBindings
	theme       renderer.Theme
	pageSize    int
	showHelp    bool

	// Callbacks
	onSelect  func(row table.Row)
	onSort    func(columnIndex int, desc bool)
	onSearch  func(term string)
	onRefresh func()

	// Dimensions
	width  int
	height int
}

// NewTable creates a new table model from a slice of data
func NewTable[T any](data []T) *TableModel {
	// Convert slice to interface{} slice
	interfaceData := make([]interface{}, len(data))
	for i, item := range data {
		interfaceData[i] = item
	}

	return NewTableFromInterface(interfaceData)
}

// NewTableFromInterface creates a new table model from interface{} data
func NewTableFromInterface(data []interface{}) *TableModel {
	tbl := table.New()
	if len(data) > 0 {
		tbl.SetData(data)
	}

	return &TableModel{
		table:         tbl,
		filteredTable: nil,
		renderer:      renderer.NewTableRenderer(80, 24),
		keyBindings:   DefaultKeyBindings(),
		theme:         renderer.DefaultTheme,
		pageSize:      10,
		currentPage:   0,
		selectedRow:   0,
		searchMode:    false,
		searchTerm:    "",
		ready:         false,
		showHelp:      false,
	}
}

// NewTableWithColumns creates a new table model with predefined columns
func NewTableWithColumns(data []map[string]interface{}, columns []table.Column) *TableModel {
	tbl := table.NewWithColumns(columns)

	// Convert data to interface{} slice
	interfaceData := make([]interface{}, len(data))
	for i, item := range data {
		interfaceData[i] = item
	}

	tbl.SetData(interfaceData)

	return &TableModel{
		table:         tbl,
		filteredTable: nil,
		renderer:      renderer.NewTableRenderer(80, 24),
		keyBindings:   DefaultKeyBindings(),
		theme:         renderer.DefaultTheme,
		pageSize:      10,
		currentPage:   0,
		selectedRow:   0,
		searchMode:    false,
		searchTerm:    "",
		ready:         false,
		showHelp:      false,
	}
}

// Builder pattern methods for configuration

// WithPageSize sets the page size
func (m *TableModel) WithPageSize(size int) *TableModel {
	m.pageSize = size
	if m.table != nil {
		m.table.PageSize = size
	}
	return m
}

// WithTheme sets the theme
func (m *TableModel) WithTheme(theme renderer.Theme) *TableModel {
	m.theme = theme
	if m.renderer != nil {
		m.renderer.SetTheme(theme)
	}
	return m
}

// WithKeyBindings sets custom key bindings
func (m *TableModel) WithKeyBindings(bindings KeyBindings) *TableModel {
	m.keyBindings = bindings
	return m
}

// WithSorting enables or disables sorting
func (m *TableModel) WithSorting(enabled bool) *TableModel {
	if m.table != nil {
		for i := range m.table.Columns {
			m.table.Columns[i].Sortable = enabled
		}
	}
	return m
}

// WithSearch enables or disables search
func (m *TableModel) WithSearch(enabled bool) *TableModel {
	if m.table != nil {
		for i := range m.table.Columns {
			m.table.Columns[i].Searchable = enabled
		}
	}
	return m
}

// WithOnSelect sets a callback for row selection
func (m *TableModel) WithOnSelect(callback func(row table.Row)) *TableModel {
	m.onSelect = callback
	return m
}

// WithOnSort sets a callback for sorting changes
func (m *TableModel) WithOnSort(callback func(columnIndex int, desc bool)) *TableModel {
	m.onSort = callback
	return m
}

// WithOnSearch sets a callback for search changes
func (m *TableModel) WithOnSearch(callback func(term string)) *TableModel {
	m.onSearch = callback
	return m
}

// WithOnRefresh sets a callback for refresh requests
func (m *TableModel) WithOnRefresh(callback func()) *TableModel {
	m.onRefresh = callback
	return m
}

// Bubble Tea interface implementation

// Init initializes the model
func (m *TableModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		if m.renderer != nil {
			m.renderer.UpdateSize(msg.Width, msg.Height)
		}

		// Optimize page size for terminal
		if m.pageSize == 10 { // Only adjust if using default
			optimalSize := m.renderer.GetOptimalPageSize()
			m.pageSize = optimalSize
			if m.table != nil {
				m.table.PageSize = optimalSize
			}
		}

		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress handles key press events
func (m *TableModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Handle search mode input
	if m.searchMode {
		return m.handleSearchInput(key)
	}

	// Handle normal mode keys
	switch {
	case m.keyBindings.IsQuit(key):
		return m, tea.Quit

	case m.keyBindings.IsUp(key):
		if m.selectedRow > 0 {
			m.selectedRow--
		}

	case m.keyBindings.IsDown(key):
		currentTable := m.getCurrentTable()
		if currentTable != nil {
			pageData := currentTable.GetPage(m.currentPage)
			if m.selectedRow < len(pageData)-1 {
				m.selectedRow++
			}
		}

	case m.keyBindings.IsLeft(key) || m.keyBindings.IsPageUp(key):
		if m.currentPage > 0 {
			m.currentPage--
			m.selectedRow = 0
		}

	case m.keyBindings.IsRight(key) || m.keyBindings.IsPageDown(key):
		currentTable := m.getCurrentTable()
		if currentTable != nil && m.currentPage < currentTable.GetTotalPages()-1 {
			m.currentPage++
			m.selectedRow = 0
		}

	case m.keyBindings.IsHome(key):
		m.currentPage = 0
		m.selectedRow = 0

	case m.keyBindings.IsEnd(key):
		currentTable := m.getCurrentTable()
		if currentTable != nil {
			m.currentPage = currentTable.GetTotalPages() - 1
			m.selectedRow = 0
		}

	case m.keyBindings.IsSearch(key):
		m.searchMode = true
		m.searchTerm = ""

	case m.keyBindings.IsPageSizeUp(key):
		m.adjustPageSize(m.pageSize + 5)

	case m.keyBindings.IsPageSizeDown(key):
		newSize := m.pageSize - 5
		if newSize < 5 {
			newSize = 5
		}
		m.adjustPageSize(newSize)

	case m.keyBindings.IsResetPage(key):
		if m.renderer != nil {
			optimalSize := m.renderer.GetOptimalPageSize()
			m.adjustPageSize(optimalSize)
		}

	case m.keyBindings.IsClearSort(key):
		if m.table != nil {
			m.table.ClearSort()
			m.currentPage = 0
			m.selectedRow = 0
		}

	case m.keyBindings.IsHelp(key):
		m.showHelp = !m.showHelp

	case m.keyBindings.IsRefresh(key):
		if m.onRefresh != nil {
			m.onRefresh()
		}

	default:
		// Check for sort keys
		if colIndex := m.keyBindings.GetSortColumn(key); colIndex >= 0 && m.table != nil {
			if colIndex < len(m.table.Columns) {
				// Three-state sorting: unsorted -> asc -> desc -> unsorted
				if m.table.SortBy == colIndex {
					if !m.table.SortDesc {
						m.table.SortByColumn(colIndex, true)
					} else {
						m.table.ClearSort()
					}
				} else {
					m.table.SortByColumn(colIndex, false)
				}

				m.currentPage = 0
				m.selectedRow = 0

				if m.onSort != nil {
					m.onSort(colIndex, m.table.SortDesc)
				}
			}
		}
	}

	// Trigger selection callback if we have one
	if m.onSelect != nil {
		currentTable := m.getCurrentTable()
		if currentTable != nil {
			pageData := currentTable.GetPage(m.currentPage)
			if m.selectedRow < len(pageData) {
				m.onSelect(pageData[m.selectedRow])
			}
		}
	}

	return m, nil
}

// handleSearchInput handles input during search mode
func (m *TableModel) handleSearchInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		m.searchMode = false
		m.searchTerm = ""
		m.filteredTable = nil
		m.currentPage = 0
		m.selectedRow = 0

	case "backspace":
		if len(m.searchTerm) > 0 {
			m.searchTerm = m.searchTerm[:len(m.searchTerm)-1]
			m.updateSearch()
		}

	case "enter":
		m.searchMode = false

	default:
		// Handle character input
		if len(key) == 1 {
			char := key[0]
			if char >= 32 && char <= 126 { // Printable ASCII
				m.searchTerm += key
				m.updateSearch()
			}
		}
	}

	return m, nil
}

// updateSearch updates the filtered table based on search term
func (m *TableModel) updateSearch() {
	if m.table == nil {
		return
	}

	if m.searchTerm == "" {
		m.filteredTable = nil
	} else {
		m.filteredTable = m.table.Filter(m.searchTerm)
	}

	m.currentPage = 0
	m.selectedRow = 0

	if m.onSearch != nil {
		m.onSearch(m.searchTerm)
	}
}

// adjustPageSize adjusts the page size and recalculates pages
func (m *TableModel) adjustPageSize(newSize int) {
	m.pageSize = newSize
	if m.table != nil {
		m.table.PageSize = newSize
	}
	if m.filteredTable != nil {
		m.filteredTable.PageSize = newSize
	}

	// Reset to first page to avoid being out of bounds
	m.currentPage = 0
	m.selectedRow = 0
}

// getCurrentTable returns the current table (filtered or main)
func (m *TableModel) getCurrentTable() *table.Table {
	if m.filteredTable != nil {
		return m.filteredTable
	}
	return m.table
}

// View renders the table
func (m *TableModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	if m.showHelp {
		return m.renderHelp()
	}

	var content strings.Builder

	// Title
	title := m.theme.Header.Render("Data Table")
	content.WriteString(title + "\n\n")

	// Table content
	currentTable := m.getCurrentTable()
	if currentTable != nil && m.renderer != nil {
		tableContent := m.renderer.RenderTable(currentTable, m.currentPage, m.selectedRow)
		content.WriteString(tableContent)
	} else {
		content.WriteString("No data available")
	}

	content.WriteString("\n\n")

	// Status bar
	content.WriteString(m.renderStatusBar())

	// Search bar
	if m.searchMode {
		content.WriteString("\n")
		searchText := fmt.Sprintf("Search: %s", m.searchTerm)
		content.WriteString(m.theme.Search.Render(searchText))
	}

	return content.String()
}

// renderStatusBar renders the status bar with navigation info
func (m *TableModel) renderStatusBar() string {
	currentTable := m.getCurrentTable()
	if currentTable == nil {
		return ""
	}

	totalPages := currentTable.GetTotalPages()
	startRow := m.currentPage*m.pageSize + 1
	pageData := currentTable.GetPage(m.currentPage)
	endRow := startRow + len(pageData) - 1

	if len(pageData) == 0 {
		startRow = 0
		endRow = 0
	}

	status := fmt.Sprintf("Page %d/%d | Rows %d-%d of %d | Page Size: %d",
		m.currentPage+1, totalPages, startRow, endRow, currentTable.TotalRows, m.pageSize)

	// Add sort info
	if currentTable.SortBy >= 0 && currentTable.SortBy < len(currentTable.Columns) {
		sortDir := "↑"
		if currentTable.SortDesc {
			sortDir = "↓"
		}
		status += fmt.Sprintf(" | Sort: %s %s", currentTable.Columns[currentTable.SortBy].Header, sortDir)
	}

	// Add search info
	if m.searchTerm != "" {
		status += fmt.Sprintf(" | Search: '%s'", m.searchTerm)
	}

	return m.theme.Status.Render(status)
}

// renderHelp renders the help screen
func (m *TableModel) renderHelp() string {
	help := `
Table Navigation Help

Navigation:
  ↑/k         - Move up
  ↓/j         - Move down  
  ←/h         - Previous page
  →/l         - Next page
  Home/g      - First page
  End/G       - Last page

Sorting:
  1-9         - Sort by column (press again to reverse, third time to clear)
  c           - Clear sort

Search:
  /           - Enter search mode
  ESC         - Exit search mode
  Backspace   - Delete search characters

Page Size:
  +/=         - Increase page size
  -/_         - Decrease page size
  r           - Reset to optimal page size

Other:
  ?           - Toggle this help
  q/ESC/Ctrl+C - Quit

Press ? again to return to the table.
`

	return m.theme.Cell.Render(help)
}

// GetTable returns the underlying table (for advanced usage)
func (m *TableModel) GetTable() *table.Table {
	return m.table
}

// GetCurrentTable returns the currently displayed table (filtered or main)
func (m *TableModel) GetCurrentTable() *table.Table {
	return m.getCurrentTable()
}

// GetSelectedRow returns the currently selected row
func (m *TableModel) GetSelectedRow() (table.Row, bool) {
	currentTable := m.getCurrentTable()
	if currentTable == nil {
		return table.Row{}, false
	}

	pageData := currentTable.GetPage(m.currentPage)
	if m.selectedRow >= len(pageData) {
		return table.Row{}, false
	}

	return pageData[m.selectedRow], true
}

// SetData updates the table data
func (m *TableModel) SetData(data interface{}) error {
	if m.table == nil {
		return fmt.Errorf("table not initialized")
	}

	// Reset state
	m.currentPage = 0
	m.selectedRow = 0
	m.searchTerm = ""
	m.filteredTable = nil

	return m.table.SetData(data)
}

// RefreshData refreshes the table data (useful for live updates)
func (m *TableModel) RefreshData(data interface{}) error {
	// Preserve current state
	currentPage := m.currentPage
	selectedRow := m.selectedRow
	searchTerm := m.searchTerm

	// Update data
	err := m.SetData(data)
	if err != nil {
		return err
	}

	// Restore state if possible
	if m.table != nil {
		maxPages := m.table.GetTotalPages()
		if currentPage < maxPages {
			m.currentPage = currentPage
		}

		pageData := m.table.GetPage(m.currentPage)
		if selectedRow < len(pageData) {
			m.selectedRow = selectedRow
		}

		// Restore search if there was one
		if searchTerm != "" {
			m.searchTerm = searchTerm
			m.updateSearch()
		}
	}

	return nil
}
