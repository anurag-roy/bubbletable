package main

import (
	"fmt"
	"os"
	"strings"

	"tui-data-table/internal/table"
	"tui-data-table/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type model struct {
	ready       bool
	table       *table.Table
	generator   *table.SampleDataGenerator
	renderer    *ui.TableRenderer
	currentPage int
	selectedRow int
	termWidth   int
	termHeight  int
	// Search functionality
	searchMode    bool
	searchTerm    string
	filteredTable *table.Table
}

// Define some basic styles
var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C4A9F4")).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1).
		Width(80) // Set a fixed width for the title
)

// Init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode input first
		if m.searchMode {
			switch msg.String() {
			case "esc":
				// Exit search mode
				m.searchMode = false
				m.searchTerm = ""
				m.filteredTable = nil
				m.currentPage = 0
				m.selectedRow = 0
				return m, nil
			case "backspace":
				if len(m.searchTerm) > 0 {
					m.searchTerm = m.searchTerm[:len(m.searchTerm)-1]
					m.updateSearch()
				}
				return m, nil
			case "up", "down", "left", "right", "home", "end", "pageup", "pagedown":
				// Allow arrow keys and special navigation keys in search mode - fall through to normal handling
			default:
				// Handle character input for search
				if len(msg.String()) == 1 {
					char := msg.String()[0]
					// Only allow printable characters for search
					if char >= 32 && char <= 126 {
						m.searchTerm += msg.String()
						m.updateSearch()
					}
				}
				return m, nil
			}
		}

		// Normal mode key handling
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		// Navigation keys
		case "up", "k":
			if m.selectedRow > 0 {
				m.selectedRow--
			}
		case "down", "j":
			currentTable := m.getCurrentTable()
			if currentTable != nil {
				pageData := currentTable.GetPage(m.currentPage)
				if m.selectedRow < len(pageData)-1 {
					m.selectedRow++
				}
			}

		// Pagination keys
		case "n", "right", "l":
			currentTable := m.getCurrentTable()
			if currentTable != nil && m.currentPage < currentTable.GetTotalPages()-1 {
				m.currentPage++
				m.selectedRow = 0 // Reset row selection
			}
		case "p", "left", "h":
			if m.currentPage > 0 {
				m.currentPage--
				m.selectedRow = 0 // Reset row selection
			}

		// Page navigation
		case "home", "g":
			m.currentPage = 0
			m.selectedRow = 0
		case "end", "G":
			if m.table != nil {
				m.currentPage = m.table.GetTotalPages() - 1
				m.selectedRow = 0
			}

		// Sorting keys (1-6 for first 6 columns)
		case "1", "2", "3", "4", "5", "6":
			if m.table != nil {
				colIndex := int(msg.String()[0] - '1')
				if colIndex < len(m.table.Columns) {
					// Three-click cycle: unsorted -> ascending -> descending -> unsorted
					if m.table.SortBy == colIndex {
						if !m.table.SortDesc {
							// Currently ascending, switch to descending
							m.table.SortByColumn(colIndex, true)
						} else {
							// Currently descending, clear sort (unsorted)
							m.table.ClearSort()
						}
					} else {
						// Different column or no sort, start with ascending
						m.table.SortByColumn(colIndex, false)
					}
					m.currentPage = 0 // Reset to first page after sorting
					m.selectedRow = 0
				}
			}

		// Add page size adjustment keys
		case "+", "=":
			if m.table != nil && m.renderer != nil {
				// Increase page size by 5, but cap at optimal size
				newSize := m.table.PageSize + 5
				optimalSize := m.renderer.GetOptimalPageSize()
				if newSize > optimalSize {
					newSize = optimalSize
				}
				m.adjustPageSize(newSize)
			}
		case "-", "_":
			if m.table != nil {
				// Decrease page size by 5, minimum of 5
				newSize := m.table.PageSize - 5
				if newSize < 5 {
					newSize = 5
				}
				m.adjustPageSize(newSize)
			}
		case "r":
			// Reset to optimal page size
			if m.table != nil && m.renderer != nil {
				optimalSize := m.renderer.GetOptimalPageSize()
				m.adjustPageSize(optimalSize)
			}

		// Search functionality
		case "/":
			m.searchMode = true
			m.searchTerm = ""
		}

	case tea.WindowSizeMsg:
		m.ready = true
		m.termWidth = msg.Width
		m.termHeight = msg.Height

		// Update or create renderer
		if m.renderer != nil {
			m.renderer.UpdateSize(msg.Width, msg.Height)
		} else {
			m.renderer = ui.NewTableRenderer(msg.Width, msg.Height)
		}

		// Optimize table layout for new size
		if m.table != nil && m.renderer != nil {
			// Calculate optimal column widths based on data
			optimizedColumns := m.renderer.CalculateColumnWidths(m.table, 100)
			m.table.Columns = optimizedColumns

			// Set optimal page size
			optimalPageSize := m.renderer.GetOptimalPageSize()
			m.adjustPageSize(optimalPageSize)
		}
	}
	return m, nil
}

// adjustPageSize safely adjusts the page size and maintains current position
func (m *model) adjustPageSize(newSize int) {
	if m.table == nil {
		return
	}

	// Calculate current absolute position
	currentRowGlobal := m.currentPage*m.table.PageSize + m.selectedRow

	// Update page size
	m.table.PageSize = newSize

	// Recalculate page and row position
	if currentRowGlobal >= len(m.table.Rows) {
		currentRowGlobal = len(m.table.Rows) - 1
	}

	if currentRowGlobal >= 0 && newSize > 0 {
		m.currentPage = currentRowGlobal / newSize
		m.selectedRow = currentRowGlobal % newSize
	} else {
		m.currentPage = 0
		m.selectedRow = 0
	}

	// Ensure valid page
	totalPages := m.table.GetTotalPages()
	if m.currentPage >= totalPages {
		m.currentPage = totalPages - 1
		if m.currentPage < 0 {
			m.currentPage = 0
		}
	}
}

