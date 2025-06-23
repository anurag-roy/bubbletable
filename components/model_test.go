package components

import (
	"testing"

	"github.com/anurag-roy/bubbletable/renderer"
	"github.com/anurag-roy/bubbletable/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Test data structures
type TestEmployee struct {
	ID   int    `table:"ID,sortable,width:5"`
	Name string `table:"Name,sortable,width:20"`
}

func TestNewTable(t *testing.T) {
	employees := []TestEmployee{
		{1, "Alice"},
		{2, "Bob"},
	}

	model := NewTable(employees)

	if model == nil {
		t.Fatal("NewTable should return non-nil model")
	}

	if model.table == nil {
		t.Error("Model should have a table instance")
	}

	if model.renderer == nil {
		t.Error("Model should have a renderer instance")
	}

	if model.pageSize != 10 {
		t.Errorf("Expected default page size 10, got %d", model.pageSize)
	}

	if model.theme.Name != renderer.DefaultTheme.Name {
		t.Errorf("Expected default theme, got %s", model.theme.Name)
	}
}

func TestNewTableFromInterface(t *testing.T) {
	data := []interface{}{
		TestEmployee{1, "Alice"},
		TestEmployee{2, "Bob"},
	}

	model := NewTableFromInterface(data)

	if model == nil {
		t.Fatal("NewTableFromInterface should return non-nil model")
	}

	if model.table == nil {
		t.Error("Model should have a table instance")
	}
}

func TestNewTableWithColumns(t *testing.T) {
	columns := []table.Column{
		*table.NewColumn("id", "ID").WithType(table.Integer),
		*table.NewColumn("name", "Name").WithType(table.String),
	}

	data := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}

	model := NewTableWithColumns(data, columns)

	if model == nil {
		t.Fatal("NewTableWithColumns should return non-nil model")
	}

	if len(model.table.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(model.table.Columns))
	}
}

func TestBuilderPattern(t *testing.T) {
	employees := []TestEmployee{
		{1, "Alice"},
		{2, "Bob"},
	}

	theme := renderer.DraculaTheme
	model := NewTable(employees).
		WithPageSize(20).
		WithTheme(&theme).
		WithKeyBindings(VimKeyBindings()).
		WithSorting(true).
		WithSearch(true)

	if model.pageSize != 20 {
		t.Errorf("Expected page size 20, got %d", model.pageSize)
	}

	if model.theme.Name != "Dracula" {
		t.Errorf("Expected Dracula theme, got %s", model.theme.Name)
	}

	// Test that methods return the same instance for chaining
	originalModel := NewTable(employees)
	chainedModel := originalModel.WithPageSize(15)

	if originalModel != chainedModel {
		t.Error("Builder methods should return the same instance for chaining")
	}
}

func TestWithCallbacks(t *testing.T) {
	employees := []TestEmployee{
		{1, "Alice"},
	}

	var selectCalled bool
	var sortCalled bool
	var searchCalled bool
	var refreshCalled bool

	model := NewTable(employees).
		WithOnSelect(func(row table.Row) {
			selectCalled = true
		}).
		WithOnSort(func(columnIndex int, desc bool) {
			sortCalled = true
		}).
		WithOnSearch(func(term string) {
			searchCalled = true
		}).
		WithOnRefresh(func() {
			refreshCalled = true
		})

	// Test that callbacks are set
	if model.onSelect == nil {
		t.Error("OnSelect callback should be set")
	}
	if model.onSort == nil {
		t.Error("OnSort callback should be set")
	}
	if model.onSearch == nil {
		t.Error("OnSearch callback should be set")
	}
	if model.onRefresh == nil {
		t.Error("OnRefresh callback should be set")
	}

	// Test callback execution
	model.onSelect(table.Row{})
	if !selectCalled {
		t.Error("OnSelect callback should have been called")
	}

	model.onSort(0, false)
	if !sortCalled {
		t.Error("OnSort callback should have been called")
	}

	model.onSearch("test")
	if !searchCalled {
		t.Error("OnSearch callback should have been called")
	}

	model.onRefresh()
	if !refreshCalled {
		t.Error("OnRefresh callback should have been called")
	}
}

func TestInit(t *testing.T) {
	employees := []TestEmployee{
		{1, "Alice"},
	}

	model := NewTable(employees)
	cmd := model.Init()

	// Init should return nil for this model
	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

func TestUpdateWindowSize(t *testing.T) {
	employees := []TestEmployee{
		{1, "Alice"},
	}

	model := NewTable(employees)

	// Send window size message
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedModel, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("WindowSizeMsg should not return a command")
	}

	tableModel := updatedModel.(*TableModel)
	if tableModel.width != 120 {
		t.Errorf("Expected width 120, got %d", tableModel.width)
	}

	if tableModel.height != 40 {
		t.Errorf("Expected height 40, got %d", tableModel.height)
	}

	if !tableModel.ready {
		t.Error("Model should be ready after receiving window size")
	}
}

