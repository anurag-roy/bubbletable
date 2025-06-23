package table

import (
	"testing"
)

// Test data structures
type Employee struct {
	ID         int     `table:"ID,sortable,width:5"`
	Name       string  `table:"Name,sortable,width:20"`
	Department string  `table:"Department,sortable,width:15"`
	Salary     float64 `table:"Salary,sortable,width:12,format:currency"`
	StartDate  string  `table:"Start Date,sortable,width:12,format:date"`
	Active     bool    `table:"Active,sortable,width:8"`
}

type Product struct {
	SKU   string  `table:"SKU,!sortable,width:10"`
	Name  string  `table:"Product Name,sortable,width:25"`
	Price float64 `table:"Price,sortable,width:10,format:currency"`
}

func TestNew(t *testing.T) {
	table := New()
	if table == nil {
		t.Fatal("New() should return a non-nil table")
	}

	if table.SortBy != -1 {
		t.Errorf("Expected SortBy to be -1, got %d", table.SortBy)
	}

	if table.PageSize != 10 {
		t.Errorf("Expected PageSize to be 10, got %d", table.PageSize)
	}

	if len(table.Columns) != 0 {
		t.Errorf("Expected empty columns, got %d columns", len(table.Columns))
	}
}

func TestNewWithColumns(t *testing.T) {
	columns := []Column{
		*NewColumn("id", "ID").WithType(Integer).WithWidth(5),
		*NewColumn("name", "Name").WithType(String).WithWidth(20),
	}

	table := NewWithColumns(columns)
	if len(table.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(table.Columns))
	}

	if table.Columns[0].Key != "id" {
		t.Errorf("Expected first column key to be 'id', got '%s'", table.Columns[0].Key)
	}
}

func TestSetDataWithStruct(t *testing.T) {
	employees := []Employee{
		{1, "Alice Johnson", "Engineering", 75000.0, "2021-01-15", true},
		{2, "Bob Smith", "Marketing", 65000.0, "2020-03-20", true},
		{3, "Charlie Brown", "Sales", 55000.0, "2019-11-10", false},
	}

	table := New()
	err := table.SetData(employees)
	if err != nil {
		t.Fatalf("SetData failed: %v", err)
	}

	if len(table.Columns) != 6 {
		t.Errorf("Expected 6 columns, got %d", len(table.Columns))
	}

	if len(table.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(table.Rows))
	}

	if table.TotalRows != 3 {
		t.Errorf("Expected TotalRows to be 3, got %d", table.TotalRows)
	}

	// Test column inference from struct tags
	expectedColumns := map[string]struct {
		header   string
		sortable bool
	}{
		"ID":         {"ID", true},
		"Name":       {"Name", true},
		"Department": {"Department", true},
		"Salary":     {"Salary", true},
		"StartDate":  {"Start Date", true},
		"Active":     {"Active", true},
	}

	for _, col := range table.Columns {
		expected, exists := expectedColumns[col.Key]
		if !exists {
			t.Errorf("Unexpected column: %s", col.Key)
			continue
		}

		if col.Header != expected.header {
			t.Errorf("Column %s: expected header '%s', got '%s'", col.Key, expected.header, col.Header)
		}

		if col.Sortable != expected.sortable {
			t.Errorf("Column %s: expected sortable %t, got %t", col.Key, expected.sortable, col.Sortable)
		}
	}
}

func TestSetDataWithMap(t *testing.T) {
	data := []map[string]interface{}{
		{"id": 1, "name": "Alice", "age": 30},
		{"id": 2, "name": "Bob", "age": 25},
	}

	columns := []Column{
		*NewColumn("id", "ID").WithType(Integer),
		*NewColumn("name", "Name").WithType(String),
		*NewColumn("age", "Age").WithType(Integer),
	}

	table := NewWithColumns(columns)
	err := table.SetData(data)
	if err != nil {
		t.Fatalf("SetData failed: %v", err)
	}

	if len(table.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(table.Rows))
	}

	// Test first row data
	if len(table.Rows[0].Cells) != 3 {
		t.Errorf("Expected 3 cells in first row, got %d", len(table.Rows[0].Cells))
	}

	if table.Rows[0].Cells[0].Value != 1 {
		t.Errorf("Expected first cell value to be 1, got %v", table.Rows[0].Cells[0].Value)
	}
}

func TestSetDataInvalidInput(t *testing.T) {
	table := New()

	// Test non-slice input
	err := table.SetData("not a slice")
	if err == nil {
		t.Error("Expected error for non-slice input")
	}

	// Test nil input
	err = table.SetData(nil)
	if err == nil {
		t.Error("Expected error for nil input")
	}
}

