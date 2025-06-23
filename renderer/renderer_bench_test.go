package renderer

import (
	"testing"

	"github.com/anurag-roy/bubbletable/table"
)

// Benchmark data structure
type BenchData struct {
	ID   int    `table:"ID,sortable,width:5"`
	Name string `table:"Name,sortable,width:20"`
	Desc string `table:"Description,sortable,width:30"`
}

func generateRendererBenchData(n int) []BenchData {
	data := make([]BenchData, n)
	for i := 0; i < n; i++ {
		data[i] = BenchData{
			ID:   i + 1,
			Name: "Item " + string(rune('A'+(i%26))),
			Desc: "This is a longer description for item number " + string(rune('A'+(i%26))),
		}
	}
	return data
}

func BenchmarkNewTableRenderer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewTableRenderer(80, 24)
	}
}

func BenchmarkRenderSmallTable(b *testing.B) {
	data := generateRendererBenchData(10)
	tbl := table.New()
	tbl.SetData(data)
	renderer := NewTableRenderer(80, 24)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		renderer.RenderTable(tbl, 0, 0)
	}
}

func BenchmarkRenderMediumTable(b *testing.B) {
	data := generateRendererBenchData(100)
	tbl := table.New()
	tbl.SetData(data)
	renderer := NewTableRenderer(80, 24)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		renderer.RenderTable(tbl, 0, 0)
	}
}

func BenchmarkRenderLargeTable(b *testing.B) {
	data := generateRendererBenchData(1000)
	tbl := table.New()
	tbl.SetData(data)
	renderer := NewTableRenderer(80, 24)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		renderer.RenderTable(tbl, 0, 0)
	}
}

func BenchmarkDistributeColumnWidths(b *testing.B) {
	columns := []table.Column{
		*table.NewColumn("col1", "Column 1"),
		*table.NewColumn("col2", "Column 2"),
		*table.NewColumn("col3", "Column 3"),
		*table.NewColumn("col4", "Column 4"),
		*table.NewColumn("col5", "Column 5"),
	}
	renderer := NewTableRenderer(80, 24)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		renderer.distributeColumnWidths(columns, 70)
	}
}

func BenchmarkTruncateText(b *testing.B) {
	text := "This is a very long text that needs to be truncated for display"
	renderer := NewTableRenderer(80, 24)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		renderer.truncateText(text, 20)
	}
}

func BenchmarkSetTheme(b *testing.B) {
	renderer := NewTableRenderer(80, 24)
	themes := []Theme{DefaultTheme, DraculaTheme, MonokaiTheme}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		renderer.SetTheme(themes[i%len(themes)])
	}
}

func BenchmarkUpdateSize(b *testing.B) {
	renderer := NewTableRenderer(80, 24)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		renderer.UpdateSize(100+i%50, 30+i%20)
	}
}

func BenchmarkRenderWithSorting(b *testing.B) {
	data := generateRendererBenchData(100)
	tbl := table.New()
	tbl.SetData(data)
	tbl.SortByColumn(0, false) // Sort by ID
	renderer := NewTableRenderer(80, 24)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		renderer.RenderTable(tbl, 0, 0)
	}
}

func BenchmarkRenderWideTable(b *testing.B) {
	// Create table with many columns
	columns := make([]table.Column, 10)
	for i := 0; i < 10; i++ {
		columns[i] = *table.NewColumn("col"+string(rune('A'+i)), "Column "+string(rune('A'+i)))
	}

	data := make([]map[string]interface{}, 50)
	for i := 0; i < 50; i++ {
		row := make(map[string]interface{})
		for j := 0; j < 10; j++ {
			row["col"+string(rune('A'+j))] = "Value " + string(rune('A'+j)) + string(rune('0'+i%10))
		}
		data[i] = row
	}

	tbl := table.NewWithColumns(columns)
	tbl.SetData(data)
	renderer := NewTableRenderer(200, 50) // Wide terminal
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		renderer.RenderTable(tbl, 0, 0)
	}
}

func BenchmarkRenderNarrowTerminal(b *testing.B) {
	data := generateRendererBenchData(50)
	tbl := table.New()
	tbl.SetData(data)
	renderer := NewTableRenderer(40, 20) // Very narrow
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		renderer.RenderTable(tbl, 0, 0)
	}
}