// updateSearch applies the current search term and updates the filtered table
func (m *model) updateSearch() {
	if m.table == nil {
		return
	}

	if m.searchTerm == "" {
		// Clear search - use original table
		m.filteredTable = nil
		m.currentPage = 0
		m.selectedRow = 0
	} else {
		// Apply search filter
		m.filteredTable = m.table.Filter(m.searchTerm)
		m.currentPage = 0
		m.selectedRow = 0
	}
}

// getCurrentTable returns the table to display (filtered or original)
func (m model) getCurrentTable() *table.Table {
	if m.filteredTable != nil {
		return m.filteredTable
	}
	return m.table
}

// createFixedHeightContent ensures content fits in a fixed height container
func (m model) createFixedHeightContent(content string, targetHeight int) string {
	lines := strings.Split(content, "\n")

	// If content is taller than target, truncate it
	if len(lines) > targetHeight {
		lines = lines[:targetHeight]
		content = strings.Join(lines, "\n")
	}

	// If content is shorter than target, pad with empty lines
	if len(lines) < targetHeight {
		padding := targetHeight - len(lines)
		for i := 0; i < padding; i++ {
			content += "\n"
		}
	}

	return content
}

// createBottomStatusBar creates a combined status and help bar with flex-like layout
func (m model) createBottomStatusBar() string {
	// Left side: Status information
	var leftStatus string
	currentTable := m.getCurrentTable()
	if currentTable != nil {
		totalPages := currentTable.GetTotalPages()
		totalRows := len(currentTable.Rows)

		leftStatus = fmt.Sprintf("Page %d of %d │ %d rows │ Page size: %d",
			m.currentPage+1, totalPages, totalRows, currentTable.PageSize)

		// Add sort information if table is sorted
		if currentTable.SortBy >= 0 && currentTable.SortBy < len(currentTable.Columns) {
			sortDir := "↑"
			if currentTable.SortDesc {
				sortDir = "↓"
			}
			leftStatus += fmt.Sprintf(" │ Sorted by: %s %s",
				currentTable.Columns[currentTable.SortBy].Name, sortDir)
		}

		// Add search filter info
		if m.filteredTable != nil {
			leftStatus += fmt.Sprintf(" │ Filtered: %d of %d", len(m.filteredTable.Rows), len(m.table.Rows))
		}
	} else {
		leftStatus = "Loading..."
	}

	// Right side: Help/controls information
	var rightHelp string
	if m.searchMode {
		rightHelp = fmt.Sprintf("🔍 Search: %s (ESC to exit, Backspace to delete)", m.searchTerm)
	} else {
		rightHelp = " Navigation: ↑↓/jk │ Page: ←→/hl │ Sort: 1-6 │ Search: / │ Quit: q "
	}

	// Calculate available width and create flex-like layout
	availableWidth := m.termWidth
	if availableWidth < 50 {
		availableWidth = 50 // Minimum reasonable width
	}

	// Style both sides
	leftStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#C4A9F4"))
	rightStyle := lipgloss.NewStyle().Background(lipgloss.Color("#88F397")).Foreground(lipgloss.Color("#282A36"))

	// Calculate spacing to push right content to the right
	leftContent := leftStyle.Render(leftStatus)
	rightContent := rightStyle.Render(rightHelp)

	// Get actual rendered widths (without ANSI codes)
	leftWidth := lipgloss.Width(leftContent)
	rightWidth := lipgloss.Width(rightContent)

	// Calculate spacing needed
	spacingNeeded := availableWidth - leftWidth - rightWidth
	if spacingNeeded < 1 {
		spacingNeeded = 1 // Minimum spacing
	}

	spacing := strings.Repeat(" ", spacingNeeded)

	return leftContent + spacing + rightContent
}

// View renders the UI
func (m model) View() string {
	if !m.ready {
		return "Initializing TUI Data Table..."
	}

	// Title - make it responsive to terminal width
	titleText := "TUI Data Table"
	if m.termWidth > 0 {
		titleStyle = titleStyle.Width(m.termWidth - 2)
	}
	title := titleStyle.Render(titleText)

	// Calculate fixed content area height
	// Reserve space for title (2 lines), bottom status bar (1 line), and some padding
	reservedHeight := 4
	contentHeight := m.termHeight - reservedHeight
	if contentHeight < 10 {
		contentHeight = 10 // Minimum reasonable height
	}

	// Table content in fixed-height container
	var content string
	currentTable := m.getCurrentTable()
	if currentTable != nil && m.renderer != nil {
		tableContent := m.renderer.RenderTable(currentTable, m.currentPage, m.selectedRow)

		// Ensure content fits exactly in the allocated height
		content = m.createFixedHeightContent(tableContent, contentHeight)
	} else {
		// Loading placeholder in fixed container
		loadingBox := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(2, 4).
			Render("Loading table data...")
		content = m.createFixedHeightContent(loadingBox, contentHeight)
	}

	// Create combined status/help bottom bar
	bottomBar := m.createBottomStatusBar()

	// Combine all elements
	return lipgloss.JoinVertical(lipgloss.Left, title, content, bottomBar)
}

func main() {
	// Initialize the model with sample data
	generator := table.NewSampleDataGenerator()
	sampleTable := generator.GenerateEmployeeTable()

	// Start with a reasonable default page size
	sampleTable.PageSize = 15

	m := model{
		table:       sampleTable,
		generator:   generator,
		currentPage: 0,
		selectedRow: 0,
	}

	// Create the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
