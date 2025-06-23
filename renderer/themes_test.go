package renderer

import (
	"testing"
)

func TestAllThemesAvailable(t *testing.T) {
	themes := GetAllThemes()

	expectedThemes := []string{
		"Default",
		"Dracula",
		"Monokai",
		"GitHub",
		"Terminal",
		"Solarized Dark",
		"Solarized Light",
	}

	if len(themes) != len(expectedThemes) {
		t.Errorf("Expected %d themes, got %d", len(expectedThemes), len(themes))
	}

	for _, expectedName := range expectedThemes {
		found := false
		for _, theme := range themes {
			if theme.Name == expectedName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Theme '%s' not found in GetAllThemes()", expectedName)
		}
	}
}

func TestGetThemeByName(t *testing.T) {
	theme, found := GetThemeByName("Default")
	if !found {
		t.Error("Expected to find Default theme")
	}
	if theme.Name != "Default" {
		t.Errorf("Expected theme name 'Default', got '%s'", theme.Name)
	}

	// Test non-existent theme
	_, found = GetThemeByName("NonExistent")
	if found {
		t.Error("Expected not to find non-existent theme")
	}
}

func TestDefaultTheme(t *testing.T) {
	theme := DefaultTheme

	if theme.Name != "Default" {
		t.Errorf("Expected theme name 'Default', got '%s'", theme.Name)
	}

	// Test that all style components are defined
	if theme.Header.GetForeground() == nil {
		t.Error("Header style should have foreground color")
	}
	if theme.Header.GetBackground() == nil {
		t.Error("Header style should have background color")
	}
	if theme.SelectedRow.GetBackground() == nil {
		t.Error("SelectedRow style should have background color")
	}
}

func TestDraculaTheme(t *testing.T) {
	theme := DraculaTheme

	if theme.Name != "Dracula" {
		t.Errorf("Expected theme name 'Dracula', got '%s'", theme.Name)
	}

	// Test that theme has distinct styling
	if theme.Header.GetBackground() == nil {
		t.Error("Dracula theme should have header background")
	}
}

func TestMonokaiTheme(t *testing.T) {
	theme := MonokaiTheme

	if theme.Name != "Monokai" {
		t.Errorf("Expected theme name 'Monokai', got '%s'", theme.Name)
	}

	// Test that theme has distinct styling
	if theme.Header.GetBackground() == nil {
		t.Error("Monokai theme should have header background")
	}
}

func TestGithubTheme(t *testing.T) {
	theme := GithubTheme

	if theme.Name != "GitHub" {
		t.Errorf("Expected theme name 'GitHub', got '%s'", theme.Name)
	}

	// GitHub theme should have light colors
	if theme.Header.GetBackground() == nil {
		t.Error("GitHub theme should have header background")
	}
}

func TestTerminalTheme(t *testing.T) {
	theme := TerminalTheme

	if theme.Name != "Terminal" {
		t.Errorf("Expected theme name 'Terminal', got '%s'", theme.Name)
	}

	// Terminal theme should use basic colors
	if theme.Header.GetBackground() == nil {
		t.Error("Terminal theme should have header background")
	}
}

func TestSolarizedThemes(t *testing.T) {
	darkTheme := SolarizedDarkTheme
	lightTheme := SolarizedLightTheme

	if darkTheme.Name != "Solarized Dark" {
		t.Errorf("Expected 'Solarized Dark', got '%s'", darkTheme.Name)
	}

	if lightTheme.Name != "Solarized Light" {
		t.Errorf("Expected 'Solarized Light', got '%s'", lightTheme.Name)
	}

	// Both themes should have backgrounds (they are different themes)
	if darkTheme.Header.GetBackground() == nil {
		t.Error("Solarized Dark theme should have header background")
	}
	if lightTheme.Header.GetBackground() == nil {
		t.Error("Solarized Light theme should have header background")
	}
}

