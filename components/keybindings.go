package components

// KeyBindings represents configurable key bindings for table interactions
type KeyBindings struct {
	Up           []string
	Down         []string
	Left         []string
	Right        []string
	PageUp       []string
	PageDown     []string
	Home         []string
	End          []string
	Search       []string
	Sort         map[int][]string // Column index to keys
	Quit         []string
	Help         []string
	Refresh      []string
	PageSizeUp   []string
	PageSizeDown []string
	ResetPage    []string
	ClearSort    []string
}

// DefaultKeyBindings returns the default key bindings
func DefaultKeyBindings() KeyBindings {
	return KeyBindings{
		Up:       []string{"up", "k"},
		Down:     []string{"down", "j"},
		Left:     []string{"left", "h", "p"},
		Right:    []string{"right", "l", "n"},
		PageUp:   []string{"pageup", "ctrl+b"},
		PageDown: []string{"pagedown", "ctrl+f"},
		Home:     []string{"home", "g"},
		End:      []string{"end", "G"},
		Search:   []string{"/"},
		Sort: map[int][]string{
			0: {"1"},
			1: {"2"},
			2: {"3"},
			3: {"4"},
			4: {"5"},
			5: {"6"},
			6: {"7"},
			7: {"8"},
			8: {"9"},
		},
		Quit:         []string{"q", "esc", "ctrl+c"},
		Help:         []string{"?", "h"},
		Refresh:      []string{"r"},
		PageSizeUp:   []string{"+", "="},
		PageSizeDown: []string{"-", "_"},
		ResetPage:    []string{"r"},
		ClearSort:    []string{"c"},
	}
}

// VimKeyBindings returns vim-style key bindings
func VimKeyBindings() KeyBindings {
	bindings := DefaultKeyBindings()
	bindings.Up = []string{"k", "up"}
	bindings.Down = []string{"j", "down"}
	bindings.Left = []string{"h", "left"}
	bindings.Right = []string{"l", "right"}
	bindings.Home = []string{"gg", "home"}
	bindings.End = []string{"G", "end"}
	bindings.PageUp = []string{"ctrl+u", "pageup"}
	bindings.PageDown = []string{"ctrl+d", "pagedown"}
	return bindings
}

// EmacsKeyBindings returns emacs-style key bindings
func EmacsKeyBindings() KeyBindings {
	bindings := DefaultKeyBindings()
	bindings.Up = []string{"ctrl+p", "up"}
	bindings.Down = []string{"ctrl+n", "down"}
	bindings.Left = []string{"ctrl+b", "left"}
	bindings.Right = []string{"ctrl+f", "right"}
	bindings.Home = []string{"ctrl+a", "home"}
	bindings.End = []string{"ctrl+e", "end"}
	bindings.Search = []string{"ctrl+s", "/"}
	return bindings
}

// matchesKey checks if a key string matches any of the configured keys
func (kb KeyBindings) matchesKey(key string, keys []string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

// IsUp checks if the key is an up navigation key
func (kb KeyBindings) IsUp(key string) bool {
	return kb.matchesKey(key, kb.Up)
}

// IsDown checks if the key is a down navigation key
func (kb KeyBindings) IsDown(key string) bool {
	return kb.matchesKey(key, kb.Down)
}

// IsLeft checks if the key is a left navigation key
func (kb KeyBindings) IsLeft(key string) bool {
	return kb.matchesKey(key, kb.Left)
}

// IsRight checks if the key is a right navigation key
func (kb KeyBindings) IsRight(key string) bool {
	return kb.matchesKey(key, kb.Right)
}

// IsPageUp checks if the key is a page up key
func (kb KeyBindings) IsPageUp(key string) bool {
	return kb.matchesKey(key, kb.PageUp)
}

// IsPageDown checks if the key is a page down key
func (kb KeyBindings) IsPageDown(key string) bool {
	return kb.matchesKey(key, kb.PageDown)
}

// IsHome checks if the key is a home key
func (kb KeyBindings) IsHome(key string) bool {
	return kb.matchesKey(key, kb.Home)
}

// IsEnd checks if the key is an end key
func (kb KeyBindings) IsEnd(key string) bool {
	return kb.matchesKey(key, kb.End)
}

// IsSearch checks if the key is a search key
func (kb KeyBindings) IsSearch(key string) bool {
	return kb.matchesKey(key, kb.Search)
}

// IsQuit checks if the key is a quit key
func (kb KeyBindings) IsQuit(key string) bool {
	return kb.matchesKey(key, kb.Quit)
}

// IsHelp checks if the key is a help key
func (kb KeyBindings) IsHelp(key string) bool {
	return kb.matchesKey(key, kb.Help)
}

// IsRefresh checks if the key is a refresh key
func (kb KeyBindings) IsRefresh(key string) bool {
	return kb.matchesKey(key, kb.Refresh)
}

// IsPageSizeUp checks if the key increases page size
func (kb KeyBindings) IsPageSizeUp(key string) bool {
	return kb.matchesKey(key, kb.PageSizeUp)
}

// IsPageSizeDown checks if the key decreases page size
func (kb KeyBindings) IsPageSizeDown(key string) bool {
	return kb.matchesKey(key, kb.PageSizeDown)
}

// IsResetPage checks if the key resets to optimal page size
func (kb KeyBindings) IsResetPage(key string) bool {
	return kb.matchesKey(key, kb.ResetPage)
}

// IsClearSort checks if the key clears sorting
func (kb KeyBindings) IsClearSort(key string) bool {
	return kb.matchesKey(key, kb.ClearSort)
}

// GetSortColumn returns the column index for sorting, or -1 if not a sort key
func (kb KeyBindings) GetSortColumn(key string) int {
	for colIndex, keys := range kb.Sort {
		if kb.matchesKey(key, keys) {
			return colIndex
		}
	}
	return -1
}
