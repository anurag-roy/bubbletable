package ui

import (
	"fmt"
	"strings"
	"testing"
	"tui-data-table/internal/table"
)

func TestNewTableRenderer(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	if renderer == nil {
		t.Fatal("Expected non-nil renderer")
	}

	if renderer.terminalWidth != 80 {
		t.Errorf("Expected width 80, got %d", renderer.terminalWidth)
	}

	if renderer.terminalHeight != 24 {
		t.Errorf("Expected height 24, got %d", renderer.terminalHeight)
	}
}

func TestUpdateSize(t *testing.T) {
	renderer := NewTableRenderer(80, 24)
	renderer.UpdateSize(120, 30)

	if renderer.terminalWidth != 120 {
		t.Errorf("Expected width 120, got %d", renderer.terminalWidth)
	}

	if renderer.terminalHeight != 30 {
		t.Errorf("Expected height 30, got %d", renderer.terminalHeight)
	}
}

func TestGetOptimalPageSize(t *testing.T) {
	renderer := NewTableRenderer(80, 24)
	pageSize := renderer.GetOptimalPageSize()

	// Should be terminal height minus reserved space (10 lines)
	expected := 24 - 10
	if pageSize != expected {
		t.Errorf("Expected page size %d, got %d", expected, pageSize)
	}
}

func TestGetOptimalPageSizeWithSmallTerminal(t *testing.T) {
	renderer := NewTableRenderer(80, 10)
	pageSize := renderer.GetOptimalPageSize()

	// Should return minimum of 5 for very small terminals
	if pageSize != 5 {
		t.Errorf("Expected minimum page size 5, got %d", pageSize)
	}
}

func TestGetOptimalPageSizeWithLargeTerminal(t *testing.T) {
	renderer := NewTableRenderer(80, 100)
	pageSize := renderer.GetOptimalPageSize()

	// Should cap at maximum of 50
	if pageSize != 50 {
		t.Errorf("Expected capped page size 50, got %d", pageSize)
	}
}

func TestTruncateText(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	// Text shorter than width
	result := renderer.truncateText("hello", 10)
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}

	// Text longer than width
	result = renderer.truncateText("this is a very long text", 10)
	if result != "this is..." {
		t.Errorf("Expected 'this is...', got '%s'", result)
	}

	// Very small width
	result = renderer.truncateText("hello", 3)
	if result != "..." {
		t.Errorf("Expected '...', got '%s'", result)
	}

	// Zero width
	result = renderer.truncateText("hello", 0)
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestDistributeColumnWidths(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	columns := []table.Column{
		{Name: "Col1", Width: 20},
		{Name: "Col2", Width: 30},
		{Name: "Col3", Width: 40},
	}

	// Test with specific available width
	adjusted := renderer.distributeColumnWidths(columns, 60)
	if len(adjusted) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(adjusted))
	}

	// Total width should approximately match available width (minus separators)
	totalWidth := 0
	for _, col := range adjusted {
		totalWidth += col.Width
	}

	// Should be close to available width minus separators (2)
	expectedTotal := 60 - 2
	if totalWidth != expectedTotal {
		t.Errorf("Expected total width %d, got %d", expectedTotal, totalWidth)
	}

	// All columns should have minimum width
	for i, col := range adjusted {
		if col.Width < 5 {
			t.Errorf("Column %d width too small: got %d, minimum is 5", i, col.Width)
		}
	}
}

func TestCalculateColumnWidths(t *testing.T) {
	renderer := NewTableRenderer(100, 24)

	// Create a test table with data
	columns := []table.Column{
		{Name: "ID", Type: table.Integer, Width: 5, Sortable: true},
		{Name: "Name", Type: table.String, Width: 15, Sortable: true},
		{Name: "Department", Type: table.String, Width: 15, Sortable: true},
	}

	testTable := table.NewTable(columns)
	testTable.AddRow(1, "Alice Smith", "Engineering")
	testTable.AddRow(2, "Bob Johnson with a very long name", "Sales")
	testTable.AddRow(3, "Charlie Brown", "Human Resources Department")

	// Test column width calculation
	optimized := renderer.CalculateColumnWidths(testTable, 10)

	if len(optimized) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(optimized))
	}

	// Name column should be wider due to long content
	if optimized[1].Width < 15 {
		t.Errorf("Name column should be at least 15 wide, got %d", optimized[1].Width)
	}

	// Department column should be wider due to long content
	if optimized[2].Width < 15 {
		t.Errorf("Department column should be at least 15 wide, got %d", optimized[2].Width)
	}

	// All columns should have reasonable bounds
	for i, col := range optimized {
		if col.Width < 8 {
			t.Errorf("Column %d width too small: got %d, minimum should be 8", i, col.Width)
		}
		if col.Width > 25 {
			t.Errorf("Column %d width too large: got %d, maximum should be 25", i, col.Width)
		}
	}
}

func TestRenderTable(t *testing.T) {
	renderer := NewTableRenderer(100, 24)

	// Create a test table
	columns := []table.Column{
		{Name: "ID", Type: table.Integer, Width: 5, Sortable: true},
		{Name: "Name", Type: table.String, Width: 15, Sortable: true},
	}

	testTable := table.NewTable(columns)
	testTable.AddRow(1, "Alice")
	testTable.AddRow(2, "Bob")
	testTable.AddRow(3, "Charlie")

	// Test rendering
	result := renderer.RenderTable(testTable, 0, 1)

	if result == "" {
		t.Error("Expected non-empty table rendering")
	}

	// Check that result contains expected elements
	if !strings.Contains(result, "ID") {
		t.Error("Expected table to contain 'ID' header")
	}

	if !strings.Contains(result, "Name") {
		t.Error("Expected table to contain 'Name' header")
	}

	if !strings.Contains(result, "Alice") {
		t.Error("Expected table to contain 'Alice' data")
	}

	if !strings.Contains(result, "Page 1 of") {
		t.Error("Expected table to contain page information")
	}
}

func TestBuildTableRow(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	columns := []table.Column{
		{Name: "Col1", Width: 10},
		{Name: "Col2", Width: 15},
	}

	result := renderer.buildTableRow(columns, func(i int, col table.Column) string {
		return fmt.Sprintf("Cell%d", i)
	})

	expected := "Cell0│Cell1"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestRenderStatusLine(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	// Create a test table
	columns := []table.Column{
		{Name: "ID", Type: table.Integer, Width: 5, Sortable: true},
		{Name: "Name", Type: table.String, Width: 15, Sortable: true},
	}

	testTable := table.NewTable(columns)
	testTable.AddRow(1, "Alice")
	testTable.AddRow(2, "Bob")

	// Test status line without sorting
	status := renderer.renderStatusLine(testTable, 0)

	if !strings.Contains(status, "Page 1 of 1") {
		t.Error("Expected status to contain page information")
	}

	if !strings.Contains(status, "2 rows") {
		t.Error("Expected status to contain row count")
	}

	// Test status line with sorting
	testTable.SortByColumn(0, false)
	status = renderer.renderStatusLine(testTable, 0)

	if !strings.Contains(status, "Sorted by: ID") {
		t.Error("Expected status to contain sort information")
	}

	if !strings.Contains(status, "↑") {
		t.Error("Expected status to contain ascending sort indicator")
	}
}

func TestRenderTableWithNilTable(t *testing.T) {
	renderer := NewTableRenderer(80, 24)
	result := renderer.RenderTable(nil, 0, 0)

	if result != "No table data available" {
		t.Errorf("Expected 'No table data available', got '%s'", result)
	}
}