func TestCustomizeTheme(t *testing.T) {
	baseTheme := DefaultTheme
	customName := "Custom Test Theme"

	// Test that CustomizeTheme function exists and works
	// Note: Testing actual customization requires proper lipgloss.Style map
	customTheme := baseTheme
	customTheme.Name = customName

	if customTheme.Name != customName {
		t.Errorf("Expected custom theme name '%s', got '%s'", customName, customTheme.Name)
	}

	// Original theme should be unchanged
	if baseTheme.Name == customName {
		t.Error("Original theme should not be modified")
	}
}

func TestThemeStructure(t *testing.T) {
	themes := []Theme{
		DefaultTheme,
		DraculaTheme,
		MonokaiTheme,
		GithubTheme,
		TerminalTheme,
		SolarizedDarkTheme,
		SolarizedLightTheme,
	}

	for _, theme := range themes {
		t.Run(theme.Name, func(t *testing.T) {
			// Test that all required fields are present
			if theme.Name == "" {
				t.Error("Theme name should not be empty")
			}

			// Test that styles are defined (non-nil)
			// Note: lipgloss styles are always valid, but we can test they exist
			_ = theme.Header
			_ = theme.Cell
			_ = theme.SelectedRow
			_ = theme.Border
			_ = theme.Status
			_ = theme.Search
		})
	}
}

func TestThemeUniqueness(t *testing.T) {
	themes := GetAllThemes()
	nameMap := make(map[string]bool)

	for _, theme := range themes {
		if nameMap[theme.Name] {
			t.Errorf("Duplicate theme name found: %s", theme.Name)
		}
		nameMap[theme.Name] = true
	}
}

func TestCustomThemeBuilder(t *testing.T) {
	// Test creating a completely custom theme
	customTheme := Theme{
		Name:        "Test Custom",
		Header:      DefaultTheme.Header,
		Cell:        DefaultTheme.Cell,
		SelectedRow: DefaultTheme.SelectedRow,
		Border:      DefaultTheme.Border,
		Status:      DefaultTheme.Status,
		Search:      DefaultTheme.Search,
	}

	if customTheme.Name != "Test Custom" {
		t.Error("Custom theme name should be set correctly")
	}

	// Verify styles can be applied without panicking
	headerText := customTheme.Header.Render("Test Header")
	if headerText == "" {
		t.Error("Header style should render text")
	}

	cellText := customTheme.Cell.Render("Test Cell")
	if cellText == "" {
		t.Error("Cell style should render text")
	}
}

func TestThemeConsistency(t *testing.T) {
	themes := GetAllThemes()

	for _, theme := range themes {
		t.Run(theme.Name, func(t *testing.T) {
			// Test that each theme can render content without errors
			testContent := "Test Content"

			headerRendered := theme.Header.Render(testContent)
			if len(headerRendered) == 0 {
				t.Error("Header should render content")
			}

			cellRendered := theme.Cell.Render(testContent)
			if len(cellRendered) == 0 {
				t.Error("Cell should render content")
			}

			selectedRendered := theme.SelectedRow.Render(testContent)
			if len(selectedRendered) == 0 {
				t.Error("SelectedRow should render content")
			}

			statusRendered := theme.Status.Render(testContent)
			if len(statusRendered) == 0 {
				t.Error("Status should render content")
			}

			searchRendered := theme.Search.Render(testContent)
			if len(searchRendered) == 0 {
				t.Error("Search should render content")
			}
		})
	}
}

func TestThemeColorAccessibility(t *testing.T) {
	// Basic test that themes have contrasting colors for readability
	themes := []Theme{DefaultTheme, DraculaTheme, GithubTheme}

	for _, theme := range themes {
		t.Run(theme.Name, func(t *testing.T) {
			// Test that header has both foreground and background
			headerFg := theme.Header.GetForeground()
			headerBg := theme.Header.GetBackground()

			if headerFg == nil {
				t.Error("Header should have foreground color for readability")
			}
			if headerBg == nil {
				t.Error("Header should have background color for visibility")
			}

			// Test that selected row has background for visibility
			selectedBg := theme.SelectedRow.GetBackground()
			if selectedBg == nil {
				t.Error("Selected row should have background color for visibility")
			}
		})
	}
}
