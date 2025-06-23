package renderer

import (
	"strings"
	"testing"

	"github.com/anurag-roy/bubbletable/table"
)

func TestNewTableRenderer(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	if renderer.terminalWidth != 80 {
		t.Errorf("Expected terminalWidth 80, got %d", renderer.terminalWidth)
	}

	if renderer.terminalHeight != 24 {
		t.Errorf("Expected terminalHeight 24, got %d", renderer.terminalHeight)
	}

	if renderer.theme.Name != DefaultTheme.Name {
		t.Errorf("Expected default theme, got %s", renderer.theme.Name)
	}
}

func TestNewTableRendererWithTheme(t *testing.T) {
	renderer := NewTableRendererWithTheme(100, 30, DraculaTheme)

	if renderer.theme.Name != "Dracula" {
		t.Errorf("Expected Dracula theme, got %s", renderer.theme.Name)
	}
}

func TestUpdateSize(t *testing.T) {
	renderer := NewTableRenderer(80, 24)
	renderer.UpdateSize(120, 40)

	if renderer.terminalWidth != 120 {
		t.Errorf("Expected updated width 120, got %d", renderer.terminalWidth)
	}

	if renderer.terminalHeight != 40 {
		t.Errorf("Expected updated height 40, got %d", renderer.terminalHeight)
	}
}

func TestSetTheme(t *testing.T) {
	renderer := NewTableRenderer(80, 24)
	renderer.SetTheme(MonokaiTheme)

	if renderer.theme.Name != "Monokai" {
		t.Errorf("Expected Monokai theme, got %s", renderer.theme.Name)
	}
}

func TestDistributeColumnWidths(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	columns := []table.Column{
		*table.NewColumn("col1", "Column 1"),
		*table.NewColumn("col2", "Column 2"),
		*table.NewColumn("col3", "Column 3"),
	}

	// Test with 60 available width (3 columns, so 2 separators = 58 content width)
	adjusted := renderer.distributeColumnWidths(columns, 60)

	if len(adjusted) != 3 {
		t.Errorf("Expected 3 adjusted columns, got %d", len(adjusted))
	}

	totalWidth := 0
	for _, col := range adjusted {
		totalWidth += col.Width
		if col.Width < 5 {
			t.Errorf("Column width should be at least 5, got %d", col.Width)
		}
	}

	// Total content width should be close to available width minus separators
	expectedContentWidth := 58 // 60 - 2 separators
	if totalWidth != expectedContentWidth {
		t.Errorf("Expected total content width %d, got %d", expectedContentWidth, totalWidth)
	}
}

func TestRenderTable(t *testing.T) {
	// Create test data
	type TestData struct {
		ID   int    `table:"ID,sortable,width:5"`
		Name string `table:"Name,sortable,width:15"`
	}

	data := []TestData{
		{1, "Alice"},
		{2, "Bob"},
	}

	// Create table
	tbl := table.New()
	err := tbl.SetData(data)
	if err != nil {
		t.Fatalf("Failed to set table data: %v", err)
	}

	// Create renderer
	renderer := NewTableRenderer(50, 20)

	// Render table
	result := renderer.RenderTable(tbl, 0, 0)

	// Basic checks
	if result == "" {
		t.Error("Rendered table should not be empty")
	}

	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Errorf("Expected at least 3 lines (header, separator, data), got %d", len(lines))
	}

	// Check if header contains column names
	headerLine := lines[0]
	if !strings.Contains(headerLine, "ID") {
		t.Error("Header should contain 'ID'")
	}
	if !strings.Contains(headerLine, "Name") {
		t.Error("Header should contain 'Name'")
	}

	// Check if data rows contain actual data
	dataLines := lines[2:] // Skip header and separator
	foundAlice := false
	foundBob := false

	for _, line := range dataLines {
		if strings.Contains(line, "Alice") {
			foundAlice = true
		}
		if strings.Contains(line, "Bob") {
			foundBob = true
		}
	}

	if !foundAlice {
		t.Error("Rendered table should contain 'Alice'")
	}
	if !foundBob {
		t.Error("Rendered table should contain 'Bob'")
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
	tbl.SetData(data)
	tbl.SortByColumn(0, false) // Sort by ID ascending

	renderer := NewTableRenderer(50, 20)
	result := renderer.RenderTable(tbl, 0, 0)

	lines := strings.Split(result, "\n")
	headerLine := lines[0]

	// Check for sort indicator
	if !strings.Contains(headerLine, "↑") {
		t.Error("Header should contain ascending sort indicator '↑'")
	}

	// Sort descending
	tbl.SortByColumn(0, true)
	result = renderer.RenderTable(tbl, 0, 0)
	lines = strings.Split(result, "\n")
	headerLine = lines[0]

	if !strings.Contains(headerLine, "↓") {
		t.Error("Header should contain descending sort indicator '↓'")
	}
}

func TestRenderTableNilInput(t *testing.T) {
	renderer := NewTableRenderer(80, 24)
	result := renderer.RenderTable(nil, 0, 0)

	expected := "No table data available"
	if result != expected {
		t.Errorf("Expected %q for nil table, got %q", expected, result)
	}
}

func TestTruncateText(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	tests := []struct {
		name     string
		text     string
		width    int
		expected string
	}{
		{"short text", "hello", 10, "hello"},
		{"exact width", "12345", 5, "12345"},
		{"long text", "this is a very long text", 10, "this is..."},
		{"empty text", "", 5, ""},
		{"zero width", "hello", 0, ""},
		{"width 1", "hello", 1, "h"},
		{"width 2", "hello", 2, "he"},
		{"width 3", "hello", 3, "hel"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderer.truncateText(tt.text, tt.width)
			if result != tt.expected {
				t.Errorf("truncateText(%q, %d) = %q, expected %q", tt.text, tt.width, result, tt.expected)
			}
		})
	}
}

func TestGetOptimalPageSize(t *testing.T) {
	renderer := NewTableRenderer(80, 24)

	pageSize := renderer.GetOptimalPageSize()

	// Should be reasonable for terminal height
	if pageSize < 5 {
		t.Errorf("Optimal page size too small: %d", pageSize)
	}

	if pageSize > renderer.terminalHeight {
		t.Errorf("Optimal page size (%d) should not exceed terminal height (%d)", pageSize, renderer.terminalHeight)
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
	tbl.SetData(data)

	renderer := NewTableRenderer(50, 20)
	result := renderer.RenderTable(tbl, 0, 0)

	if !strings.Contains(result, "✅ Active") {
		t.Error("Custom renderer should produce '✅ Active'")
	}

	if !strings.Contains(result, "❌ Inactive") {
		t.Error("Custom renderer should produce '❌ Inactive'")
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
