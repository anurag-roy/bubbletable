package main

import (
	"log"

	"github.com/anurag-roy/bubbletable/components"
	"github.com/anurag-roy/bubbletable/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Employee represents a sample employee struct with table tags
type Employee struct {
	ID         int     `table:"ID,sortable,width:5"`
	Name       string  `table:"Name,sortable,width:20"`
	Department string  `table:"Department,sortable,width:15"`
	Salary     float64 `table:"Salary,sortable,width:12,format:currency"`
	StartDate  string  `table:"Start Date,sortable,width:12,format:date"`
	Active     bool    `table:"Active,sortable,width:8"`
}

func main() {
	// Sample data
	employees := []Employee{
		{1, "Alice Johnson", "Engineering", 75000.0, "2021-01-15", true},
		{2, "Bob Smith", "Marketing", 65000.0, "2020-03-20", true},
		{3, "Charlie Brown", "Sales", 55000.0, "2019-11-10", false},
		{4, "Diana Prince", "Engineering", 85000.0, "2021-06-01", true},
		{5, "Edward Norton", "Finance", 70000.0, "2020-08-15", true},
		{6, "Fiona Apple", "Design", 68000.0, "2021-02-28", true},
		{7, "George Lucas", "Engineering", 95000.0, "2018-12-01", true},
		{8, "Helen Troy", "HR", 62000.0, "2020-09-15", true},
		{9, "Ivan Drago", "Sales", 58000.0, "2019-07-20", true},
		{10, "Julia Roberts", "Marketing", 72000.0, "2021-04-10", true},
		{11, "Kevin Hart", "Engineering", 78000.0, "2020-11-25", true},
		{12, "Luna Lovegood", "Design", 64000.0, "2021-01-30", true},
		{13, "Mike Tyson", "Operations", 66000.0, "2020-05-12", true},
		{14, "Nancy Drew", "Legal", 88000.0, "2019-03-08", true},
		{15, "Oscar Wilde", "Marketing", 69000.0, "2020-10-22", false},
	}

	// Create table model with fluent API
	tableModel := components.NewTable(employees).
		WithPageSize(10).
		WithSorting(true).
		WithSearch(true).
		WithOnSelect(func(row table.Row) {
			// Optional: Handle row selection
			// log.Printf("Selected row: %+v", row)
		})

	// Create and run the Bubble Tea program
	p := tea.NewProgram(tableModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