func TestGetMethods(t *testing.T) {
	employees := []TestEmployee{
		{1, "Alice"},
		{2, "Bob"},
	}

	model := NewTable(employees)

	// Test GetTable
	tbl := model.GetTable()
	if tbl != model.table {
		t.Error("GetTable should return the internal table")
	}

	// Test GetCurrentTable (without filtering)
	currentTbl := model.GetCurrentTable()
	if currentTbl != model.table {
		t.Error("GetCurrentTable should return the main table when not filtered")
	}

	// Test GetSelectedRow
	model.selectedRow = 0
	row, exists := model.GetSelectedRow()
	if !exists {
		t.Error("GetSelectedRow should return true when valid row is selected")
	}
	if len(row.Cells) == 0 {
		t.Error("Selected row should have cells")
	}

	// Test with invalid selected row
	model.selectedRow = 999
	_, exists = model.GetSelectedRow()
	if exists {
		t.Error("GetSelectedRow should return false for invalid row index")
	}
}

func TestSetData(t *testing.T) {
	model := NewTable([]TestEmployee{})

	newEmployees := []TestEmployee{
		{1, "Alice"},
		{2, "Bob"},
		{3, "Charlie"},
	}

	err := model.SetData(newEmployees)
	if err != nil {
		t.Errorf("SetData failed: %v", err)
	}

	if len(model.table.Rows) != 3 {
		t.Errorf("Expected 3 rows after SetData, got %d", len(model.table.Rows))
	}
}

func TestRefreshData(t *testing.T) {
	originalEmployees := []TestEmployee{
		{1, "Alice"},
	}

	model := NewTable(originalEmployees)

	newEmployees := []TestEmployee{
		{1, "Alice Updated"},
		{2, "Bob"},
	}

	err := model.RefreshData(newEmployees)
	if err != nil {
		t.Errorf("RefreshData failed: %v", err)
	}

	if len(model.table.Rows) != 2 {
		t.Errorf("Expected 2 rows after RefreshData, got %d", len(model.table.Rows))
	}

	// Check that the data was actually updated
	if model.table.Rows[0].Cells[1].Value != "Alice Updated" {
		t.Error("Data should be refreshed with new values")
	}
}

func TestView(t *testing.T) {
	employees := []TestEmployee{
		{1, "Alice"},
		{2, "Bob"},
	}

	model := NewTable(employees)

	// Set up model as ready
	model.ready = true
	model.width = 80
	model.height = 24

	view := model.View()

	if view == "" {
		t.Error("View should not be empty")
	}

	// Test view when not ready
	model.ready = false
	view = model.View()

	if view == "" {
		t.Error("View should show something even when not ready")
	}
}

func TestSearchMode(t *testing.T) {
	employees := []TestEmployee{
		{1, "Alice"},
		{2, "Bob"},
	}

	model := NewTable(employees)
	model.ready = true

	// Test entering search mode
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	updatedModel, _ := model.Update(keyMsg)

	tableModel := updatedModel.(*TableModel)
	if !tableModel.searchMode {
		t.Error("Should enter search mode when '/' is pressed")
	}

	// Test search input
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}}
	updatedModel, _ = tableModel.Update(keyMsg)

	tableModel = updatedModel.(*TableModel)
	if tableModel.searchTerm != "A" {
		t.Errorf("Expected search term 'A', got '%s'", tableModel.searchTerm)
	}

	// Test exiting search mode
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = tableModel.Update(keyMsg)

	tableModel = updatedModel.(*TableModel)
	if tableModel.searchMode {
		t.Error("Should exit search mode when Enter is pressed")
	}
}

func TestPagination(t *testing.T) {
	// Create enough data for multiple pages
	employees := make([]TestEmployee, 25)
	for i := 0; i < 25; i++ {
		employees[i] = TestEmployee{i + 1, "Employee " + string(rune('A'+i%26))}
	}

	model := NewTable(employees).WithPageSize(10)
	model.ready = true

	// Test navigation to next page
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}} // right/next page
	updatedModel, _ := model.Update(keyMsg)

	tableModel := updatedModel.(*TableModel)
	if tableModel.currentPage != 1 {
		t.Errorf("Expected current page 1, got %d", tableModel.currentPage)
	}

	// Test navigation to previous page
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}} // left/prev page
	updatedModel, _ = tableModel.Update(keyMsg)

	tableModel = updatedModel.(*TableModel)
	if tableModel.currentPage != 0 {
		t.Errorf("Expected current page 0, got %d", tableModel.currentPage)
	}
}

func TestSortingIntegration(t *testing.T) {
	employees := []TestEmployee{
		{3, "Charlie"},
		{1, "Alice"},
		{2, "Bob"},
	}

	model := NewTable(employees)
	model.ready = true

	// Test sorting by first column (ID)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	updatedModel, _ := model.Update(keyMsg)

	tableModel := updatedModel.(*TableModel)

	// Check if table is sorted
	if tableModel.table.SortBy != 0 {
		t.Errorf("Expected table to be sorted by column 0, got %d", tableModel.table.SortBy)
	}

	// Check if data is actually sorted
	page := tableModel.table.GetPage(0)
	if len(page) > 0 && page[0].Cells[0].Value != 1 {
		t.Error("First row should have ID 1 after sorting")
	}
}

