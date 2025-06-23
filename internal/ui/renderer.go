package ui

import (
	"fmt"
	"math"
	"strings"

	"tui-data-table/internal/table"

	"github.com/charmbracelet/lipgloss"
)

// TableRenderer handles rendering tables with proper styling and layout
type TableRenderer struct {
	terminalWidth  int
	terminalHeight int
	styles         TableStyles
}

// TableStyles contains all styling for table components
type TableStyles struct {
	Header      lipgloss.Style
	Cell        lipgloss.Style
	SelectedRow lipgloss.Style
	Border      lipgloss.Style
	Status      lipgloss.Style
}

// NewTableRenderer creates a new table renderer with default styles
func NewTableRenderer(width, height int) *TableRenderer {
	return &TableRenderer{
		terminalWidth:  width,
		terminalHeight: height,
		styles:         getDefaultTableStyles(),
	}
}

// getDefaultTableStyles returns the default styling for table components
func getDefaultTableStyles() TableStyles {
	return TableStyles{
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282A36")).
			Background(lipgloss.Color("#C4A9F4")).
			Bold(true).
			Padding(0, 1),

		Cell: lipgloss.NewStyle().
			Padding(0, 1),

		SelectedRow: lipgloss.NewStyle().
			Background(lipgloss.Color("#44475A")).
			Foreground(lipgloss.Color("#F8F8F2")).
			Padding(0, 1),

		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#874BFD")),

		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4")).
			Italic(true),
	}
}

// UpdateSize updates the terminal dimensions
func (r *TableRenderer) UpdateSize(width, height int) {
	r.terminalWidth = width
	r.terminalHeight = height
}

// RenderTable renders a complete table with headers, data, and status
func (r *TableRenderer) RenderTable(tbl *table.Table, currentPage int, selectedRow int) string {
	if tbl == nil {
		return "No table data available"
	}

	// Use nearly full terminal width - only subtract 2 for minimal border
	availableWidth := r.terminalWidth - 2
	if availableWidth < 30 {
		availableWidth = 30 // Minimum reasonable width
	}

	// Adjust column widths to use full available space
	adjustedColumns := r.distributeColumnWidths(tbl.Columns, availableWidth)

	// Build table components
	var tableRows []string

	// Header row with sort indicators
	headerRow := r.buildTableRow(adjustedColumns, func(i int, col table.Column) string {
		headerText := col.Name

		// Add sort indicator if this column is currently sorted
		sortIndicator := ""
		if tbl.SortBy == i {
			if tbl.SortDesc {
				sortIndicator = "↓"
			} else {
				sortIndicator = "↑"
			}
		}

		// Calculate available space for header text
		availableWidth := col.Width
		if sortIndicator != "" {
			availableWidth -= 2 // Reserve space for " " + arrow
		}

		// Truncate header text if needed to make room for sort indicator
		truncatedHeader := headerText
		if len(headerText) > availableWidth {
			truncatedHeader = r.truncateText(headerText, availableWidth)
		}

		// Combine header text and sort indicator, right-aligned
		finalHeader := truncatedHeader
		if sortIndicator != "" {
			padding := col.Width - len(truncatedHeader) - len(sortIndicator)
			if padding > 0 {
				finalHeader += strings.Repeat(" ", padding) + sortIndicator
			} else {
				finalHeader += sortIndicator
			}
		}

		return r.styles.Header.Width(col.Width).Render(finalHeader)
	})
	tableRows = append(tableRows, headerRow)

	// Separator row
	separatorRow := r.buildSeparatorRow(adjustedColumns)
	tableRows = append(tableRows, separatorRow)

	// Data rows
	pageData := tbl.GetPage(currentPage)
	for rowIndex, row := range pageData {
		isSelected := rowIndex == selectedRow

		dataRow := r.buildTableRow(adjustedColumns, func(colIndex int, col table.Column) string {
			cellValue := ""
			if colIndex < len(row.Cells) {
				// Use the cell value directly from the row
				cell := row.Cells[colIndex]
				// Format the cell value - simple formatting for now
				if cell.Type == table.Float {
					if f, ok := cell.Value.(float64); ok {
						cellValue = fmt.Sprintf("%.2f", f)
					} else {
						cellValue = fmt.Sprintf("%v", cell.Value)
					}
				} else {
					cellValue = fmt.Sprintf("%v", cell.Value)
				}
			}

			content := r.truncateText(cellValue, col.Width)

			if isSelected {
				return r.styles.SelectedRow.Width(col.Width).Render(content)
			}
			return r.styles.Cell.Width(col.Width).Render(content)
		})

		tableRows = append(tableRows, dataRow)
	}

	// Join all table content
	tableContent := strings.Join(tableRows, "\n")

	// Return just the table content without status line
	// (status is now handled by the main application's bottom bar)
	return tableContent
}

