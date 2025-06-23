# TUI Data Table

A terminal-based interactive data table application built with Go and Charmbracelet's Bubble Tea framework.

![TUI Data Table Screenshot](screenshot.png)

## Features

- ✅ **Interactive TUI** - Beautiful terminal interface with responsive design
- ✅ **Data Model** - Robust table structures with multiple data types
- ✅ **Sample Data** - Built-in generators for testing (Employee, Product tables)
- ✅ **Pagination** - Navigate through large datasets with arrow keys or `h`/`l`/`n`/`p`
- ✅ **Sorting** - Sort columns by pressing number keys (1-6), with 3-state cycling (unsorted → ascending → descending)
- ✅ **Searching** - Real-time search with `/` key, press `ESC` to exit search mode
- ✅ **Row Selection** - Navigate rows with `↑`/`↓` or `j`/`k` keys
- ✅ **Dynamic Page Sizing** - Adjust page size with `+`/`-` keys, `r` to reset to optimal size
- ✅ **Responsive Layout** - Automatically adapts to terminal size changes
- 🔲 **CSV Import** - Load data from CSV files
- 🔲 **Configurable** - Customizable settings and key bindings

## Getting Started

### Prerequisites

- Go 1.21 or later

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```

### Running the Application

```bash
# Build and run
go build -o tui-data-table .
./tui-data-table

# Or run directly
go run .
```

### Controls

#### Navigation

- `↑`/`↓` or `j`/`k` - Navigate rows
- `←`/`→` or `h`/`l` - Previous/Next page
- `n` - Next page
- `p` - Previous page
- `Home`/`g` - Go to first page
- `End`/`G` - Go to last page

#### Sorting

- `1`-`6` - Sort by column (press again to reverse, press third time to clear sort)

#### Page Size

- `+`/`=` - Increase page size
- `-`/`_` - Decrease page size
- `r` - Reset to optimal page size

#### Search

- `/` - Enter search mode
- Type to search across all columns
- `ESC` - Exit search mode
- `Backspace` - Delete search characters

#### General

- `q`/`ESC`/`Ctrl+C` - Quit the application

## Project Structure

```
tui-data-table/
├── main.go           # Main application entry point
├── cmd/              # Command-line tools
├── internal/         # Private application code
│   ├── table/        # Table data structures and logic
│   └── ui/           # UI rendering and layout
├── pkg/              # Public library code
└── README.md         # This file
```

## Data Model

The application supports multiple data types with proper formatting and sorting:

- **String** - Text data with case-insensitive sorting
- **Integer** - Numeric values with proper numeric sorting
- **Float** - Decimal numbers with 2-decimal formatting
- **Date** - Date values with YYYY-MM-DD format
- **Boolean** - True/false values displayed as "true"/"false"

### Sample Data

Built-in sample data generators include:

- **Employee Table** - 50 employees with ID, Name, Department, Salary, Start Date, Active status
- **Product Table** - 100 products with SKU, Name, Category, Price, Stock, Available status

The default view shows the Employee table with realistic sample data including various departments, salary ranges, and hire dates.

## Tech Stack

- **Go** - Main programming language
- **Bubble Tea** - TUI framework for interactive terminal applications
- **Lip Gloss** - Terminal styling and layout
- **Charmbracelet** ecosystem for beautiful terminal UIs

## License

[MIT](LICENSE)
