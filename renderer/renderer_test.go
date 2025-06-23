package renderer

import (
	"strings"
	"testing"

	"github.com/anurag-roy/bubbletable/table"
)

// TestNewTableRenderer tests renderer creation
func TestNewTableRenderer(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	if renderer == nil {
		t.Fatal("Expected renderer to be created")
	}

	if renderer.terminalWidth != 80 {
		t.Errorf("Expected width 80, got %d", renderer.terminalWidth)
	}

	if renderer.terminalHeight != 24 {
		t.Errorf("Expected height 24, got %d", renderer.terminalHeight)
	}
}

// TestNewTableRendererWithTheme tests renderer creation with theme
func TestNewTableRendererWithTheme(t *testing.T) {
	theme := DraculaTheme
	renderer := NewTableRendererWithTheme(100, 30, &theme)

	if renderer == nil {
		t.Fatal("Expected renderer to be created")
	}

	if renderer.theme.Name != "Dracula" {
		t.Errorf("Expected theme name 'Dracula', got '%s'", renderer.theme.Name)
	}
}

// TestSetTheme tests theme setting
func TestSetTheme(t *testing.T) {
	renderer := NewTableRenderer(80, 24)
	theme := MonokaiTheme
	renderer.SetTheme(&theme)

	if renderer.theme.Name != "Monokai" {
		t.Errorf("Expected theme name 'Monokai', got '%s'", renderer.theme.Name)
	}
}

// TestUpdateSize tests size updating
func TestUpdateSize(t *testing.T) {
	renderer := NewTableRenderer(80, 24)
	renderer.UpdateSize(100, 30)

	if renderer.terminalWidth != 100 {
		t.Errorf("Expected width 100, got %d", renderer.terminalWidth)
	}

	if renderer.terminalHeight != 30 {
		t.Errorf("Expected height 30, got %d", renderer.terminalHeight)
	}
}

// TestRenderTable tests table rendering
func TestRenderTable(t *testing.T) {
	// Create a simple table
	tbl := table.New()
	columns := []table.Column{
		{Header: "Name", Key: "name", Width: 10, Formatter: table.DefaultFormatter},
		{Header: "Age", Key: "age", Width: 5, Formatter: table.DefaultFormatter},
	}
	tbl.Columns = columns

	// Add test data
	err := tbl.SetData([]map[string]interface{}{
		{"name": "Alice", "age": 30},
		{"name": "Bob", "age": 25},
	})

	if err != nil {
		t.Fatalf("Failed to set data: %v", err)
	}

	renderer := NewTableRenderer(80, 24)
	output := renderer.RenderTable(tbl, 0, 0)

	if output == "" {
		t.Error("Expected non-empty output")
	}

	// Basic checks
	if len(output) < 10 {
		t.Error("Output seems too short")
	}
}

