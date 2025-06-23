package renderer

import (
	"fmt"
	"strings"

	"github.com/anurag-roy/bubbletable/table"
)

// TableRenderer handles rendering tables with proper styling and layout
type TableRenderer struct {
	terminalWidth  int
	terminalHeight int
	theme          Theme
}

// NewTableRenderer creates a new table renderer with the default theme
func NewTableRenderer(width, height int) *TableRenderer {
	return &TableRenderer{
		terminalWidth:  width,
		terminalHeight: height,
		theme:          DefaultTheme,
	}
}

// NewTableRendererWithTheme creates a new table renderer with a custom theme
func NewTableRendererWithTheme(width, height int, theme Theme) *TableRenderer {
	return &TableRenderer{
		terminalWidth:  width,
		terminalHeight: height,
		theme:          theme,
	}
}

// SetTheme updates the renderer's theme
func (r *TableRenderer) SetTheme(theme Theme) {
	r.theme = theme
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
		headerText := col.Header

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

		return r.theme.Header.Width(col.Width).Render(finalHeader)
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
			var cellVal interface{}
			if colIndex < len(row.Cells) {
				cell := row.Cells[colIndex]
				cellVal = cell.Value
				// Use the column's formatter
				cellValue = col.Formatter(cell.Value)
			}

			content := r.truncateText(cellValue, col.Width)

			// Use custom renderer if available
			if col.Renderer != nil {
				content = col.Renderer(cellVal, isSelected)
				content = r.truncateText(content, col.Width)
			}

			if isSelected {
				return r.theme.SelectedRow.Width(col.Width).Render(content)
			}
			return r.theme.Cell.Width(col.Width).Render(content)
		})

		tableRows = append(tableRows, dataRow)
	}

	// Join all table content
	tableContent := strings.Join(tableRows, "\n")

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

// buildTableRow builds a table row using the provided cell renderer function
func (r *TableRenderer) buildTableRow(columns []table.Column, cellRenderer func(int, table.Column) string) string {
	var cells []string
	for i, col := range columns {
		cells = append(cells, cellRenderer(i, col))
	}
	return strings.Join(cells, "│")
}

// buildSeparatorRow builds a separator row between header and data
func (r *TableRenderer) buildSeparatorRow(columns []table.Column) string {
	var separators []string
	for _, col := range columns {
		separator := strings.Repeat("─", col.Width+2) // +2 for padding
		separators = append(separators, separator)
	}
	return strings.Join(separators, "┼")
}

// renderStatusLine renders the status line with page information
func (r *TableRenderer) renderStatusLine(tbl *table.Table, currentPage int) string {
	if tbl == nil {
		return ""
	}

	totalPages := tbl.GetTotalPages()
	startRow := currentPage*tbl.PageSize + 1
	endRow := startRow + len(tbl.GetPage(currentPage)) - 1

	status := fmt.Sprintf("Page %d of %d | Rows %d-%d of %d | Page Size: %d",
		currentPage+1, totalPages, startRow, endRow, tbl.TotalRows, tbl.PageSize)

	// Add sort information
	if tbl.SortBy >= 0 && tbl.SortBy < len(tbl.Columns) {
		sortDir := "asc"
		if tbl.SortDesc {
			sortDir = "desc"
		}
		status += fmt.Sprintf(" | Sorted by %s (%s)", tbl.Columns[tbl.SortBy].Header, sortDir)
	}

	return r.theme.Status.Render(status)
}

// truncateText truncates text to fit within the specified width
func (r *TableRenderer) truncateText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	if width <= 3 {
		return text[:width]
	}

	return text[:width-3] + "..."
}

// GetOptimalPageSize calculates the optimal page size based on terminal height
func (r *TableRenderer) GetOptimalPageSize() int {
	// Reserve space for header, separator, status, and some padding
	reservedLines := 10
	availableLines := r.terminalHeight - reservedLines

	if availableLines < 5 {
		return 5 // Minimum page size
	}

	return availableLines
}

// GetMaxTableHeight returns the maximum height for table content
func (r *TableRenderer) GetMaxTableHeight() int {
	return r.terminalHeight - 5 // Reserve space for status and padding
}

// GetTableCapacity returns how many rows can fit in the current terminal
func (r *TableRenderer) GetTableCapacity() int {
	return r.GetOptimalPageSize()
}

// CalculateColumnWidths calculates optimal column widths based on content
func (r *TableRenderer) CalculateColumnWidths(tbl *table.Table, maxSampleRows int) []table.Column {
	if tbl == nil || len(tbl.Columns) == 0 {
		return []table.Column{}
	}

	columns := make([]table.Column, len(tbl.Columns))
	copy(columns, tbl.Columns)

	// Sample up to maxSampleRows to determine optimal widths
	sampleSize := len(tbl.Rows)
	if sampleSize > maxSampleRows {
		sampleSize = maxSampleRows
	}

	// Calculate max width needed for each column
	for i, col := range columns {
		maxWidth := len(col.Header) // Start with header width

		// Check sample data
		for j := 0; j < sampleSize; j++ {
			if i < len(tbl.Rows[j].Cells) {
				cellValue := col.Formatter(tbl.Rows[j].Cells[i].Value)
				if len(cellValue) > maxWidth {
					maxWidth = len(cellValue)
				}
			}
		}

		// Set reasonable bounds
		if maxWidth < 5 {
			maxWidth = 5
		}
		if maxWidth > 50 {
			maxWidth = 50
		}

		columns[i].Width = maxWidth
	}

	return columns
}
