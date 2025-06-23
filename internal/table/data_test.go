package table

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTable(t *testing.T) {
	columns := []Column{
		{Name: "ID", Type: Integer, Width: 5, Sortable: true},
		{Name: "Name", Type: String, Width: 20, Sortable: true},
	}

	table := NewTable(columns)

	if len(table.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(table.Columns))
	}

	if table.PageSize != 10 {
		t.Errorf("Expected default page size 10, got %d", table.PageSize)
	}

	if table.SortBy != -1 {
		t.Errorf("Expected SortBy to be -1, got %d", table.SortBy)
	}

	if table.TotalRows != 0 {
		t.Errorf("Expected TotalRows to be 0, got %d", table.TotalRows)
	}
}

func TestAddRow(t *testing.T) {
	columns := []Column{
		{Name: "ID", Type: Integer, Width: 5, Sortable: true},
		{Name: "Name", Type: String, Width: 20, Sortable: true},
		{Name: "Active", Type: Boolean, Width: 8, Sortable: true},
	}

	table := NewTable(columns)

	// Test successful row addition
	err := table.AddRow(1, "John Doe", true)
	if err != nil {
		t.Errorf("Unexpected error adding row: %v", err)
	}

	if table.TotalRows != 1 {
		t.Errorf("Expected TotalRows to be 1, got %d", table.TotalRows)
	}

	if len(table.Rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(table.Rows))
	}

	// Test error with wrong number of values
	err = table.AddRow(2, "Jane Doe") // Missing boolean value
	if err == nil {
		t.Error("Expected error when adding row with wrong number of values")
	}
}

func TestGetPage(t *testing.T) {
	columns := []Column{
		{Name: "ID", Type: Integer, Width: 5, Sortable: true},
		{Name: "Name", Type: String, Width: 20, Sortable: true},
	}

	table := NewTable(columns)
	table.PageSize = 3

	// Add test data
	for i := 1; i <= 10; i++ {
		table.AddRow(i, fmt.Sprintf("Person %d", i))
	}

	// Test first page
	page0 := table.GetPage(0)
	if len(page0) != 3 {
		t.Errorf("Expected 3 rows on first page, got %d", len(page0))
	}

	// Test middle page
	page1 := table.GetPage(1)
	if len(page1) != 3 {
		t.Errorf("Expected 3 rows on second page, got %d", len(page1))
	}

	// Test last page (partial)
	page3 := table.GetPage(3)
	if len(page3) != 1 {
		t.Errorf("Expected 1 row on last page, got %d", len(page3))
	}

	// Test page beyond data
	page10 := table.GetPage(10)
	if len(page10) != 0 {
		t.Errorf("Expected 0 rows for page beyond data, got %d", len(page10))
	}
}

func TestGetTotalPages(t *testing.T) {
	columns := []Column{
		{Name: "ID", Type: Integer, Width: 5, Sortable: true},
	}

	table := NewTable(columns)
	table.PageSize = 3

	// Empty table
	if table.GetTotalPages() != 1 {
		t.Errorf("Expected 1 page for empty table, got %d", table.GetTotalPages())
	}

	// Add 7 rows (should be 3 pages with page size 3)
	for i := 1; i <= 7; i++ {
		table.AddRow(i)
	}

	if table.GetTotalPages() != 3 {
		t.Errorf("Expected 3 pages for 7 rows with page size 3, got %d", table.GetTotalPages())
	}

	// Add 2 more rows (should still be 3 pages)
	table.AddRow(8)
	table.AddRow(9)

	if table.GetTotalPages() != 3 {
		t.Errorf("Expected 3 pages for 9 rows with page size 3, got %d", table.GetTotalPages())
	}

	// Add 1 more row (should be 4 pages)
	table.AddRow(10)

	if table.GetTotalPages() != 4 {
		t.Errorf("Expected 4 pages for 10 rows with page size 3, got %d", table.GetTotalPages())
	}
}

func TestSortByColumn(t *testing.T) {
	columns := []Column{
		{Name: "ID", Type: Integer, Width: 5, Sortable: true},
		{Name: "Name", Type: String, Width: 20, Sortable: true},
		{Name: "Score", Type: Float, Width: 8, Sortable: false},
	}

	table := NewTable(columns)

	// Add test data
	table.AddRow(3, "Charlie", 85.5)
	table.AddRow(1, "Alice", 92.0)
	table.AddRow(2, "Bob", 78.5)

	// Test sorting by ID (ascending)
	err := table.SortByColumn(0, false)
	if err != nil {
		t.Errorf("Unexpected error sorting by ID: %v", err)
	}

	if table.Rows[0].Cells[0].Value != 1 {
		t.Errorf("Expected first row ID to be 1, got %v", table.Rows[0].Cells[0].Value)
	}

	// Test sorting by Name (descending)
	err = table.SortByColumn(1, true)
	if err != nil {
		t.Errorf("Unexpected error sorting by Name: %v", err)
	}

	if table.Rows[0].Cells[1].Value != "Charlie" {
		t.Errorf("Expected first row Name to be Charlie, got %v", table.Rows[0].Cells[1].Value)
	}

	// Test error for non-sortable column
	err = table.SortByColumn(2, false)
	if err == nil {
		t.Error("Expected error when sorting non-sortable column")
	}

	// Test error for invalid column index
	err = table.SortByColumn(10, false)
	if err == nil {
		t.Error("Expected error for invalid column index")
	}
}

