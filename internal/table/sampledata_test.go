package table

import (
	"testing"
)

func TestNewSampleDataGenerator(t *testing.T) {
	generator := NewSampleDataGenerator()
	if generator == nil {
		t.Error("Expected non-nil generator")
	}
	if generator.rand == nil {
		t.Error("Expected non-nil random generator")
	}
}

func TestGenerateEmployeeTable(t *testing.T) {
	generator := NewSampleDataGenerator()
	table := generator.GenerateEmployeeTable()

	if table == nil {
		t.Fatal("Expected non-nil table")
	}

	// Check columns
	expectedColumns := 6
	if len(table.Columns) != expectedColumns {
		t.Errorf("Expected %d columns, got %d", expectedColumns, len(table.Columns))
	}

	// Check column names
	expectedNames := []string{"ID", "Name", "Department", "Salary", "Start Date", "Active"}
	for i, expected := range expectedNames {
		if i >= len(table.Columns) {
			t.Errorf("Missing column at index %d", i)
			continue
		}
		if table.Columns[i].Name != expected {
			t.Errorf("Expected column name '%s', got '%s'", expected, table.Columns[i].Name)
		}
	}

	// Check that data was generated
	if len(table.Rows) == 0 {
		t.Error("Expected some rows to be generated")
	}

	// Check first row has correct number of cells
	if len(table.Rows) > 0 && len(table.Rows[0].Cells) != expectedColumns {
		t.Errorf("Expected %d cells in first row, got %d", expectedColumns, len(table.Rows[0].Cells))
	}

	// Check data types
	if len(table.Rows) > 0 {
		row := table.Rows[0]
		expectedTypes := []DataType{Integer, String, String, Float, Date, Boolean}
		for i, expectedType := range expectedTypes {
			if i >= len(row.Cells) {
				continue
			}
			if row.Cells[i].Type != expectedType {
				t.Errorf("Expected type %d for column %d, got %d", expectedType, i, row.Cells[i].Type)
			}
		}
	}
}

func TestGenerateProductTable(t *testing.T) {
	generator := NewSampleDataGenerator()
	table := generator.GenerateProductTable()

	if table == nil {
		t.Fatal("Expected non-nil table")
	}

	// Check columns
	expectedColumns := 6
	if len(table.Columns) != expectedColumns {
		t.Errorf("Expected %d columns, got %d", expectedColumns, len(table.Columns))
	}

	// Check that data was generated
	if len(table.Rows) == 0 {
		t.Error("Expected some rows to be generated")
	}

	// Check page size is set correctly
	if table.PageSize != 15 {
		t.Errorf("Expected page size 15, got %d", table.PageSize)
	}
}

func TestGenerateSampleTable(t *testing.T) {
	generator := NewSampleDataGenerator()

	// Test employee table
	empTable := generator.GenerateSampleTable("employees")
	if empTable == nil {
		t.Error("Expected non-nil employee table")
	}

	// Test product table
	prodTable := generator.GenerateSampleTable("products")
	if prodTable == nil {
		t.Error("Expected non-nil product table")
	}

	// Test default fallback (unknown table name)
	defaultTable := generator.GenerateSampleTable("unknown")
	if defaultTable == nil {
		t.Error("Expected non-nil default table")
	}
}

func TestGetSampleTableNames(t *testing.T) {
	names := GetSampleTableNames()
	if len(names) == 0 {
		t.Error("Expected some sample table names")
	}

	// Check for expected names
	expectedNames := map[string]bool{
		"employees": false,
		"products":  false,
	}

	for _, name := range names {
		if _, exists := expectedNames[name]; exists {
			expectedNames[name] = true
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected to find sample table name '%s'", name)
		}
	}
}
