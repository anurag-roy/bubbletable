package table

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DataType represents the type of data in a column
type DataType int

const (
	String DataType = iota
	Integer
	Float
	Date
	Boolean
)

// Formatter is a function that formats a value for display
type Formatter func(value interface{}) string

// CellRenderer is a function that renders a cell value with styling information
type CellRenderer func(value interface{}, selected bool) string

// Accessor is a function that extracts a value from a data row
type Accessor func(data interface{}) interface{}

// Column represents a table column with metadata and behavior
type Column struct {
	Key        string
	Header     string
	Type       DataType
	Width      int
	Sortable   bool
	Searchable bool
	Formatter  Formatter
	Renderer   CellRenderer
	Accessor   Accessor
}

// NewColumn creates a new column with the given key and header
func NewColumn(key, header string) *Column {
	return &Column{
		Key:        key,
		Header:     header,
		Type:       String,
		Width:      15,
		Sortable:   true,
		Searchable: true,
		Formatter:  DefaultFormatter,
		Renderer:   nil,
		Accessor:   nil,
	}
}

// WithType sets the column data type
func (c *Column) WithType(dataType DataType) *Column {
	c.Type = dataType
	return c
}

// WithWidth sets the column display width
func (c *Column) WithWidth(width int) *Column {
	c.Width = width
	return c
}

// WithSortable sets whether the column is sortable
func (c *Column) WithSortable(sortable bool) *Column {
	c.Sortable = sortable
	return c
}

// WithSearchable sets whether the column is searchable
func (c *Column) WithSearchable(searchable bool) *Column {
	c.Searchable = searchable
	return c
}

// WithFormatter sets a custom formatter for the column
func (c *Column) WithFormatter(formatter Formatter) *Column {
	c.Formatter = formatter
	return c
}

// WithRenderer sets a custom renderer for the column
func (c *Column) WithRenderer(renderer CellRenderer) *Column {
	c.Renderer = renderer
	return c
}

// WithAccessor sets a custom accessor function for the column
func (c *Column) WithAccessor(accessor Accessor) *Column {
	c.Accessor = accessor
	return c
}

// Cell represents a single cell value with type information
type Cell struct {
	Value interface{}
	Type  DataType
}

// Row represents a table row
type Row struct {
	ID    int
	Cells []Cell
	Data  interface{} // Original data for custom accessors
}

// Table represents the complete table data structure
type Table struct {
	Columns       []Column
	Rows          []Row
	UnsortedOrder []Row // Store original row order for unsort functionality
	SortBy        int   // Column index for sorting (-1 if not sorted)
	SortDesc      bool  // Sort direction (true for descending)
	PageSize      int
	TotalRows     int
	originalData  []interface{} // Store original data for re-processing
}

// New creates a new empty table
func New() *Table {
	return &Table{
		Columns:       make([]Column, 0),
		Rows:          make([]Row, 0),
		UnsortedOrder: make([]Row, 0),
		SortBy:        -1,
		SortDesc:      false,
		PageSize:      10,
		TotalRows:     0,
		originalData:  make([]interface{}, 0),
	}
}

// NewWithColumns creates a new table with predefined columns
func NewWithColumns(columns []Column) *Table {
	table := New()
	table.Columns = columns
	return table
}

// WithColumns sets the table columns (builder pattern)
func (t *Table) WithColumns(columns []Column) *Table {
	t.Columns = columns
	return t
}

// WithPageSize sets the page size (builder pattern)
func (t *Table) WithPageSize(pageSize int) *Table {
	t.PageSize = pageSize
	return t
}

// WithData sets the table data from a slice of any type (builder pattern)
func (t *Table) WithData(data interface{}) *Table {
	t.SetData(data)
	return t
}