// TestRenderEmptyTable tests rendering empty table
func TestRenderEmptyTable(t *testing.T) {
	renderer := NewTableRenderer(80, 24)
	output := renderer.RenderTable(nil, 0, 0)

	expected := "No data to display"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

// TestGetOptimalPageSize tests page size calculation
func TestGetOptimalPageSize(t *testing.T) {
	renderer := NewTableRenderer(80, 24)
	pageSize := renderer.GetOptimalPageSize()

	if pageSize < 5 {
		t.Errorf("Page size too small: %d", pageSize)
	}

	// Test with small terminal
	renderer.UpdateSize(80, 10)
	pageSize = renderer.GetOptimalPageSize()

	if pageSize != 5 {
		t.Errorf("Expected minimum page size 5, got %d", pageSize)
	}
}

// TestDistributeColumnWidths tests column width distribution
func TestDistributeColumnWidths(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	columns := []table.Column{
		{Header: "Name", Width: 10},
		{Header: "Age", Width: 5},
		{Header: "City", Width: 15},
	}

	adjusted := renderer.distributeColumnWidths(columns, 60)

	// Check total width allocation
	totalWidth := 0
	for _, col := range adjusted {
		totalWidth += col.Width
		if col.Width < 5 {
			t.Errorf("Column width too small: %d", col.Width)
		}
	}

	// Should use available space (minus separators)
	expectedContent := 60 - (len(columns) - 1)
	if totalWidth != expectedContent {
		t.Errorf("Expected total width %d, got %d", expectedContent, totalWidth)
	}
}

// TestTruncateText tests text truncation
func TestTruncateText(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	tests := []struct {
		input    string
		width    int
		expected string
	}{
		{"Hello", 10, "Hello"},
		{"Hello World", 5, "He..."},
		{"Hi", 3, "Hi"},
		{"Test", 2, "Te"},
		{"A", 1, "A"},
	}

	for _, test := range tests {
		result := renderer.truncateText(test.input, test.width)
		if result != test.expected {
			t.Errorf("Input: %s, Width: %d, Expected: %s, Got: %s",
				test.input, test.width, test.expected, result)
		}
	}
}

func TestRenderTableWithSorting(t *testing.T) {
	type TestData struct {
		ID   int    `table:"ID,sortable,width:5"`
		Name string `table:"Name,sortable,width:15"`
	}

	data := []TestData{
		{2, "Bob"},
		{1, "Alice"},
	}

	tbl := table.New()
	if err := tbl.SetData(data); err != nil {
		t.Fatalf("Failed to set data: %v", err)
	}

	if err := tbl.SortByColumn(0, false); err != nil {
		t.Fatalf("Failed to sort: %v", err)
	}

	renderer := NewTableRenderer(50, 20)
	result := renderer.RenderTable(tbl, 0, 0)

	// Just verify the table renders without error and contains data
	if len(result) == 0 {
		t.Error("Rendered table should not be empty")
	}

	// Verify data is present
	if !strings.Contains(result, "Alice") || !strings.Contains(result, "Bob") {
		t.Error("Rendered table should contain the test data")
	}
}

func TestMinimumColumnWidth(t *testing.T) {
	renderer := NewTableRenderer(20, 24) // Very narrow terminal

	columns := []table.Column{
		*table.NewColumn("col1", "Column 1"),
		*table.NewColumn("col2", "Column 2"),
		*table.NewColumn("col3", "Column 3"),
		*table.NewColumn("col4", "Column 4"),
	}

	adjusted := renderer.distributeColumnWidths(columns, 20)

	// All columns should have minimum width of 5
	for i, col := range adjusted {
		if col.Width < 5 {
			t.Errorf("Column %d width should be at least 5, got %d", i, col.Width)
		}
	}
}

func TestGetMaxTableHeight(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	maxHeight := renderer.GetMaxTableHeight()

	// Should be less than terminal height (leaving room for status, etc.)
	if maxHeight >= renderer.terminalHeight {
		t.Errorf("Max table height (%d) should be less than terminal height (%d)", maxHeight, renderer.terminalHeight)
	}

	if maxHeight < 5 {
		t.Errorf("Max table height too small: %d", maxHeight)
	}
}

func TestBuildTableRow(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	columns := []table.Column{
		*table.NewColumn("col1", "Col1").WithWidth(10),
		*table.NewColumn("col2", "Col2").WithWidth(15),
	}

	result := renderer.buildTableRow(columns, func(i int, col table.Column) string {
		return strings.Repeat("X", col.Width)
	})

	expected := strings.Repeat("X", 10) + "│" + strings.Repeat("X", 15)
	if result != expected {
		t.Errorf("buildTableRow result = %q, expected %q", result, expected)
	}
}

func TestRenderTableWithCustomRenderers(t *testing.T) {
	// Create table with custom renderer
	columns := []table.Column{
		*table.NewColumn("status", "Status").WithRenderer(func(val interface{}, selected bool) string {
			if val == true {
				return "✅ Active"
			}
			return "❌ Inactive"
		}),
	}

	data := []map[string]interface{}{
		{"status": true},
		{"status": false},
	}

	tbl := table.NewWithColumns(columns)
	if err := tbl.SetData(data); err != nil {
		t.Fatalf("Failed to set data: %v", err)
	}

	renderer := NewTableRenderer(50, 20)
	result := renderer.RenderTable(tbl, 0, 0)

	if !strings.Contains(result, "✅ Active") {
		t.Error("Custom renderer should produce '✅ Active'")
	}

	if !strings.Contains(result, "❌ Inactive") {
		t.Error("Custom renderer should produce '❌ Inactive'")
	}
}