func TestSortByColumn(t *testing.T) {
	employees := []Employee{
		{3, "Charlie Brown", "Sales", 55000.0, "2019-11-10", false},
		{1, "Alice Johnson", "Engineering", 75000.0, "2021-01-15", true},
		{2, "Bob Smith", "Marketing", 65000.0, "2020-03-20", true},
	}

	table := New()
	table.SetData(employees)

	// Sort by ID (ascending)
	err := table.SortByColumn(0, false)
	if err != nil {
		t.Fatalf("SortByColumn failed: %v", err)
	}

	if table.SortBy != 0 {
		t.Errorf("Expected SortBy to be 0, got %d", table.SortBy)
	}

	if table.SortDesc != false {
		t.Errorf("Expected SortDesc to be false, got %t", table.SortDesc)
	}

	// Check if rows are sorted by ID
	expectedIDs := []int{1, 2, 3}
	for i, expectedID := range expectedIDs {
		if table.Rows[i].Cells[0].Value != expectedID {
			t.Errorf("Row %d: expected ID %d, got %v", i, expectedID, table.Rows[i].Cells[0].Value)
		}
	}

	// Sort by ID (descending)
	err = table.SortByColumn(0, true)
	if err != nil {
		t.Fatalf("SortByColumn descending failed: %v", err)
	}

	if table.SortDesc != true {
		t.Errorf("Expected SortDesc to be true, got %t", table.SortDesc)
	}

	// Check if rows are sorted by ID descending
	expectedIDsDesc := []int{3, 2, 1}
	for i, expectedID := range expectedIDsDesc {
		if table.Rows[i].Cells[0].Value != expectedID {
			t.Errorf("Row %d: expected ID %d, got %v", i, expectedID, table.Rows[i].Cells[0].Value)
		}
	}
}

func TestSortByColumnInvalidIndex(t *testing.T) {
	table := New()
	table.SetData([]Employee{{1, "Alice", "Engineering", 75000.0, "2021-01-15", true}})

	err := table.SortByColumn(10, false)
	if err == nil {
		t.Error("Expected error for invalid column index")
	}

	err = table.SortByColumn(-1, false)
	if err == nil {
		t.Error("Expected error for negative column index")
	}
}

func TestClearSort(t *testing.T) {
	employees := []Employee{
		{3, "Charlie Brown", "Sales", 55000.0, "2019-11-10", false},
		{1, "Alice Johnson", "Engineering", 75000.0, "2021-01-15", true},
		{2, "Bob Smith", "Marketing", 65000.0, "2020-03-20", true},
	}

	table := New()
	table.SetData(employees)

	// Sort first
	table.SortByColumn(0, false)

	// Clear sort
	table.ClearSort()

	if table.SortBy != -1 {
		t.Errorf("Expected SortBy to be -1 after clear, got %d", table.SortBy)
	}

	// Check if original order is restored
	expectedIDs := []int{3, 1, 2}
	for i, expectedID := range expectedIDs {
		if table.Rows[i].Cells[0].Value != expectedID {
			t.Errorf("Row %d: expected ID %d, got %v", i, expectedID, table.Rows[i].Cells[0].Value)
		}
	}
}

func TestFilter(t *testing.T) {
	employees := []Employee{
		{1, "Alice Johnson", "Engineering", 75000.0, "2021-01-15", true},
		{2, "Bob Smith", "Marketing", 65000.0, "2020-03-20", true},
		{3, "Charlie Brown", "Sales", 55000.0, "2019-11-10", false},
	}

	table := New()
	table.SetData(employees)

	// Filter by name
	filtered := table.Filter("Alice")
	if filtered == nil {
		t.Fatal("Filter returned nil")
	}

	if len(filtered.Rows) != 1 {
		t.Errorf("Expected 1 filtered row, got %d", len(filtered.Rows))
	}

	if filtered.Rows[0].Cells[1].Value != "Alice Johnson" {
		t.Errorf("Expected filtered row to contain Alice Johnson, got %v", filtered.Rows[0].Cells[1].Value)
	}

	// Filter by department
	filtered = table.Filter("Engineering")
	if len(filtered.Rows) != 1 {
		t.Errorf("Expected 1 filtered row for Engineering, got %d", len(filtered.Rows))
	}

	// Filter with no matches
	filtered = table.Filter("NonExistent")
	if len(filtered.Rows) != 0 {
		t.Errorf("Expected 0 filtered rows for non-existent term, got %d", len(filtered.Rows))
	}

	// Case insensitive search
	filtered = table.Filter("alice")
	if len(filtered.Rows) != 1 {
		t.Errorf("Expected 1 filtered row for case insensitive search, got %d", len(filtered.Rows))
	}
}

func TestGetPage(t *testing.T) {
	employees := make([]Employee, 25)
	for i := 0; i < 25; i++ {
		employees[i] = Employee{
			ID:   i + 1,
			Name: "Employee " + string(rune('A'+i%26)),
		}
	}

	table := New().WithPageSize(10)
	table.SetData(employees)

	// Test first page
	page := table.GetPage(0)
	if len(page) != 10 {
		t.Errorf("Expected 10 rows in first page, got %d", len(page))
	}

	if page[0].Cells[0].Value != 1 {
		t.Errorf("Expected first row ID to be 1, got %v", page[0].Cells[0].Value)
	}

	// Test second page
	page = table.GetPage(1)
	if len(page) != 10 {
		t.Errorf("Expected 10 rows in second page, got %d", len(page))
	}

	if page[0].Cells[0].Value != 11 {
		t.Errorf("Expected first row of second page ID to be 11, got %v", page[0].Cells[0].Value)
	}

	// Test last page (partial)
	page = table.GetPage(2)
	if len(page) != 5 {
		t.Errorf("Expected 5 rows in last page, got %d", len(page))
	}

	// Test invalid page
	page = table.GetPage(10)
	if len(page) != 0 {
		t.Errorf("Expected 0 rows for invalid page, got %d", len(page))
	}
}

