package table

import (
	"fmt"
	"testing"
)

// Benchmark data structure
type BenchEmployee struct {
	ID         int     `table:"ID,sortable,width:5"`
	Name       string  `table:"Name,sortable,width:20"`
	Department string  `table:"Department,sortable,width:15"`
	Salary     float64 `table:"Salary,sortable,width:12,format:currency"`
	Active     bool    `table:"Active,sortable,width:8"`
}

func generateBenchData(n int) []BenchEmployee {
	employees := make([]BenchEmployee, n)
	departments := []string{"Engineering", "Marketing", "Sales", "HR", "Finance"}

	for i := 0; i < n; i++ {
		employees[i] = BenchEmployee{
			ID:         i + 1,
			Name:       fmt.Sprintf("Employee_%d", i+1),
			Department: departments[i%len(departments)],
			Salary:     float64(50000 + (i*1000)%50000),
			Active:     i%3 == 0,
		}
	}
	return employees
}

func BenchmarkSetData100(b *testing.B) {
	data := generateBenchData(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table := New()
		table.SetData(data)
	}
}

func BenchmarkSetData1000(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table := New()
		table.SetData(data)
	}
}

func BenchmarkSetData10000(b *testing.B) {
	data := generateBenchData(10000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table := New()
		table.SetData(data)
	}
}

func BenchmarkSortByColumn(b *testing.B) {
	data := generateBenchData(1000)
	table := New()
	table.SetData(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table.SortByColumn(0, false)
	}
}

func BenchmarkSortByColumnLarge(b *testing.B) {
	data := generateBenchData(10000)
	table := New()
	table.SetData(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table.SortByColumn(1, false) // Sort by name
	}
}

func BenchmarkFilter(b *testing.B) {
	data := generateBenchData(1000)
	table := New()
	table.SetData(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table.Filter("Engineering")
	}
}

func BenchmarkFilterLarge(b *testing.B) {
	data := generateBenchData(10000)
	table := New()
	table.SetData(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table.Filter("Employee_5")
	}
}

func BenchmarkGetPage(b *testing.B) {
	data := generateBenchData(1000)
	table := New().WithPageSize(50)
	table.SetData(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table.GetPage(i % 20) // Cycle through pages
	}
}

func BenchmarkGetPageLarge(b *testing.B) {
	data := generateBenchData(10000)
	table := New().WithPageSize(100)
	table.SetData(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table.GetPage(i % 100) // Cycle through pages
	}
}

func BenchmarkComplexWorkflow(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table := New().WithPageSize(50)
		table.SetData(data)
		table.SortByColumn(0, false)
		filtered := table.Filter("Engineering")
		filtered.GetPage(0)
	}
}

func BenchmarkMapData(b *testing.B) {
	data := make([]map[string]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = map[string]interface{}{
			"id":   i + 1,
			"name": fmt.Sprintf("Employee_%d", i+1),
			"dept": "Engineering",
		}
	}

	columns := []Column{
		*NewColumn("id", "ID").WithType(Integer),
		*NewColumn("name", "Name").WithType(String),
		*NewColumn("dept", "Department").WithType(String),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		table := NewWithColumns(columns)
		table.SetData(data)
	}
}

func BenchmarkColumnCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		col := NewColumn("test", "Test").
			WithType(String).
			WithWidth(20).
			WithSortable(true).
			WithSearchable(true).
			WithFormatter(DefaultFormatter)
		_ = col
	}
}

func BenchmarkCellComparison(b *testing.B) {
	cell1 := Cell{Value: "Alice", Type: String}
	cell2 := Cell{Value: "Bob", Type: String}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		compareCells(cell1, cell2)
	}
}

func BenchmarkIntegerComparison(b *testing.B) {
	cell1 := Cell{Value: 100, Type: Integer}
	cell2 := Cell{Value: 200, Type: Integer}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		compareCells(cell1, cell2)
	}
}

func BenchmarkFloatComparison(b *testing.B) {
	cell1 := Cell{Value: 100.5, Type: Float}
	cell2 := Cell{Value: 200.7, Type: Float}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		compareCells(cell1, cell2)
	}
}
