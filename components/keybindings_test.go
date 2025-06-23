package components

import (
	"testing"
)

func TestDefaultKeyBindings(t *testing.T) {
	kb := DefaultKeyBindings()

	// Test navigation keys
	if !kb.IsUp("up") {
		t.Error("Default bindings should recognize 'up' as up key")
	}
	if !kb.IsUp("k") {
		t.Error("Default bindings should recognize 'k' as up key")
	}

	if !kb.IsDown("down") {
		t.Error("Default bindings should recognize 'down' as down key")
	}
	if !kb.IsDown("j") {
		t.Error("Default bindings should recognize 'j' as down key")
	}

	if !kb.IsLeft("left") {
		t.Error("Default bindings should recognize 'left' as left key")
	}
	if !kb.IsLeft("h") {
		t.Error("Default bindings should recognize 'h' as left key")
	}

	if !kb.IsRight("right") {
		t.Error("Default bindings should recognize 'right' as right key")
	}
	if !kb.IsRight("l") {
		t.Error("Default bindings should recognize 'l' as right key")
	}

	// Test quit keys
	if !kb.IsQuit("q") {
		t.Error("Default bindings should recognize 'q' as quit key")
	}
	if !kb.IsQuit("esc") {
		t.Error("Default bindings should recognize 'esc' as quit key")
	}
	if !kb.IsQuit("ctrl+c") {
		t.Error("Default bindings should recognize 'ctrl+c' as quit key")
	}

	// Test search key
	if !kb.IsSearch("/") {
		t.Error("Default bindings should recognize '/' as search key")
	}

	// Test help key
	if !kb.IsHelp("?") {
		t.Error("Default bindings should recognize '?' as help key")
	}
}

func TestVimKeyBindings(t *testing.T) {
	kb := VimKeyBindings()

	// Test vim-specific keys
	if !kb.IsUp("k") {
		t.Error("Vim bindings should recognize 'k' as up key")
	}
	if !kb.IsDown("j") {
		t.Error("Vim bindings should recognize 'j' as down key")
	}
	if !kb.IsLeft("h") {
		t.Error("Vim bindings should recognize 'h' as left key")
	}
	if !kb.IsRight("l") {
		t.Error("Vim bindings should recognize 'l' as right key")
	}

	// Test vim page navigation
	if !kb.IsPageUp("ctrl+u") {
		t.Error("Vim bindings should recognize 'ctrl+u' as page up key")
	}
	if !kb.IsPageDown("ctrl+d") {
		t.Error("Vim bindings should recognize 'ctrl+d' as page down key")
	}

	// Test vim home/end
	if !kb.IsHome("gg") {
		t.Error("Vim bindings should recognize 'gg' as home key")
	}
	if !kb.IsEnd("G") {
		t.Error("Vim bindings should recognize 'G' as end key")
	}

	// Should still support standard keys
	if !kb.IsUp("up") {
		t.Error("Vim bindings should still support 'up' key")
	}
}

func TestEmacsKeyBindings(t *testing.T) {
	kb := EmacsKeyBindings()

	// Test emacs-specific keys
	if !kb.IsUp("ctrl+p") {
		t.Error("Emacs bindings should recognize 'ctrl+p' as up key")
	}
	if !kb.IsDown("ctrl+n") {
		t.Error("Emacs bindings should recognize 'ctrl+n' as down key")
	}
	if !kb.IsLeft("ctrl+b") {
		t.Error("Emacs bindings should recognize 'ctrl+b' as left key")
	}
	if !kb.IsRight("ctrl+f") {
		t.Error("Emacs bindings should recognize 'ctrl+f' as right key")
	}

	// Test emacs home/end
	if !kb.IsHome("ctrl+a") {
		t.Error("Emacs bindings should recognize 'ctrl+a' as home key")
	}
	if !kb.IsEnd("ctrl+e") {
		t.Error("Emacs bindings should recognize 'ctrl+e' as end key")
	}

	// Test emacs search
	if !kb.IsSearch("ctrl+s") {
		t.Error("Emacs bindings should recognize 'ctrl+s' as search key")
	}

	// Should still support standard keys
	if !kb.IsUp("up") {
		t.Error("Emacs bindings should still support 'up' key")
	}
}

func TestSortColumns(t *testing.T) {
	kb := DefaultKeyBindings()

	// Test sort column mapping
	if kb.GetSortColumn("1") != 0 {
		t.Error("Key '1' should map to column 0")
	}
	if kb.GetSortColumn("2") != 1 {
		t.Error("Key '2' should map to column 1")
	}
	if kb.GetSortColumn("9") != 8 {
		t.Error("Key '9' should map to column 8")
	}

	// Test invalid sort key
	if kb.GetSortColumn("0") != -1 {
		t.Error("Key '0' should return -1 (invalid)")
	}
	if kb.GetSortColumn("a") != -1 {
		t.Error("Key 'a' should return -1 (invalid)")
	}
}

func TestPageSizeKeys(t *testing.T) {
	kb := DefaultKeyBindings()

	// Test page size adjustment keys
	if !kb.IsPageSizeUp("+") {
		t.Error("Default bindings should recognize '+' as page size up key")
	}
	if !kb.IsPageSizeUp("=") {
		t.Error("Default bindings should recognize '=' as page size up key")
	}

	if !kb.IsPageSizeDown("-") {
		t.Error("Default bindings should recognize '-' as page size down key")
	}
	if !kb.IsPageSizeDown("_") {
		t.Error("Default bindings should recognize '_' as page size down key")
	}
}

