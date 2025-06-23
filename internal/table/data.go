package table

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DataType represents the type of data in a column
type DataType int

const (
	String DataType = iota
	Integer
	Float
	Date
	Boolean
)

// Column represents a table column with metadata
type Column struct {
	Name     string
	Type     DataType
	Width    int // Display width
	Sortable bool
}

// Cell represents a single cell value with type information
type Cell struct {
	Value interface{}
	Type  DataType
}

// Row represents a table row
type Row struct {
	ID    int
	Cells []Cell
}

// Table represents the complete table data structure
type Table struct {
	Columns       []Column
	Rows          []Row
	UnsortedOrder []Row // Store original row order for unsort functionality
	SortBy        int   // Column index for sorting (-1 if not sorted)
	SortDesc      bool  // Sort direction (true for descending)
	PageSize      int
	TotalRows     int
}

// NewTable creates a new table with the given columns
func NewTable(columns []Column) *Table {
	return &Table{
		Columns:       columns,
		Rows:          make([]Row, 0),
		UnsortedOrder: make([]Row, 0),
		SortBy:        -1,
		SortDesc:      false,
		PageSize:      10,
		TotalRows:     0,
	}
}

// AddRow adds a new row to the table
func (t *Table) AddRow(values ...interface{}) error {
	if len(values) != len(t.Columns) {
		return fmt.Errorf("expected %d values, got %d", len(t.Columns), len(values))
	}

	cells := make([]Cell, len(values))
	for i, value := range values {
		cells[i] = Cell{
			Value: value,
			Type:  t.Columns[i].Type,
		}
	}

	row := Row{
		ID:    t.TotalRows,
		Cells: cells,
	}

	t.Rows = append(t.Rows, row)
	// Maintain original order for unsort functionality
	t.UnsortedOrder = append(t.UnsortedOrder, row)
	t.TotalRows++
	return nil
}

// GetPage returns a slice of rows for the given page number (0-indexed)
func (t *Table) GetPage(pageNum int) []Row {
	start := pageNum * t.PageSize
	end := start + t.PageSize

	if start >= len(t.Rows) {
		return []Row{}
	}

	if end > len(t.Rows) {
		end = len(t.Rows)
	}

	return t.Rows[start:end]
}

// GetTotalPages returns the total number of pages
func (t *Table) GetTotalPages() int {
	if t.PageSize <= 0 {
		return 1
	}
	if len(t.Rows) == 0 {
		return 1
	}
	return (len(t.Rows) + t.PageSize - 1) / t.PageSize
}

// SortByColumn sorts the table by the specified column
func (t *Table) SortByColumn(columnIndex int, descending bool) error {
	if columnIndex < 0 || columnIndex >= len(t.Columns) {
		return fmt.Errorf("invalid column index: %d", columnIndex)
	}

	if !t.Columns[columnIndex].Sortable {
		return fmt.Errorf("column %s is not sortable", t.Columns[columnIndex].Name)
	}

	t.SortBy = columnIndex
	t.SortDesc = descending

	sort.Slice(t.Rows, func(i, j int) bool {
		cellI := t.Rows[i].Cells[columnIndex]
		cellJ := t.Rows[j].Cells[columnIndex]

		result := compareCells(cellI, cellJ)
		if descending {
			return result > 0
		}
		return result < 0
	})

	return nil
}

// Filter returns a new table with rows matching the search term
func (t *Table) Filter(searchTerm string) *Table {
	if searchTerm == "" {
		return t
	}

	filtered := NewTable(t.Columns)
	filtered.PageSize = t.PageSize
	searchTerm = strings.ToLower(searchTerm)

	for _, row := range t.Rows {
		if t.rowMatchesSearch(row, searchTerm) {
			filtered.Rows = append(filtered.Rows, row)
			filtered.UnsortedOrder = append(filtered.UnsortedOrder, row)
			filtered.TotalRows++
		}
	}

	// Preserve sort state from original table
	filtered.SortBy = t.SortBy
	filtered.SortDesc = t.SortDesc

	return filtered
}

