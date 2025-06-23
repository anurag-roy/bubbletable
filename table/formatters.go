package table

import (
	"fmt"
	"time"
)

// DefaultFormatter is the default formatter that converts any value to string
func DefaultFormatter(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float64:
		return fmt.Sprintf("%.2f", v)
	case float32:
		return fmt.Sprintf("%.2f", v)
	case time.Time:
		return v.Format("2006-01-02")
	default:
		return fmt.Sprintf("%v", value)
	}
}

// CurrencyFormatter formats numeric values as currency
func CurrencyFormatter(value interface{}) string {
	switch v := value.(type) {
	case float64:
		return fmt.Sprintf("$%.2f", v)
	case float32:
		return fmt.Sprintf("$%.2f", v)
	case int:
		return fmt.Sprintf("$%d.00", v)
	case int64:
		return fmt.Sprintf("$%d.00", v)
	default:
		// Try to convert to float
		if str := fmt.Sprintf("%v", value); str != "" {
			return "$" + str
		}
		return "$0.00"
	}
}

// PercentFormatter formats numeric values as percentages
func PercentFormatter(value interface{}) string {
	switch v := value.(type) {
	case float64:
		return fmt.Sprintf("%.1f%%", v*100)
	case float32:
		return fmt.Sprintf("%.1f%%", v*100)
	default:
		return fmt.Sprintf("%v%%", value)
	}
}

// DateFormatter formats date values consistently
func DateFormatter(value interface{}) string {
	switch v := value.(type) {
	case time.Time:
		return v.Format("2006-01-02")
	case string:
		// Try to parse the string as a date
		if t, err := time.Parse("2006-01-02", v); err == nil {
			return t.Format("2006-01-02")
		}
		if t, err := time.Parse("01/02/2006", v); err == nil {
			return t.Format("2006-01-02")
		}
		if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			return t.Format("2006-01-02")
		}
		return v
	default:
		return fmt.Sprintf("%v", value)
	}
}

// TimeFormatter formats time values with time included
func TimeFormatter(value interface{}) string {
	switch v := value.(type) {
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	case string:
		// Try to parse and reformat
		if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
		return v
	default:
		return fmt.Sprintf("%v", value)
	}
}

// BooleanFormatter formats boolean values with custom true/false strings
func BooleanFormatter(trueStr, falseStr string) Formatter {
	return func(value interface{}) string {
		switch v := value.(type) {
		case bool:
			if v {
				return trueStr
			}
			return falseStr
		case string:
			if v == "true" || v == "1" || v == "yes" {
				return trueStr
			}
			return falseStr
		default:
			return falseStr
		}
	}
}

// NumberWithCommasFormatter formats large numbers with comma separators
func NumberWithCommasFormatter(value interface{}) string {
	switch v := value.(type) {
	case int:
		return addCommas(fmt.Sprintf("%d", v))
	case int64:
		return addCommas(fmt.Sprintf("%d", v))
	case float64:
		return addCommas(fmt.Sprintf("%.0f", v))
	case float32:
		return addCommas(fmt.Sprintf("%.0f", v))
	default:
		return fmt.Sprintf("%v", value)
	}
}

// addCommas adds comma separators to a numeric string
func addCommas(s string) string {
	if len(s) <= 3 {
		return s
	}

	// Handle negative numbers
	negative := false
	if s[0] == '-' {
		negative = true
		s = s[1:]
	}

	n := len(s)
	result := make([]byte, 0, n+n/3)

	for i, char := range []byte(s) {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, char)
	}

	if negative {
		return "-" + string(result)
	}
	return string(result)
}

// TruncateFormatter creates a formatter that truncates strings to a maximum length
func TruncateFormatter(maxLength int) Formatter {
	return func(value interface{}) string {
		str := fmt.Sprintf("%v", value)
		if len(str) <= maxLength {
			return str
		}
		if maxLength <= 3 {
			return str[:maxLength]
		}
		return str[:maxLength-3] + "..."
	}
}

// PrefixFormatter creates a formatter that adds a prefix to values
func PrefixFormatter(prefix string) Formatter {
	return func(value interface{}) string {
		return prefix + fmt.Sprintf("%v", value)
	}
}

// SuffixFormatter creates a formatter that adds a suffix to values
func SuffixFormatter(suffix string) Formatter {
	return func(value interface{}) string {
		return fmt.Sprintf("%v", value) + suffix
	}
}