func TestPageNavigationKeys(t *testing.T) {
	kb := DefaultKeyBindings()

	// Test page navigation
	if !kb.IsPageUp("pageup") {
		t.Error("Default bindings should recognize 'pageup' as page up key")
	}
	if !kb.IsPageUp("ctrl+b") {
		t.Error("Default bindings should recognize 'ctrl+b' as page up key")
	}

	if !kb.IsPageDown("pagedown") {
		t.Error("Default bindings should recognize 'pagedown' as page down key")
	}
	if !kb.IsPageDown("ctrl+f") {
		t.Error("Default bindings should recognize 'ctrl+f' as page down key")
	}

	// Test home/end
	if !kb.IsHome("home") {
		t.Error("Default bindings should recognize 'home' as home key")
	}
	if !kb.IsHome("g") {
		t.Error("Default bindings should recognize 'g' as home key")
	}

	if !kb.IsEnd("end") {
		t.Error("Default bindings should recognize 'end' as end key")
	}
	if !kb.IsEnd("G") {
		t.Error("Default bindings should recognize 'G' as end key")
	}
}

func TestRefreshAndClearKeys(t *testing.T) {
	kb := DefaultKeyBindings()

	// Test refresh key
	if !kb.IsRefresh("r") {
		t.Error("Default bindings should recognize 'r' as refresh key")
	}

	// Test clear sort key
	if !kb.IsClearSort("c") {
		t.Error("Default bindings should recognize 'c' as clear sort key")
	}
}

func TestMatchesKey(t *testing.T) {
	// Test the internal matchesKey function through public methods
	testKeys := []string{"test1", "test2", "test3"}

	// This tests the internal logic via public methods
	customKb := KeyBindings{
		Up: testKeys,
	}

	if !customKb.IsUp("test1") {
		t.Error("Should match first key in list")
	}
	if !customKb.IsUp("test2") {
		t.Error("Should match second key in list")
	}
	if !customKb.IsUp("test3") {
		t.Error("Should match third key in list")
	}
	if customKb.IsUp("test4") {
		t.Error("Should not match key not in list")
	}
}

func TestAllKeyBindingMethods(t *testing.T) {
	// Test that all methods work without panicking
	kb := DefaultKeyBindings()
	testMethods := []struct {
		name   string
		method func(string) bool
		key    string
	}{
		{"IsUp", kb.IsUp, "up"},
		{"IsDown", kb.IsDown, "down"},
		{"IsLeft", kb.IsLeft, "left"},
		{"IsRight", kb.IsRight, "right"},
		{"IsPageUp", kb.IsPageUp, "pageup"},
		{"IsPageDown", kb.IsPageDown, "pagedown"},
		{"IsHome", kb.IsHome, "home"},
		{"IsEnd", kb.IsEnd, "end"},
		{"IsSearch", kb.IsSearch, "/"},
		{"IsQuit", kb.IsQuit, "q"},
		{"IsHelp", kb.IsHelp, "?"},
		{"IsRefresh", kb.IsRefresh, "r"},
		{"IsPageSizeUp", kb.IsPageSizeUp, "+"},
		{"IsPageSizeDown", kb.IsPageSizeDown, "-"},
		{"IsClearSort", kb.IsClearSort, "c"},
	}

	for _, test := range testMethods {
		t.Run(test.name, func(t *testing.T) {
			// Should not panic and should return true for expected keys
			result := test.method(test.key)
			if !result {
				t.Errorf("%s should return true for key %q", test.name, test.key)
			}

			// Should return false for unexpected keys
			result = test.method("invalid_key_xyz")
			if result {
				t.Errorf("%s should return false for invalid key", test.name)
			}
		})
	}
}

func TestCustomKeyBindings(t *testing.T) {
	// Test creating custom key bindings
	customKb := KeyBindings{
		Up:    []string{"w", "up"},
		Down:  []string{"s", "down"},
		Left:  []string{"a", "left"},
		Right: []string{"d", "right"},
		Quit:  []string{"q", "escape"},
		Sort: map[int][]string{
			0: {"1", "!"},
			1: {"2", "@"},
		},
	}

	// Test custom navigation
	if !customKb.IsUp("w") {
		t.Error("Custom bindings should recognize 'w' as up key")
	}
	if !customKb.IsDown("s") {
		t.Error("Custom bindings should recognize 's' as down key")
	}
	if !customKb.IsLeft("a") {
		t.Error("Custom bindings should recognize 'a' as left key")
	}
	if !customKb.IsRight("d") {
		t.Error("Custom bindings should recognize 'd' as right key")
	}

	// Test custom quit
	if !customKb.IsQuit("escape") {
		t.Error("Custom bindings should recognize 'escape' as quit key")
	}

	// Test custom sort keys
	if customKb.GetSortColumn("!") != 0 {
		t.Error("Custom sort key '!' should map to column 0")
	}
	if customKb.GetSortColumn("@") != 1 {
		t.Error("Custom sort key '@' should map to column 1")
	}
}

func TestEmptyKeyBindings(t *testing.T) {
	// Test behavior with empty key bindings
	emptyKb := KeyBindings{}

	// Should return false for all keys
	if emptyKb.IsUp("up") {
		t.Error("Empty bindings should return false for all keys")
	}
	if emptyKb.IsQuit("q") {
		t.Error("Empty bindings should return false for all keys")
	}
	if emptyKb.GetSortColumn("1") != -1 {
		t.Error("Empty bindings should return -1 for sort columns")
	}
}
