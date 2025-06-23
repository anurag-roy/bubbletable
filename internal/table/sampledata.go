package table

import (
	"fmt"
	"math/rand"
	"time"
)

// SampleDataGenerator provides methods to generate test data
type SampleDataGenerator struct {
	rand *rand.Rand
}

// NewSampleDataGenerator creates a new sample data generator
func NewSampleDataGenerator() *SampleDataGenerator {
	return &SampleDataGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateEmployeeTable creates a sample employee table
func (g *SampleDataGenerator) GenerateEmployeeTable() *Table {
	columns := []Column{
		{Name: "ID", Type: Integer, Width: 5, Sortable: true},
		{Name: "Name", Type: String, Width: 20, Sortable: true},
		{Name: "Department", Type: String, Width: 15, Sortable: true},
		{Name: "Salary", Type: Float, Width: 10, Sortable: true},
		{Name: "Start Date", Type: Date, Width: 12, Sortable: true},
		{Name: "Active", Type: Boolean, Width: 8, Sortable: true},
	}

	table := NewTable(columns)
	table.PageSize = 10

	// Sample data
	names := []string{
		"Alice Johnson", "Bob Smith", "Charlie Brown", "Diana Prince",
		"Edward Norton", "Fiona Apple", "George Lucas", "Helen Troy",
		"Ivan Drago", "Julia Roberts", "Kevin Hart", "Luna Lovegood",
		"Mike Tyson", "Nancy Drew", "Oscar Wilde", "Penny Lane",
		"Quincy Jones", "Rachel Green", "Steve Jobs", "Tina Turner",
		"Uma Thurman", "Victor Hugo", "Wendy Darling", "Xavier Woods",
		"Yoda Master", "Zoe Saldana", "Aaron Paul", "Bella Swan",
		"Clark Kent", "Daisy Miller", "Ethan Hunt", "Faith Hill",
		"Goku Son", "Hermione Granger", "Indiana Jones", "Jessica Alba",
	}

	departments := []string{
		"Engineering", "Marketing", "Sales", "HR", "Finance",
		"Design", "Operations", "Customer Success", "Legal", "R&D",
	}

	// Generate sample employees
	for i := 0; i < 50; i++ {
		id := i + 1
		name := names[i%len(names)]
		department := departments[g.rand.Intn(len(departments))]
		salary := 45000.0 + g.rand.Float64()*85000.0 // $45k to $130k

		// Random start date within last 5 years
		startDate := time.Now().AddDate(-g.rand.Intn(5), -g.rand.Intn(12), -g.rand.Intn(28))

		active := g.rand.Float64() > 0.1 // 90% active employees

		err := table.AddRow(id, name, department, salary, startDate.Format("2006-01-02"), active)
		if err != nil {
			fmt.Printf("Error adding row: %v\n", err)
		}
	}

	return table
}

// GenerateProductTable creates a sample product table
func (g *SampleDataGenerator) GenerateProductTable() *Table {
	columns := []Column{
		{Name: "SKU", Type: String, Width: 10, Sortable: true},
		{Name: "Product Name", Type: String, Width: 25, Sortable: true},
		{Name: "Category", Type: String, Width: 15, Sortable: true},
		{Name: "Price", Type: Float, Width: 8, Sortable: true},
		{Name: "Stock", Type: Integer, Width: 8, Sortable: true},
		{Name: "Available", Type: Boolean, Width: 10, Sortable: true},
	}

	table := NewTable(columns)
	table.PageSize = 15

	products := []string{
		"Wireless Headphones", "Smart Watch", "Laptop Stand", "USB-C Cable",
		"Mechanical Keyboard", "Gaming Mouse", "Monitor", "Desk Lamp",
		"Phone Case", "Tablet Stand", "Webcam", "Microphone",
		"Speaker Set", "Power Bank", "Charging Pad", "Bluetooth Adapter",
		"Cable Organizer", "Desk Pad", "Ergonomic Chair", "Standing Desk",
		"Notebook", "Pen Set", "Sticky Notes", "Highlighters",
		"Coffee Mug", "Water Bottle", "Lunch Box", "Backpack",
		"Travel Adapter", "Portable SSD", "Memory Card", "Screen Protector",
	}

	categories := []string{
		"Electronics", "Accessories", "Office", "Furniture", "Stationery",
		"Kitchen", "Travel", "Storage", "Audio", "Computing",
	}

	for i := 0; i < 100; i++ {
		sku := fmt.Sprintf("SKU-%04d", i+1)
		product := products[i%len(products)]
		if i >= len(products) {
			product = fmt.Sprintf("%s v%d", product, (i/len(products))+1)
		}
		category := categories[g.rand.Intn(len(categories))]
		price := 9.99 + g.rand.Float64()*490.0 // $9.99 to $499.99
		stock := g.rand.Intn(200)
		available := stock > 0

		err := table.AddRow(sku, product, category, price, stock, available)
		if err != nil {
			fmt.Printf("Error adding row: %v\n", err)
		}
	}

	return table
}

// GenerateFinancialTable creates a sample financial data table
func (g *SampleDataGenerator) GenerateFinancialTable() *Table {
	columns := []Column{
		{Name: "Date", Type: Date, Width: 12, Sortable: true},
		{Name: "Symbol", Type: String, Width: 8, Sortable: true},
		{Name: "Open", Type: Float, Width: 10, Sortable: true},
		{Name: "High", Type: Float, Width: 10, Sortable: true},
		{Name: "Low", Type: Float, Width: 10, Sortable: true},
		{Name: "Close", Type: Float, Width: 10, Sortable: true},
		{Name: "Volume", Type: Integer, Width: 12, Sortable: true},
	}

	table := NewTable(columns)
	table.PageSize = 20

	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "META", "NVDA", "NFLX"}

	// Generate 30 days of data for each symbol
	for _, symbol := range symbols {
		basePrice := 100.0 + g.rand.Float64()*400.0 // $100 to $500

		for i := 0; i < 30; i++ {
			date := time.Now().AddDate(0, 0, -i)

			// Simulate price movement
			change := (g.rand.Float64() - 0.5) * 20.0 // ±$10 change
			open := basePrice + change

			high := open + g.rand.Float64()*10.0
			low := open - g.rand.Float64()*10.0
			close := low + g.rand.Float64()*(high-low)

			volume := 1000000 + g.rand.Intn(10000000) // 1M to 11M shares

			basePrice = close // Use close as next day's base

			err := table.AddRow(
				date.Format("2006-01-02"),
				symbol,
				open,
				high,
				low,
				close,
				volume,
			)
			if err != nil {
				fmt.Printf("Error adding row: %v\n", err)
			}
		}
	}

	return table
}