// SetData populates the table with data from a slice of any type
func (t *Table) SetData(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	// Clear existing data
	t.Rows = make([]Row, 0)
	t.UnsortedOrder = make([]Row, 0)
	t.TotalRows = 0
	t.originalData = make([]interface{}, 0)

	// If no columns are defined, try to infer them from the data
	if len(t.Columns) == 0 && v.Len() > 0 {
		firstItem := v.Index(0).Interface()
		columns, err := t.inferColumnsFromStruct(firstItem)
		if err != nil {
			return err
		}
		t.Columns = columns
	}

	// Process each item in the slice
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		t.originalData = append(t.originalData, item)

		err := t.addRowFromData(item, i)
		if err != nil {
			return fmt.Errorf("error processing row %d: %v", i, err)
		}
	}

	return nil
}

// addRowFromData adds a row from arbitrary data
func (t *Table) addRowFromData(data interface{}, id int) error {
	cells := make([]Cell, len(t.Columns))

	for i, col := range t.Columns {
		var value interface{}
		var err error

		if col.Accessor != nil {
			// Use custom accessor
			value = col.Accessor(data)
		} else {
			// Try to extract value based on column key
			value, err = t.extractValueFromData(data, col.Key)
			if err != nil {
				value = ""
			}
		}

		cells[i] = Cell{
			Value: value,
			Type:  col.Type,
		}
	}

	row := Row{
		ID:    id,
		Cells: cells,
		Data:  data,
	}

	t.Rows = append(t.Rows, row)
	t.UnsortedOrder = append(t.UnsortedOrder, row)
	t.TotalRows++

	return nil
}

// extractValueFromData extracts a value from data using reflection
func (t *Table) extractValueFromData(data interface{}, key string) (interface{}, error) {
	v := reflect.ValueOf(data)

	// Handle maps
	if v.Kind() == reflect.Map {
		mapValue := v.MapIndex(reflect.ValueOf(key))
		if mapValue.IsValid() {
			return mapValue.Interface(), nil
		}
		return nil, fmt.Errorf("key %s not found in map", key)
	}

	// Handle structs
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		field := v.FieldByName(key)
		if field.IsValid() {
			return field.Interface(), nil
		}

		// Try case-insensitive search
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			fieldType := t.Field(i)
			if strings.EqualFold(fieldType.Name, key) {
				return v.Field(i).Interface(), nil
			}
		}

		return nil, fmt.Errorf("field %s not found in struct", key)
	}

	return nil, fmt.Errorf("cannot extract value from type %T", data)
}

// inferColumnsFromStruct infers columns from a struct using reflection and struct tags
func (t *Table) inferColumnsFromStruct(data interface{}) ([]Column, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Map {
		return t.inferColumnsFromMap(data)
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("cannot infer columns from type %T", data)
	}

	var columns []Column
	structType := v.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		col := Column{
			Key:        field.Name,
			Header:     field.Name,
			Sortable:   true,
			Searchable: true,
			Formatter:  DefaultFormatter,
		}

		// Parse struct tag for configuration
		if tag := field.Tag.Get("table"); tag != "" {
			col = t.parseStructTag(col, tag)
		}

		// Infer type from Go type
		col.Type = t.inferDataType(field.Type)

		// Set default width based on type
		col.Width = t.getDefaultWidth(col.Type)

		columns = append(columns, col)
	}

	return columns, nil
}

// inferColumnsFromMap infers columns from a map
func (t *Table) inferColumnsFromMap(data interface{}) ([]Column, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Map {
		return nil, fmt.Errorf("expected map, got %T", data)
	}

	var columns []Column
	for _, key := range v.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())
		value := v.MapIndex(key)

		col := Column{
			Key:        keyStr,
			Header:     keyStr,
			Sortable:   true,
			Searchable: true,
			Formatter:  DefaultFormatter,
			Type:       t.inferDataType(value.Type()),
		}

		col.Width = t.getDefaultWidth(col.Type)
		columns = append(columns, col)
	}

	return columns, nil
}