func TestFilter(t *testing.T) {
	columns := []Column{
		{Name: "ID", Type: Integer, Width: 5, Sortable: true},
		{Name: "Name", Type: String, Width: 20, Sortable: true},
		{Name: "Email", Type: String, Width: 25, Sortable: true},
	}

	table := NewTable(columns)

	// Add test data
	table.AddRow(1, "Alice Johnson", "alice@example.com")
	table.AddRow(2, "Bob Smith", "bob@test.com")
	table.AddRow(3, "Charlie Brown", "charlie@example.com")
	table.AddRow(4, "Diana Prince", "diana@company.com")

	// Test filtering by name
	filtered := table.Filter("alice")
	if len(filtered.Rows) != 1 {
		t.Errorf("Expected 1 row when filtering by 'alice', got %d", len(filtered.Rows))
	}

	// Test filtering by email domain
	filtered = table.Filter("example.com")
	if len(filtered.Rows) != 2 {
		t.Errorf("Expected 2 rows when filtering by 'example.com', got %d", len(filtered.Rows))
	}

	// Test case insensitive filtering
	filtered = table.Filter("ALICE")
	if len(filtered.Rows) != 1 {
		t.Errorf("Expected 1 row when filtering by 'ALICE' (case insensitive), got %d", len(filtered.Rows))
	}

	// Test no matches
	filtered = table.Filter("xyz")
	if len(filtered.Rows) != 0 {
		t.Errorf("Expected 0 rows when filtering by 'xyz', got %d", len(filtered.Rows))
	}

	// Test empty search term returns original table
	filtered = table.Filter("")
	if filtered != table {
		t.Error("Expected original table when search term is empty")
	}
}

func TestCompareCells(t *testing.T) {
	// Test string comparison
	cellA := Cell{Value: "apple", Type: String}
	cellB := Cell{Value: "banana", Type: String}
	result := compareCells(cellA, cellB)
	if result >= 0 {
		t.Error("Expected 'apple' < 'banana'")
	}

	// Test integer comparison
	cellA = Cell{Value: 5, Type: Integer}
	cellB = Cell{Value: 10, Type: Integer}
	result = compareCells(cellA, cellB)
	if result >= 0 {
		t.Error("Expected 5 < 10")
	}

	// Test float comparison
	cellA = Cell{Value: 3.14, Type: Float}
	cellB = Cell{Value: 2.71, Type: Float}
	result = compareCells(cellA, cellB)
	if result <= 0 {
		t.Error("Expected 3.14 > 2.71")
	}

	// Test boolean comparison
	cellA = Cell{Value: true, Type: Boolean}
	cellB = Cell{Value: false, Type: Boolean}
	result = compareCells(cellA, cellB)
	if result <= 0 {
		t.Error("Expected true > false")
	}
}

func TestFormatCellValue(t *testing.T) {
	// Test string formatting
	cell := Cell{Value: "hello", Type: String}
	formatted := formatCellValue(cell)
	if formatted != "hello" {
		t.Errorf("Expected 'hello', got '%s'", formatted)
	}

	// Test float formatting
	cell = Cell{Value: 3.14159, Type: Float}
	formatted = formatCellValue(cell)
	if formatted != "3.14" {
		t.Errorf("Expected '3.14', got '%s'", formatted)
	}

	// Test date formatting
	cell = Cell{Value: "2023-12-25", Type: Date}
	formatted = formatCellValue(cell)
	if formatted != "2023-12-25" {
		t.Errorf("Expected '2023-12-25', got '%s'", formatted)
	}
}

func TestParseDate(t *testing.T) {
	// Test ISO format
	date, err := parseDate("2023-12-25")
	if err != nil {
		t.Errorf("Unexpected error parsing ISO date: %v", err)
	}
	expected := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	if date.Format("2006-01-02") != expected.Format("2006-01-02") {
		t.Errorf("Expected %s, got %s", expected.Format("2006-01-02"), date.Format("2006-01-02"))
	}

	// Test US format
	date, err = parseDate("12/25/2023")
	if err != nil {
		t.Errorf("Unexpected error parsing US date: %v", err)
	}

	// Test invalid date
	_, err = parseDate("invalid-date")
	if err == nil {
		t.Error("Expected error parsing invalid date")
	}
}