// GenerateCustomTable creates a customizable table with specified parameters
func (g *SampleDataGenerator) GenerateCustomTable(rows int, columns []Column) *Table {
	table := NewTable(columns)
	table.PageSize = 15

	for i := 0; i < rows; i++ {
		values := make([]interface{}, len(columns))

		for j, col := range columns {
			values[j] = g.generateValueForType(col.Type, i)
		}

		err := table.AddRow(values...)
		if err != nil {
			fmt.Printf("Error adding row: %v\n", err)
		}
	}

	return table
}

// generateValueForType generates a sample value for a given data type
func (g *SampleDataGenerator) generateValueForType(dataType DataType, index int) interface{} {
	switch dataType {
	case String:
		options := []string{
			"Alpha", "Beta", "Gamma", "Delta", "Epsilon", "Zeta",
			"Eta", "Theta", "Iota", "Kappa", "Lambda", "Mu",
		}
		return fmt.Sprintf("%s-%d", options[g.rand.Intn(len(options))], index+1)

	case Integer:
		return g.rand.Intn(1000) + 1

	case Float:
		return g.rand.Float64() * 1000.0

	case Date:
		days := g.rand.Intn(365 * 2) // Random date within last 2 years
		date := time.Now().AddDate(0, 0, -days)
		return date.Format("2006-01-02")

	case Boolean:
		return g.rand.Float64() > 0.5

	default:
		return fmt.Sprintf("Value-%d", index+1)
	}
}

// GenerateSampleTable creates a sample table by name
func (g *SampleDataGenerator) GenerateSampleTable(name string) *Table {
	switch name {
	case "employees":
		return g.GenerateEmployeeTable()
	case "products":
		return g.GenerateProductTable()
	case "financial":
		return g.GenerateFinancialTable()
	default:
		return g.GenerateEmployeeTable() // Default fallback
	}
}

// GetSampleTableNames returns a list of available sample table types
func GetSampleTableNames() []string {
	return []string{
		"employees",
		"products",
		"financial",
	}
}