// parseStructTag parses struct tag for column configuration
func (t *Table) parseStructTag(col Column, tag string) Column {
	parts := strings.Split(tag, ",")

	// First part is the header name
	if len(parts) > 0 && parts[0] != "" {
		col.Header = parts[0]
	}

	// Parse additional options
	for i := 1; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])

		switch {
		case part == "sortable":
			col.Sortable = true
		case part == "!sortable":
			col.Sortable = false
		case part == "searchable":
			col.Searchable = true
		case part == "!searchable":
			col.Searchable = false
		case strings.HasPrefix(part, "width:"):
			if width, err := strconv.Atoi(strings.TrimPrefix(part, "width:")); err == nil {
				col.Width = width
			}
		case strings.HasPrefix(part, "format:"):
			format := strings.TrimPrefix(part, "format:")
			col.Formatter = t.getFormatterByName(format)
		}
	}

	return col
}

// inferDataType infers DataType from Go reflect.Type
func (t *Table) inferDataType(goType reflect.Type) DataType {
	switch goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Integer
	case reflect.Float32, reflect.Float64:
		return Float
	case reflect.Bool:
		return Boolean
	case reflect.String:
		return String
	default:
		// Check if it's a time.Time
		if goType == reflect.TypeOf(time.Time{}) {
			return Date
		}
		return String
	}
}

// getDefaultWidth returns default width for a data type
func (t *Table) getDefaultWidth(dataType DataType) int {
	switch dataType {
	case Integer:
		return 8
	case Float:
		return 10
	case Boolean:
		return 8
	case Date:
		return 12
	default:
		return 15
	}
}

// getFormatterByName returns a formatter by name
func (t *Table) getFormatterByName(name string) Formatter {
	switch name {
	case "currency":
		return CurrencyFormatter
	case "date":
		return DateFormatter
	case "percent":
		return PercentFormatter
	default:
		return DefaultFormatter
	}
}

// AddRow adds a new row to the table with explicit values
func (t *Table) AddRow(values ...interface{}) error {
	if len(values) != len(t.Columns) {
		return fmt.Errorf("expected %d values, got %d", len(t.Columns), len(values))
	}

	cells := make([]Cell, len(values))
	for i, value := range values {
		cells[i] = Cell{
			Value: value,
			Type:  t.Columns[i].Type,
		}
	}

	row := Row{
		ID:    t.TotalRows,
		Cells: cells,
		Data:  values, // Store the raw values as data
	}

	t.Rows = append(t.Rows, row)
	t.UnsortedOrder = append(t.UnsortedOrder, row)
	t.TotalRows++
	return nil
}

// GetPage returns a slice of rows for the given page number (0-indexed)
func (t *Table) GetPage(pageNum int) []Row {
	start := pageNum * t.PageSize
	end := start + t.PageSize

	if start >= len(t.Rows) {
		return []Row{}
	}

	if end > len(t.Rows) {
		end = len(t.Rows)
	}

	return t.Rows[start:end]
}

// GetTotalPages returns the total number of pages
func (t *Table) GetTotalPages() int {
	if t.PageSize <= 0 {
		return 1
	}
	if len(t.Rows) == 0 {
		return 1
	}
	return (len(t.Rows) + t.PageSize - 1) / t.PageSize
}

// SortByColumn sorts the table by the specified column
func (t *Table) SortByColumn(columnIndex int, descending bool) error {
	if columnIndex < 0 || columnIndex >= len(t.Columns) {
		return fmt.Errorf("invalid column index: %d", columnIndex)
	}

	if !t.Columns[columnIndex].Sortable {
		return fmt.Errorf("column %s is not sortable", t.Columns[columnIndex].Header)
	}

	t.SortBy = columnIndex
	t.SortDesc = descending

	sort.Slice(t.Rows, func(i, j int) bool {
		cellI := t.Rows[i].Cells[columnIndex]
		cellJ := t.Rows[j].Cells[columnIndex]

		result := compareCells(cellI, cellJ)
		if descending {
			return result > 0
		}
		return result < 0
	})

	return nil
}

// ClearSort clears any active sorting and restores original order
func (t *Table) ClearSort() {
	t.SortBy = -1
	t.SortDesc = false
	// Restore original order
	t.Rows = make([]Row, len(t.UnsortedOrder))
	copy(t.Rows, t.UnsortedOrder)
}

