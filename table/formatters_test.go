package table

import (
	"testing"
	"time"
)

func TestDefaultFormatter(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"nil", nil, ""},
		{"string", "hello", "hello"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"float64", 123.456, "123.46"},
		{"float32", float32(123.456), "123.46"},
		{"int", 42, "42"},
		{"time", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), "2023-01-15"},
		{"struct", struct{ Name string }{"test"}, "{test}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultFormatter(tt.input)
			if result != tt.expected {
				t.Errorf("DefaultFormatter(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCurrencyFormatter(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"float64", 99.99, "$99.99"},
		{"float32", float32(123.45), "$123.45"},
		{"int", 100, "$100.00"},
		{"int64", int64(250), "$250.00"},
		{"zero", 0.0, "$0.00"},
		{"negative", -50.25, "$-50.25"},
		{"large number", 1234567.89, "$1234567.89"},
		{"string fallback", "invalid", "$invalid"},
		{"nil fallback", nil, "$<nil>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CurrencyFormatter(tt.input)
			if result != tt.expected {
				t.Errorf("CurrencyFormatter(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPercentFormatter(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"float64 decimal", 0.25, "25.0%"},
		{"float32 decimal", float32(0.75), "75.0%"},
		{"float64 whole", 1.0, "100.0%"},
		{"zero", 0.0, "0.0%"},
		{"over 100%", 1.5, "150.0%"},
		{"small decimal", 0.001, "0.1%"},
		{"fallback", "50", "50%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PercentFormatter(tt.input)
			if result != tt.expected {
				t.Errorf("PercentFormatter(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDateFormatter(t *testing.T) {
	testDate := time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"time.Time", testDate, "2023-12-25"},
		{"ISO date string", "2023-01-15", "2023-01-15"},
		{"US date string", "01/15/2023", "2023-01-15"},
		{"datetime string", "2023-01-15 10:30:00", "2023-01-15"},
		{"invalid date string", "not a date", "not a date"},
		{"empty string", "", ""},
		{"other type", 12345, "12345"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DateFormatter(tt.input)
			if result != tt.expected {
				t.Errorf("DateFormatter(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTimeFormatter(t *testing.T) {
	testTime := time.Date(2023, 12, 25, 14, 30, 45, 0, time.UTC)

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"time.Time", testTime, "2023-12-25 14:30:45"},
		{"datetime string", "2023-01-15 10:30:00", "2023-01-15 10:30:00"},
		{"invalid time string", "not a time", "not a time"},
		{"other type", 12345, "12345"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimeFormatter(tt.input)
			if result != tt.expected {
				t.Errorf("TimeFormatter(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBooleanFormatter(t *testing.T) {
	formatter := BooleanFormatter("✅ Yes", "❌ No")

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"bool true", true, "✅ Yes"},
		{"bool false", false, "❌ No"},
		{"string true", "true", "✅ Yes"},
		{"string 1", "1", "✅ Yes"},
		{"string yes", "yes", "✅ Yes"},
		{"string false", "false", "❌ No"},
		{"string 0", "0", "❌ No"},
		{"string no", "no", "❌ No"},
		{"other string", "maybe", "❌ No"},
		{"int", 123, "❌ No"},
		{"nil", nil, "❌ No"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter(tt.input)
			if result != tt.expected {
				t.Errorf("BooleanFormatter(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNumberWithCommasFormatter(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"small int", 123, "123"},
		{"int with commas", 1234, "1,234"},
		{"large int", 1234567, "1,234,567"},
		{"int64", int64(9876543210), "9,876,543,210"},
		{"float64", 12345.67, "12,346"},        // Rounded
		{"float32", float32(9876.54), "9,877"}, // Rounded
		{"zero", 0, "0"},
		{"negative", -1234567, "-1,234,567"},
		{"fallback", "text", "text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NumberWithCommasFormatter(tt.input)
			if result != tt.expected {
				t.Errorf("NumberWithCommasFormatter(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAddCommas(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"short number", "123", "123"},
		{"four digits", "1234", "1,234"},
		{"seven digits", "1234567", "1,234,567"},
		{"ten digits", "1234567890", "1,234,567,890"},
		{"negative", "-1234567", "-1,234,567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addCommas(tt.input)
			if result != tt.expected {
				t.Errorf("addCommas(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTruncateFormatter(t *testing.T) {
	formatter := TruncateFormatter(10)

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"short string", "hello", "hello"},
		{"exact length", "1234567890", "1234567890"},
		{"long string", "this is a very long string", "this is..."},
		{"very short limit", "hello", "hello"}, // Test with different formatter
		{"number", 1234567890123, "1234567..."},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter(tt.input)
			if result != tt.expected {
				t.Errorf("TruncateFormatter(10)(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}

	// Test very short truncate
	shortFormatter := TruncateFormatter(3)
	result := shortFormatter("hello")
	if result != "hel" {
		t.Errorf("TruncateFormatter(3)(\"hello\") = %q, expected \"hel\"", result)
	}
}

func TestPrefixFormatter(t *testing.T) {
	formatter := PrefixFormatter("👤 ")

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "Alice", "👤 Alice"},
		{"number", 123, "👤 123"},
		{"empty", "", "👤 "},
		{"nil", nil, "👤 <nil>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter(tt.input)
			if result != tt.expected {
				t.Errorf("PrefixFormatter(\"👤 \")(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSuffixFormatter(t *testing.T) {
	formatter := SuffixFormatter(" units")

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "100", "100 units"},
		{"number", 50, "50 units"},
		{"empty", "", " units"},
		{"nil", nil, "<nil> units"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter(tt.input)
			if result != tt.expected {
				t.Errorf("SuffixFormatter(\" units\")(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatterChaining(t *testing.T) {
	// Test combining formatters (though they don't chain directly, we can test usage patterns)

	// Test BooleanFormatter with custom strings
	activeFormatter := BooleanFormatter("Active", "Inactive")
	if activeFormatter(true) != "Active" {
		t.Error("Custom boolean formatter failed for true")
	}
	if activeFormatter(false) != "Inactive" {
		t.Error("Custom boolean formatter failed for false")
	}

	// Test TruncateFormatter with different lengths
	longTruncator := TruncateFormatter(20)
	shortTruncator := TruncateFormatter(5)

	longText := "This is a very long text that should be truncated"

	longResult := longTruncator(longText)
	shortResult := shortTruncator(longText)

	if len(longResult) > 20 {
		t.Errorf("Long truncator should limit to 20 chars, got %d", len(longResult))
	}

	if len(shortResult) > 5 {
		t.Errorf("Short truncator should limit to 5 chars, got %d", len(shortResult))
	}

	if longResult == shortResult {
		t.Error("Different truncate lengths should produce different results")
	}
}

func TestEdgeCases(t *testing.T) {
	// Test with extremely large numbers
	hugeNumber := int64(9223372036854775807) // Max int64
	result := NumberWithCommasFormatter(hugeNumber)
	expected := "9,223,372,036,854,775,807"
	if result != expected {
		t.Errorf("NumberWithCommasFormatter(max int64) = %q, expected %q", result, expected)
	}

	// Test currency with very small amounts
	smallAmount := 0.001
	result = CurrencyFormatter(smallAmount)
	if result != "$0.00" {
		t.Errorf("CurrencyFormatter(0.001) = %q, expected \"$0.00\"", result)
	}

	// Test percent with zero
	result = PercentFormatter(0)
	if result != "0%" {
		t.Errorf("PercentFormatter(0) = %q, expected \"0%%\"", result)
	}

	// Test truncate with zero length (edge case)
	zeroTruncator := TruncateFormatter(0)
	result = zeroTruncator("hello")
	if result != "" {
		t.Errorf("TruncateFormatter(0)(\"hello\") = %q, expected \"\"", result)
	}
}
