// Package table provides the headless core functionality for BubbleTable.
//
// This package contains the core table data structures and operations without any
// UI dependencies, making it suitable for custom implementations and testing.
//
// # Core Types
//
//   - Table: The main table data structure that manages rows, columns, and state
//   - Column: Defines column metadata, formatting, and behavior
//   - Row: Represents a single row of data
//   - ColumnType: Enumeration of supported data types
//
// # Key Features
//
//   - Generic data support with reflection-based struct parsing
//   - Struct tag configuration for column metadata
//   - Multi-type sorting with custom comparison functions
//   - Case-insensitive filtering and search
//   - Efficient pagination for large datasets
//   - Built-in formatters for common data types
//   - Custom formatters and renderers support
//
// # Usage Examples
//
// Creating a table with struct data:
//
//	type Person struct {
//	    Name string `table:"Full Name,sortable,width:25"`
//	    Age  int    `table:"Age,sortable,width:10"`
//	}
//
//	people := []Person{
//	    {"Alice", 30},
//	    {"Bob", 25},
//	}
//
//	tbl := table.New().WithData(people)
//
// Creating a table with map data and custom columns:
//
//	columns := []table.Column{
//	    *table.NewColumn("name", "Name").WithType(table.String).WithWidth(20),
//	    *table.NewColumn("salary", "Salary").WithType(table.Float).WithFormatter(table.CurrencyFormatter),
//	}
//
//	data := []map[string]interface{}{
//	    {"name": "Alice", "salary": 75000.0},
//	    {"name": "Bob", "salary": 65000.0},
//	}
//
//	tbl := table.New().WithColumns(columns).WithData(data)
//
// # Struct Tags
//
// The table package supports comprehensive struct tag configuration:
//
//	type Example struct {
//	    ID          int     `table:"ID,sortable,width:5"`
//	    Name        string  `table:"Full Name,sortable,searchable,width:25"`
//	    Price       float64 `table:"Price,sortable,width:12,format:currency"`
//	    InStock     bool    `table:"Available,sortable,width:10"`
//	    Description string  `table:"Description,!sortable,!searchable,width:30"`
//	}
//
// Supported tag options:
//   - sortable/!sortable: Enable/disable sorting
//   - searchable/!searchable: Enable/disable search filtering
//   - width:N: Set column width in characters
//   - format:type: Apply built-in formatters (currency, date, percent, etc.)
//
// # Formatters
//
// Built-in formatters handle common data presentation needs:
//   - CurrencyFormatter: Formats numbers as currency ($1,234.56)
//   - PercentFormatter: Formats decimals as percentages (0.15 → 15%)
//   - DateFormatter: Formats time.Time as dates
//   - TimeFormatter: Formats time.Time as times
//   - BooleanFormatter: Formats booleans as Yes/No
//   - NumberWithCommasFormatter: Adds thousand separators
//
// Custom formatters can be created by implementing the FormatterFunc type.
//
// # Performance
//
// The table package is optimized for performance with large datasets:
//   - Efficient pagination avoids loading all data into memory
//   - Sorting uses optimized comparison functions for each data type
//   - Filtering operates on paginated data for better responsiveness
//   - Reflection-based operations are cached to minimize overhead
//
// # Thread Safety
//
// The table package is not thread-safe by design, as it's intended for use
// in single-threaded TUI applications. If concurrent access is needed,
// external synchronization should be implemented.
package table