// rowMatchesSearch checks if a row contains the search term in any cell
func (t *Table) rowMatchesSearch(row Row, searchTerm string) bool {
	for _, cell := range row.Cells {
		cellStr := strings.ToLower(formatCellValue(cell))
		if strings.Contains(cellStr, searchTerm) {
			return true
		}
	}
	return false
}

// compareCells compares two cells for sorting purposes
func compareCells(a, b Cell) int {
	switch a.Type {
	case String:
		aStr := fmt.Sprintf("%v", a.Value)
		bStr := fmt.Sprintf("%v", b.Value)
		return strings.Compare(strings.ToLower(aStr), strings.ToLower(bStr))

	case Integer:
		aInt, _ := strconv.Atoi(fmt.Sprintf("%v", a.Value))
		bInt, _ := strconv.Atoi(fmt.Sprintf("%v", b.Value))
		return aInt - bInt

	case Float:
		aFloat, _ := strconv.ParseFloat(fmt.Sprintf("%v", a.Value), 64)
		bFloat, _ := strconv.ParseFloat(fmt.Sprintf("%v", b.Value), 64)
		if aFloat < bFloat {
			return -1
		} else if aFloat > bFloat {
			return 1
		}
		return 0

	case Date:
		aTime, aErr := parseDate(a.Value)
		bTime, bErr := parseDate(b.Value)
		if aErr != nil || bErr != nil {
			// Fallback to string comparison if date parsing fails
			return strings.Compare(fmt.Sprintf("%v", a.Value), fmt.Sprintf("%v", b.Value))
		}
		if aTime.Before(bTime) {
			return -1
		} else if aTime.After(bTime) {
			return 1
		}
		return 0

	case Boolean:
		aBool := fmt.Sprintf("%v", a.Value) == "true"
		bBool := fmt.Sprintf("%v", b.Value) == "true"
		if aBool == bBool {
			return 0
		} else if aBool {
			return 1
		}
		return -1

	default:
		return strings.Compare(fmt.Sprintf("%v", a.Value), fmt.Sprintf("%v", b.Value))
	}
}

// formatCellValue formats a cell value for display
func formatCellValue(cell Cell) string {
	switch cell.Type {
	case Date:
		if date, err := parseDate(cell.Value); err == nil {
			return date.Format("2006-01-02")
		}
		return fmt.Sprintf("%v", cell.Value)
	case Float:
		if f, ok := cell.Value.(float64); ok {
			return fmt.Sprintf("%.2f", f)
		}
		return fmt.Sprintf("%v", cell.Value)
	default:
		return fmt.Sprintf("%v", cell.Value)
	}
}

// parseDate attempts to parse various date formats
func parseDate(value interface{}) (time.Time, error) {
	str := fmt.Sprintf("%v", value)

	// Try common date formats
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"2006-01-02 15:04:05",
		"Jan 2, 2006",
		"January 2, 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", str)
}

// GetColumnNames returns a slice of column names
func (t *Table) GetColumnNames() []string {
	names := make([]string, len(t.Columns))
	for i, col := range t.Columns {
		names[i] = col.Name
	}
	return names
}

// GetCellValue returns the formatted value of a cell at the given row and column
func (t *Table) GetCellValue(rowIndex, columnIndex int) string {
	if rowIndex >= len(t.Rows) || columnIndex >= len(t.Rows[rowIndex].Cells) {
		return ""
	}
	return formatCellValue(t.Rows[rowIndex].Cells[columnIndex])
}

// ClearSort removes sorting and restores original order
func (t *Table) ClearSort() {
	if len(t.UnsortedOrder) > 0 {
		// Restore original order
		t.Rows = make([]Row, len(t.UnsortedOrder))
		copy(t.Rows, t.UnsortedOrder)
	}
	t.SortBy = -1
	t.SortDesc = false
}
