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
	keyBindings *KeyBindings
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

// NewTableFromInterface creates a new table model from a slice of interface{}
func NewTableFromInterface(data []interface{}) *TableModel {
	tbl := table.New()
	if len(data) > 0 {
		_ = tbl.SetData(data) // Ignore error for initialization
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

	_ = tbl.SetData(interfaceData) // Ignore error for initialization

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
func (m *TableModel) WithTheme(theme *renderer.Theme) *TableModel {
	m.theme = *theme
	if m.renderer != nil {
		m.renderer.SetTheme(theme)
	}
	return m
}

// WithKeyBindings sets custom key bindings
func (m *TableModel) WithKeyBindings(bindings *KeyBindings) *TableModel {
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

// handleKeyPress is split into smaller functions to reduce complexity
func (m *TableModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Handle search mode input
	if m.searchMode {
		return m.handleSearchInput(key)
	}

	// Handle help mode - only allow help and quit keys
	if m.showHelp {
		return m.handleHelpInput(key)
	}

	// Handle different key categories
	if handled, model := m.handleNavigationKeys(key); handled {
		return model, nil
	}

	if handled, model := m.handleActionKeys(key); handled {
		return model, nil
	}

	if handled, model := m.handleSortKeys(key); handled {
		return model, nil
	}

	// Trigger selection callback if we have one
	m.triggerSelectionCallback()

	return m, nil
}

// handleNavigationKeys handles navigation key presses
func (m *TableModel) handleNavigationKeys(key string) (bool, tea.Model) {
	if m.keyBindings.IsQuit(key) {
		return true, m // Will be handled by caller with tea.Quit
	}

	if handled, model := m.handleVerticalNavigation(key); handled {
		return true, model
	}

	if handled, model := m.handleHorizontalNavigation(key); handled {
		return true, model
	}

	if handled, model := m.handleJumpNavigation(key); handled {
		return true, model
	}

	return false, m
}

// handleVerticalNavigation handles up and down key presses
func (m *TableModel) handleVerticalNavigation(key string) (bool, tea.Model) {
	if m.keyBindings.IsUp(key) {
		if m.selectedRow > 0 {
			m.selectedRow--
		}
		return true, m
	}

	if m.keyBindings.IsDown(key) {
		currentTable := m.getCurrentTable()
		if currentTable != nil {
			pageData := currentTable.GetPage(m.currentPage)
			if m.selectedRow < len(pageData)-1 {
				m.selectedRow++
			}
		}
		return true, m
	}

	return false, m
}

// handleHorizontalNavigation handles left, right, page up, and page down key presses
func (m *TableModel) handleHorizontalNavigation(key string) (bool, tea.Model) {
	if m.keyBindings.IsLeft(key) || m.keyBindings.IsPageUp(key) {
		if m.currentPage > 0 {
			m.currentPage--
			m.selectedRow = 0
		}
		return true, m
	}

	if m.keyBindings.IsRight(key) || m.keyBindings.IsPageDown(key) {
		currentTable := m.getCurrentTable()
		if currentTable != nil && m.currentPage < currentTable.GetTotalPages()-1 {
			m.currentPage++
			m.selectedRow = 0
		}
		return true, m
	}

	return false, m
}

// handleJumpNavigation handles home and end key presses
func (m *TableModel) handleJumpNavigation(key string) (bool, tea.Model) {
	if m.keyBindings.IsHome(key) {
		m.currentPage = 0
		m.selectedRow = 0
		return true, m
	}

	if m.keyBindings.IsEnd(key) {
		currentTable := m.getCurrentTable()
		if currentTable != nil {
			m.currentPage = currentTable.GetTotalPages() - 1
			if m.currentPage < 0 {
				m.currentPage = 0
			}
			m.selectedRow = 0
		}
		return true, m
	}

	return false, m
}

// handleActionKeys handles action key presses
func (m *TableModel) handleActionKeys(key string) (bool, tea.Model) {
	switch {
	case m.keyBindings.IsSearch(key):
		m.searchMode = true
		return true, m

	case m.keyBindings.IsHelp(key):
		m.showHelp = !m.showHelp
		return true, m

	case m.keyBindings.IsRefresh(key):
		if m.onRefresh != nil {
			m.onRefresh()
		}
		return true, m

	case m.keyBindings.IsPageSizeUp(key):
		m.adjustPageSize(m.pageSize + 5)
		return true, m

	case m.keyBindings.IsPageSizeDown(key):
		if m.pageSize > 5 {
			m.adjustPageSize(m.pageSize - 5)
		}
		return true, m

	case m.keyBindings.IsResetPage(key):
		if m.renderer != nil {
			m.adjustPageSize(m.renderer.GetOptimalPageSize())
		}
		return true, m

	case m.keyBindings.IsClearSort(key):
		if m.table != nil {
			m.table.ClearSort()
			m.currentPage = 0
			m.selectedRow = 0
		}
		return true, m
	}

	return false, m
}

// handleSortKeys handles sorting key presses
func (m *TableModel) handleSortKeys(key string) (bool, tea.Model) {
	if colIndex := m.keyBindings.GetSortColumn(key); colIndex >= 0 && m.table != nil {
		if colIndex < len(m.table.Columns) {
			// Three-state sorting: unsorted -> asc -> desc -> unsorted
			if m.table.SortBy == colIndex {
				if !m.table.SortDesc {
					_ = m.table.SortByColumn(colIndex, true) // Ignore error
				} else {
					m.table.ClearSort()
				}
			} else {
				_ = m.table.SortByColumn(colIndex, false) // Ignore error
			}

			m.currentPage = 0
			m.selectedRow = 0

			if m.onSort != nil {
				m.onSort(colIndex, m.table.SortDesc)
			}
		}
		return true, m
	}

	return false, m
}

// triggerSelectionCallback triggers the selection callback if configured
func (m *TableModel) triggerSelectionCallback() {
	if m.onSelect != nil {
		currentTable := m.getCurrentTable()
		if currentTable != nil {
			pageData := currentTable.GetPage(m.currentPage)
			if m.selectedRow < len(pageData) {
				m.onSelect(pageData[m.selectedRow])
			}
		}
	}
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
		if m.searchTerm != "" {
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

// handleHelpInput handles input during help mode
func (m *TableModel) handleHelpInput(key string) (tea.Model, tea.Cmd) {
	switch {
	case m.keyBindings.IsHelp(key):
		// Toggle help off
		m.showHelp = false
		return m, nil
	case m.keyBindings.IsQuit(key):
		// Allow quitting from help
		return m, nil
	default:
		// Ignore all other keys in help mode
		return m, nil
	}
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

Actions:
  /           - Search
  ?           - Toggle help
  r           - Refresh
  +/=         - Increase page size
  -/_         - Decrease page size
  c           - Clear sort
  q/Esc       - Quit

Sorting:
  1-9         - Sort by column 1-9

Search Mode:
  Esc         - Exit search
  Backspace   - Delete character
  Enter       - Apply search
`

	return m.theme.Cell.Render(help)
}

// Public API methods

// GetTable returns the underlying table
func (m *TableModel) GetTable() *table.Table {
	return m.table
}

// GetCurrentTable returns the current effective table (filtered or main)
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

// SetData sets new data for the table
func (m *TableModel) SetData(data interface{}) error {
	if m.table == nil {
		return fmt.Errorf("table is not initialized")
	}

	err := m.table.SetData(data)
	if err != nil {
		return err
	}

	// Reset state
	m.currentPage = 0
	m.selectedRow = 0
	m.filteredTable = nil
	m.searchTerm = ""
	m.searchMode = false

	return nil
}

// RefreshData refreshes the table with new data
func (m *TableModel) RefreshData(data interface{}) error {
	return m.SetData(data)
}