// distributeColumnWidths distributes available width across columns intelligently
func (r *TableRenderer) distributeColumnWidths(columns []table.Column, availableWidth int) []table.Column {
	if len(columns) == 0 {
		return columns
	}

	adjusted := make([]table.Column, len(columns))
	copy(adjusted, columns)

	// Account for separators between columns
	separatorWidth := len(columns) - 1
	contentWidth := availableWidth - separatorWidth

	// Ensure minimum content width
	if contentWidth < len(columns)*5 {
		contentWidth = len(columns) * 5
	}

	// Calculate ideal width per column
	baseWidth := contentWidth / len(columns)
	remainder := contentWidth % len(columns)

	// Distribute width, giving extra to first few columns
	for i := range adjusted {
		adjusted[i].Width = baseWidth
		if i < remainder {
			adjusted[i].Width++
		}

		// Ensure minimum width
		if adjusted[i].Width < 5 {
			adjusted[i].Width = 5
		}
	}

	return adjusted
}

// buildTableRow builds a single table row using a cell renderer function
func (r *TableRenderer) buildTableRow(columns []table.Column, cellRenderer func(int, table.Column) string) string {
	var cells []string
	for i, col := range columns {
		cell := cellRenderer(i, col)
		cells = append(cells, cell)
	}
	return strings.Join(cells, "│")
}

// buildSeparatorRow builds the separator line between header and data
func (r *TableRenderer) buildSeparatorRow(columns []table.Column) string {
	var parts []string
	for _, col := range columns {
		// Create separator that matches the actual cell width (no extra padding needed here)
		separator := strings.Repeat("─", col.Width)
		parts = append(parts, separator)
	}
	return strings.Join(parts, "┼")
}

// renderStatusLine renders the status/info line below the table
func (r *TableRenderer) renderStatusLine(tbl *table.Table, currentPage int) string {
	totalPages := tbl.GetTotalPages()
	totalRows := len(tbl.Rows)

	statusText := fmt.Sprintf("Page %d of %d │ %d rows │ Page size: %d",
		currentPage+1, totalPages, totalRows, tbl.PageSize)

	// Add sort information if table is sorted
	if tbl.SortBy >= 0 && tbl.SortBy < len(tbl.Columns) {
		sortDir := "↑"
		if tbl.SortDesc {
			sortDir = "↓"
		}
		statusText += fmt.Sprintf(" │ Sorted by: %s %s",
			tbl.Columns[tbl.SortBy].Name, sortDir)
	}

	return r.styles.Status.Render(statusText)
}

// truncateText truncates text to fit within the specified width
func (r *TableRenderer) truncateText(text string, width int) string {
	if width <= 0 {
		return ""
	}

	if len(text) <= width {
		return text
	}

	if width <= 3 {
		return strings.Repeat(".", width)
	}

	return text[:width-3] + "..."
}

// GetOptimalPageSize calculates the ideal number of rows based on terminal height
func (r *TableRenderer) GetOptimalPageSize() int {
	// Reserve space for: title (2), help text (2), status (2), header (1), separator (1), borders (2)
	reservedLines := 10

	availableLines := r.terminalHeight - reservedLines
	if availableLines < 5 {
		return 5 // Minimum reasonable page size
	}

	// Cap at reasonable maximum to avoid performance issues
	maxPageSize := 50
	if availableLines > maxPageSize {
		return maxPageSize
	}

	return availableLines
}

// GetMaxTableHeight returns the maximum height available for table content
func (r *TableRenderer) GetMaxTableHeight() int {
	return r.GetOptimalPageSize() + 2 // Add header and separator back
}

// GetTableCapacity returns how many rows can fit in the available space
func (r *TableRenderer) GetTableCapacity() int {
	return r.GetOptimalPageSize()
}

// CalculateColumnWidths analyzes table data to suggest optimal column widths
func (r *TableRenderer) CalculateColumnWidths(tbl *table.Table, maxSampleRows int) []table.Column {
	if tbl == nil || len(tbl.Columns) == 0 {
		return tbl.Columns
	}

	adjusted := make([]table.Column, len(tbl.Columns))
	copy(adjusted, tbl.Columns)

	// Analyze actual data to determine optimal widths
	sampleSize := int(math.Min(float64(len(tbl.Rows)), float64(maxSampleRows)))

	for colIndex, col := range adjusted {
		maxWidth := len(col.Name) // Start with header width

		// Sample data to find reasonable column width
		for i := 0; i < sampleSize; i++ {
			if i < len(tbl.Rows) {
				cellValue := tbl.GetCellValue(tbl.Rows[i].ID, colIndex)
				if len(cellValue) > maxWidth {
					maxWidth = len(cellValue)
				}
			}
		}

		// Set reasonable bounds
		if maxWidth < 8 {
			maxWidth = 8
		}
		if maxWidth > 25 {
			maxWidth = 25
		}

		adjusted[colIndex].Width = maxWidth
	}

	return adjusted
}
