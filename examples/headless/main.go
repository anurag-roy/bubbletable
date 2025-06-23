package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/anurag-roy/bubbletable/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CustomModel demonstrates headless usage with full control over rendering
type CustomModel struct {
	table       *table.Table
	currentPage int
	selectedRow int
}

// Sample data structure
type Task struct {
	ID        int    `table:"ID,sortable,width:5"`
	Title     string `table:"Title,sortable,width:30"`
	Status    string `table:"Status,sortable,width:12"`
	Priority  string `table:"Priority,sortable,width:10"`
	Assignee  string `table:"Assignee,sortable,width:15"`
	Completed bool   `table:"Done,sortable,width:6"`
}

func NewCustomModel() *CustomModel {
	// Create sample tasks
	tasks := []Task{
		{1, "Implement user authentication", "In Progress", "High", "Alice", false},
		{2, "Design database schema", "Completed", "Medium", "Bob", true},
		{3, "Write API documentation", "Todo", "Low", "Charlie", false},
		{4, "Set up CI/CD pipeline", "In Progress", "High", "Diana", false},
		{5, "Create unit tests", "Todo", "Medium", "Edward", false},
		{6, "Optimize database queries", "Completed", "High", "Fiona", true},
		{7, "Implement caching layer", "In Progress", "Medium", "George", false},
		{8, "Fix security vulnerabilities", "Todo", "Critical", "Helen", false},
		{9, "Update user interface", "In Progress", "Low", "Ivan", false},
		{10, "Deploy to production", "Todo", "Critical", "Julia", false},
	}

	// Create headless table with custom configuration
	tbl := table.New().
		WithData(tasks).
		WithPageSize(5)

	// Customize columns with custom formatters
	for i, col := range tbl.Columns {
		switch col.Key {
		case "Status":
			tbl.Columns[i].Formatter = func(value interface{}) string {
				status := value.(string)
				switch status {
				case "Completed":
					return "✅ " + status
				case "In Progress":
					return "🔄 " + status
				case "Todo":
					return "📋 " + status
				default:
					return status
				}
			}
		case "Priority":
			tbl.Columns[i].Formatter = func(value interface{}) string {
				priority := value.(string)
				switch priority {
				case "Critical":
					return "🔴 " + priority
				case "High":
					return "🟠 " + priority
				case "Medium":
					return "🟡 " + priority
				case "Low":
					return "🟢 " + priority
				default:
					return priority
				}
			}
		case "Completed":
			tbl.Columns[i].Formatter = table.BooleanFormatter("✅", "❌")
		}
	}

	return &CustomModel{
		table:       tbl,
		currentPage: 0,
		selectedRow: 0,
	}
}

func (m *CustomModel) Init() tea.Cmd {
	return nil
}

func (m *CustomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.selectedRow > 0 {
				m.selectedRow--
			}
		case "down", "j":
			pageData := m.table.GetPage(m.currentPage)
			if m.selectedRow < len(pageData)-1 {
				m.selectedRow++
			}
		case "left", "h":
			if m.currentPage > 0 {
				m.currentPage--
				m.selectedRow = 0
			}
		case "right", "l":
			if m.currentPage < m.table.GetTotalPages()-1 {
				m.currentPage++
				m.selectedRow = 0
			}
		case "1", "2", "3", "4", "5", "6":
			// Sort by column
			colIndex := int(msg.String()[0] - '1')
			if colIndex < len(m.table.Columns) {
				if m.table.SortBy == colIndex {
					if !m.table.SortDesc {
						m.table.SortByColumn(colIndex, true)
					} else {
						m.table.ClearSort()
					}
				} else {
					m.table.SortByColumn(colIndex, false)
				}
				m.currentPage = 0
				m.selectedRow = 0
			}
		}
	}
	return m, nil
}

func (m *CustomModel) View() string {
	// Custom rendering using the headless table data
	var b strings.Builder

	// Custom header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D4CDB")).
		Padding(0, 2)

	b.WriteString(headerStyle.Render("🚀 Task Management Dashboard"))
	b.WriteString("\n\n")

	// Render table headers
	headerRow := ""
	for i, col := range m.table.Columns {
		header := col.Header
		if m.table.SortBy == i {
			if m.table.SortDesc {
				header += " ↓"
			} else {
				header += " ↑"
			}
		}

		cellStyle := lipgloss.NewStyle().
			Width(col.Width).
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#E6E6FA")).
			Padding(0, 1)

		headerRow += cellStyle.Render(header) + " "
	}
	b.WriteString(headerRow)
	b.WriteString("\n")

	// Render data rows
	pageData := m.table.GetPage(m.currentPage)
	for rowIndex, row := range pageData {
		rowStr := ""
		isSelected := rowIndex == m.selectedRow

		for colIndex, cell := range row.Cells {
			if colIndex >= len(m.table.Columns) {
				continue
			}

			col := m.table.Columns[colIndex]
			value := col.Formatter(cell.Value)

			cellStyle := lipgloss.NewStyle().
				Width(col.Width).
				Padding(0, 1)

			if isSelected {
				cellStyle = cellStyle.
					Background(lipgloss.Color("#7D4CDB")).
					Foreground(lipgloss.Color("#FFFFFF")).
					Bold(true)
			} else {
				cellStyle = cellStyle.
					Background(lipgloss.Color("#F8F8FF"))
			}

			rowStr += cellStyle.Render(value) + " "
		}

		b.WriteString(rowStr)
		b.WriteString("\n")
	}

	// Custom status bar
	b.WriteString("\n")
	statusStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#666666"))

	totalPages := m.table.GetTotalPages()
	status := fmt.Sprintf("Page %d of %d | %d total tasks | Use 1-6 to sort, h/l for pages, k/j for rows, q to quit",
		m.currentPage+1, totalPages, m.table.TotalRows)

	if m.table.SortBy >= 0 {
		sortDir := "ascending"
		if m.table.SortDesc {
			sortDir = "descending"
		}
		status += fmt.Sprintf(" | Sorted by %s (%s)",
			m.table.Columns[m.table.SortBy].Header, sortDir)
	}

	b.WriteString(statusStyle.Render(status))

	return b.String()
}

func main() {
	model := NewCustomModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