// Filter returns a new table with rows matching the search term
func (t *Table) Filter(searchTerm string) *Table {
	if searchTerm == "" {
		return t
	}

	filtered := NewWithColumns(t.Columns)
	filtered.PageSize = t.PageSize
	searchTerm = strings.ToLower(searchTerm)

	for _, row := range t.Rows {
		if t.rowMatchesSearch(row, searchTerm) {
			filtered.Rows = append(filtered.Rows, row)
			filtered.UnsortedOrder = append(filtered.UnsortedOrder, row)
			filtered.TotalRows++
		}
	}

	// Preserve sort state from original table
	filtered.SortBy = t.SortBy
	filtered.SortDesc = t.SortDesc

	return filtered
}

// rowMatchesSearch checks if a row contains the search term in any searchable cell
func (t *Table) rowMatchesSearch(row Row, searchTerm string) bool {
	for i, cell := range row.Cells {
		if i < len(t.Columns) && t.Columns[i].Searchable {
			cellStr := strings.ToLower(t.formatCellValue(cell, i))
			if strings.Contains(cellStr, searchTerm) {
				return true
			}
		}
	}
	return false
}

// formatCellValue formats a cell value using the column's formatter
func (t *Table) formatCellValue(cell Cell, columnIndex int) string {
	if columnIndex < len(t.Columns) && t.Columns[columnIndex].Formatter != nil {
		return t.Columns[columnIndex].Formatter(cell.Value)
	}
	return DefaultFormatter(cell.Value)
}

// GetCellValue returns the formatted value of a cell
func (t *Table) GetCellValue(rowIndex, columnIndex int) string {
	if rowIndex < 0 || rowIndex >= len(t.Rows) {
		return ""
	}
	if columnIndex < 0 || columnIndex >= len(t.Columns) {
		return ""
	}

	cell := t.Rows[rowIndex].Cells[columnIndex]
	return t.formatCellValue(cell, columnIndex)
}

// GetColumnNames returns the headers of all columns
func (t *Table) GetColumnNames() []string {
	names := make([]string, len(t.Columns))
	for i, col := range t.Columns {
		names[i] = col.Header
	}
	return names
}

// compareCells compares two cells for sorting purposes
func compareCells(a, b Cell) int {
	switch a.Type {
	case String:
		aStr := fmt.Sprintf("%v", a.Value)
		bStr := fmt.Sprintf("%v", b.Value)
		return strings.Compare(strings.ToLower(aStr), strings.ToLower(bStr))

	case Integer:
		aInt, _ := strconv.Atoi(fmt.Sprintf("%v", a.Value))
		bInt, _ := strconv.Atoi(fmt.Sprintf("%v", b.Value))
		return aInt - bInt

	case Float:
		aFloat, _ := strconv.ParseFloat(fmt.Sprintf("%v", a.Value), 64)
		bFloat, _ := strconv.ParseFloat(fmt.Sprintf("%v", b.Value), 64)
		if aFloat < bFloat {
			return -1
		} else if aFloat > bFloat {
			return 1
		}
		return 0

	case Boolean:
		aBool := fmt.Sprintf("%v", a.Value) == "true"
		bBool := fmt.Sprintf("%v", b.Value) == "true"
		if aBool == bBool {
			return 0
		} else if aBool {
			return 1
		}
		return -1

	case Date:
		aTime, aErr := parseDate(a.Value)
		bTime, bErr := parseDate(b.Value)
		if aErr != nil || bErr != nil {
			// Fall back to string comparison if date parsing fails
			aStr := fmt.Sprintf("%v", a.Value)
			bStr := fmt.Sprintf("%v", b.Value)
			return strings.Compare(aStr, bStr)
		}
		if aTime.Before(bTime) {
			return -1
		} else if aTime.After(bTime) {
			return 1
		}
		return 0

	default:
		aStr := fmt.Sprintf("%v", a.Value)
		bStr := fmt.Sprintf("%v", b.Value)
		return strings.Compare(strings.ToLower(aStr), strings.ToLower(bStr))
	}
}

// parseDate parses various date formats
func parseDate(value interface{}) (time.Time, error) {
	str := fmt.Sprintf("%v", value)

	// Try common date formats
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"01/02/2006",
		"01-02-2006",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", str)
}