func TestEmptyData(t *testing.T) {
	model := NewTable([]TestEmployee{})

	if model.table.TotalRows != 0 {
		t.Errorf("Expected 0 total rows for empty data, got %d", model.table.TotalRows)
	}

	// Should not panic when getting selected row with empty data
	_, exists := model.GetSelectedRow()
	if exists {
		t.Error("GetSelectedRow should return false for empty data")
	}

	// View should not panic with empty data
	model.ready = true
	view := model.View()
	if view == "" {
		t.Error("View should not be empty even with no data")
	}
}

// TestWithTheme tests theme setting
func TestWithTheme(t *testing.T) {
	employees := []TestEmployee{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	theme := renderer.DraculaTheme
	model := NewTable(employees).WithPageSize(20).WithTheme(&theme)

	if model.theme.Name != "Dracula" {
		t.Errorf("Expected theme name 'Dracula', got '%s'", model.theme.Name)
	}
}

func TestHelpMenu(t *testing.T) {
	employees := []TestEmployee{
		{1, "Alice"},
		{2, "Bob"},
	}

	model := NewTable(employees)
	model.ready = true

	// Test initial state - help should not be showing
	if model.showHelp {
		t.Error("Help should not be showing initially")
	}

	// Test entering help mode
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updatedModel, _ := model.Update(keyMsg)
	tableModel := updatedModel.(*TableModel)

	if !tableModel.showHelp {
		t.Error("Help should be showing after pressing '?'")
	}

	// Test that help screen is rendered
	view := tableModel.View()
	if view == "Loading..." {
		t.Error("Help view should not show loading message")
	}

	// Help view should contain help text
	expectedHelpText := "Table Navigation Help"
	if !contains(view, expectedHelpText) {
		t.Errorf("Help view should contain '%s', got: %s", expectedHelpText, view)
	}

	// Test exiting help mode with '?' key
	updatedModel2, _ := tableModel.Update(keyMsg)
	tableModel2 := updatedModel2.(*TableModel)

	if tableModel2.showHelp {
		t.Error("Help should not be showing after pressing '?' again")
	}

	// Test that normal table view is restored
	normalView := tableModel2.View()
	if contains(normalView, expectedHelpText) {
		t.Error("Normal view should not contain help text")
	}
}

func TestHelpMenuKeyHandling(t *testing.T) {
	employees := []TestEmployee{
		{1, "Alice"},
		{2, "Bob"},
	}

	model := NewTable(employees)
	model.ready = true

	// Enter help mode
	helpKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updatedModel, _ := model.Update(helpKeyMsg)
	tableModel := updatedModel.(*TableModel)

	if !tableModel.showHelp {
		t.Error("Help should be showing after pressing '?'")
	}

	// Test that other keys are ignored in help mode
	navigationKeys := []tea.KeyMsg{
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
		{Type: tea.KeyLeft},
		{Type: tea.KeyRight},
		{Type: tea.KeyRunes, Runes: []rune{'j'}},
		{Type: tea.KeyRunes, Runes: []rune{'k'}},
		{Type: tea.KeyRunes, Runes: []rune{'h'}},
		{Type: tea.KeyRunes, Runes: []rune{'l'}},
		{Type: tea.KeyRunes, Runes: []rune{'1'}},
		{Type: tea.KeyRunes, Runes: []rune{'/'}},
	}

	// Store original state
	originalPage := tableModel.currentPage
	originalRow := tableModel.selectedRow
	originalSearchMode := tableModel.searchMode

	for _, keyMsg := range navigationKeys {
		updatedModel, _ := tableModel.Update(keyMsg)
		tableModel = updatedModel.(*TableModel)

		// Help should still be showing
		if !tableModel.showHelp {
			t.Errorf("Help should still be showing after pressing key: %v", keyMsg)
		}

		// State should not change
		if tableModel.currentPage != originalPage {
			t.Errorf("Current page should not change in help mode, was %d, now %d", originalPage, tableModel.currentPage)
		}
		if tableModel.selectedRow != originalRow {
			t.Errorf("Selected row should not change in help mode, was %d, now %d", originalRow, tableModel.selectedRow)
		}
		if tableModel.searchMode != originalSearchMode {
			t.Errorf("Search mode should not change in help mode, was %t, now %t", originalSearchMode, tableModel.searchMode)
		}
	}

	// Test that quit key still works in help mode
	quitKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, _ = tableModel.Update(quitKeyMsg)
	tableModel = updatedModel.(*TableModel)

	// Help should still be showing (quit is handled by the caller)
	if !tableModel.showHelp {
		t.Error("Help should still be showing after quit key (quit is handled by caller)")
	}

	// Test that escape key still works in help mode
	escKeyMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ = tableModel.Update(escKeyMsg)
	tableModel = updatedModel.(*TableModel)

	// Help should still be showing (quit is handled by the caller)
	if !tableModel.showHelp {
		t.Error("Help should still be showing after escape key (quit is handled by caller)")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