func TestGetTotalPages(t *testing.T) {
	table := New().WithPageSize(10)

	// Test with exact multiple
	employees := make([]Employee, 20)
	table.SetData(employees)
	if table.GetTotalPages() != 2 {
		t.Errorf("Expected 2 total pages for 20 items, got %d", table.GetTotalPages())
	}

	// Test with remainder
	employees = make([]Employee, 25)
	table.SetData(employees)
	if table.GetTotalPages() != 3 {
		t.Errorf("Expected 3 total pages for 25 items, got %d", table.GetTotalPages())
	}

	// Test with no data
	table.SetData([]Employee{})
	if table.GetTotalPages() != 1 {
		t.Errorf("Expected 1 total page for no data (empty page), got %d", table.GetTotalPages())
	}
}

func TestNewColumn(t *testing.T) {
	col := NewColumn("test", "Test Header")

	if col.Key != "test" {
		t.Errorf("Expected key 'test', got '%s'", col.Key)
	}

	if col.Header != "Test Header" {
		t.Errorf("Expected header 'Test Header', got '%s'", col.Header)
	}

	if col.Type != String {
		t.Errorf("Expected default type String, got %v", col.Type)
	}

	if col.Width != 15 {
		t.Errorf("Expected default width 15, got %d", col.Width)
	}

	if !col.Sortable {
		t.Error("Expected default sortable to be true")
	}

	if !col.Searchable {
		t.Error("Expected default searchable to be true")
	}
}

func TestColumnBuilderPattern(t *testing.T) {
	col := NewColumn("price", "Price").
		WithType(Float).
		WithWidth(12).
		WithSortable(true).
		WithSearchable(false).
		WithFormatter(CurrencyFormatter)

	if col.Type != Float {
		t.Errorf("Expected type Float, got %v", col.Type)
	}

	if col.Width != 12 {
		t.Errorf("Expected width 12, got %d", col.Width)
	}

	if !col.Sortable {
		t.Error("Expected sortable to be true")
	}

	if col.Searchable {
		t.Error("Expected searchable to be false")
	}

	// Test formatter
	formatted := col.Formatter(99.99)
	if formatted != "$99.99" {
		t.Errorf("Expected formatted value '$99.99', got '%s'", formatted)
	}
}

func TestStructTagParsing(t *testing.T) {
	products := []Product{
		{"SKU-001", "Test Product", 99.99},
	}

	table := New()
	table.SetData(products)

	// Find SKU column - should not be sortable due to !sortable tag
	var skuCol *Column
	for i, col := range table.Columns {
		if col.Key == "SKU" {
			skuCol = &table.Columns[i]
			break
		}
	}

	if skuCol == nil {
		t.Fatal("SKU column not found")
	}

	if skuCol.Sortable {
		t.Error("SKU column should not be sortable due to !sortable tag")
	}

	// Width should be set based on struct tag parsing
	if skuCol.Width < 5 {
		t.Errorf("Expected SKU column width to be at least 5, got %d", skuCol.Width)
	}
}

func TestCompareCells(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Cell
		expected int
	}{
		{
			name:     "Integer ascending",
			a:        Cell{Value: 1, Type: Integer},
			b:        Cell{Value: 2, Type: Integer},
			expected: -1,
		},
		{
			name:     "Integer descending",
			a:        Cell{Value: 2, Type: Integer},
			b:        Cell{Value: 1, Type: Integer},
			expected: 1,
		},
		{
			name:     "String ascending",
			a:        Cell{Value: "apple", Type: String},
			b:        Cell{Value: "banana", Type: String},
			expected: -1,
		},
		{
			name:     "Float ascending",
			a:        Cell{Value: 1.5, Type: Float},
			b:        Cell{Value: 2.5, Type: Float},
			expected: -1,
		},
		{
			name:     "Boolean true > false",
			a:        Cell{Value: true, Type: Boolean},
			b:        Cell{Value: false, Type: Boolean},
			expected: 1,
		},
		{
			name:     "Equal values",
			a:        Cell{Value: "same", Type: String},
			b:        Cell{Value: "same", Type: String},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareCells(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareCells(%v, %v) = %d, expected %d", tt.a.Value, tt.b.Value, result, tt.expected)
			}
		})
	}
}

func TestBuilderPatternChaining(t *testing.T) {
	table := New().
		WithPageSize(25).
		WithData([]Employee{
			{1, "Alice", "Engineering", 75000.0, "2021-01-15", true},
		})

	if table.PageSize != 25 {
		t.Errorf("Expected PageSize 25, got %d", table.PageSize)
	}

	if len(table.Rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(table.Rows))
	}
}
