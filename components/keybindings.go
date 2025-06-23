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
	Quit         []string
	Help         []string
	Refresh      []string
	PageSizeUp   []string
	PageSizeDown []string
	ResetPage    []string
	ClearSort    []string
	Sort1        []string
	Sort2        []string
	Sort3        []string
	Sort4        []string
	Sort5        []string
	Sort6        []string
	Sort7        []string
	Sort8        []string
	Sort9        []string
}

// DefaultKeyBindings returns the default key bindings
func DefaultKeyBindings() *KeyBindings {
	return &KeyBindings{
		Quit:         []string{"q", "esc", "ctrl+c"},
		Up:           []string{"up", "k"},
		Down:         []string{"down", "j"},
		Left:         []string{"left", "h"},
		Right:        []string{"right", "l"},
		PageUp:       []string{"pageup", "pgup", "ctrl+b"},
		PageDown:     []string{"pagedown", "pgdown", "ctrl+f"},
		Home:         []string{"home", "g"},
		End:          []string{"end", "G"},
		Search:       []string{"/"},
		Help:         []string{"?"},
		Refresh:      []string{"r"},
		PageSizeUp:   []string{"+", "="},
		PageSizeDown: []string{"-", "_"},
		ResetPage:    []string{"0"},
		ClearSort:    []string{"c"},
		Sort1:        []string{"1"},
		Sort2:        []string{"2"},
		Sort3:        []string{"3"},
		Sort4:        []string{"4"},
		Sort5:        []string{"5"},
		Sort6:        []string{"6"},
		Sort7:        []string{"7"},
		Sort8:        []string{"8"},
		Sort9:        []string{"9"},
	}
}

// VimKeyBindings returns Vim-style key bindings
func VimKeyBindings() *KeyBindings {
	return &KeyBindings{
		Quit:         []string{"q", "esc"},
		Up:           []string{"k", "up"},
		Down:         []string{"j", "down"},
		Left:         []string{"h", "left"},
		Right:        []string{"l", "right"},
		PageUp:       []string{"ctrl+u", "pageup", "pgup"},
		PageDown:     []string{"ctrl+d", "pagedown", "pgdown"},
		Home:         []string{"gg", "home"},
		End:          []string{"G", "end"},
		Search:       []string{"/"},
		Help:         []string{"?"},
		Refresh:      []string{"r"},
		PageSizeUp:   []string{"+"},
		PageSizeDown: []string{"-"},
		ResetPage:    []string{"0"},
		ClearSort:    []string{"c"},
		Sort1:        []string{"1"},
		Sort2:        []string{"2"},
		Sort3:        []string{"3"},
		Sort4:        []string{"4"},
		Sort5:        []string{"5"},
		Sort6:        []string{"6"},
		Sort7:        []string{"7"},
		Sort8:        []string{"8"},
		Sort9:        []string{"9"},
	}
}

// EmacsKeyBindings returns Emacs-style key bindings
func EmacsKeyBindings() *KeyBindings {
	return &KeyBindings{
		Quit:         []string{"ctrl+x ctrl+c", "ctrl+g"},
		Up:           []string{"ctrl+p", "up"},
		Down:         []string{"ctrl+n", "down"},
		Left:         []string{"ctrl+b", "left"},
		Right:        []string{"ctrl+f", "right"},
		PageUp:       []string{"alt+v", "pageup", "pgup"},
		PageDown:     []string{"ctrl+v", "pagedown", "pgdown"},
		Home:         []string{"ctrl+a", "alt+<", "home"},
		End:          []string{"ctrl+e", "alt+>", "end"},
		Search:       []string{"ctrl+s"},
		Help:         []string{"ctrl+h"},
		Refresh:      []string{"ctrl+l"},
		PageSizeUp:   []string{"ctrl++"},
		PageSizeDown: []string{"ctrl+-"},
		ResetPage:    []string{"ctrl+0"},
		ClearSort:    []string{"ctrl+c"},
		Sort1:        []string{"ctrl+1"},
		Sort2:        []string{"ctrl+2"},
		Sort3:        []string{"ctrl+3"},
		Sort4:        []string{"ctrl+4"},
		Sort5:        []string{"ctrl+5"},
		Sort6:        []string{"ctrl+6"},
		Sort7:        []string{"ctrl+7"},
		Sort8:        []string{"ctrl+8"},
		Sort9:        []string{"ctrl+9"},
	}
}

// matchesKey checks if a key string matches any of the configured keys
func (kb *KeyBindings) matchesKey(key string, keys []string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

// IsUp checks if the key is an up navigation key
func (kb *KeyBindings) IsUp(key string) bool {
	return kb.matchesKey(key, kb.Up)
}

// IsDown checks if the key is a down navigation key
func (kb *KeyBindings) IsDown(key string) bool {
	return kb.matchesKey(key, kb.Down)
}

// IsLeft checks if the key is a left navigation key
func (kb *KeyBindings) IsLeft(key string) bool {
	return kb.matchesKey(key, kb.Left)
}

// IsRight checks if the key is a right navigation key
func (kb *KeyBindings) IsRight(key string) bool {
	return kb.matchesKey(key, kb.Right)
}

// IsPageUp checks if the key is a page up key
func (kb *KeyBindings) IsPageUp(key string) bool {
	return kb.matchesKey(key, kb.PageUp)
}

// IsPageDown checks if the key is a page down key
func (kb *KeyBindings) IsPageDown(key string) bool {
	return kb.matchesKey(key, kb.PageDown)
}

// IsHome checks if the key is a home key
func (kb *KeyBindings) IsHome(key string) bool {
	return kb.matchesKey(key, kb.Home)
}

// IsEnd checks if the key is an end key
func (kb *KeyBindings) IsEnd(key string) bool {
	return kb.matchesKey(key, kb.End)
}

// IsSearch checks if the key is a search key
func (kb *KeyBindings) IsSearch(key string) bool {
	return kb.matchesKey(key, kb.Search)
}

// IsQuit checks if the key is a quit key
func (kb *KeyBindings) IsQuit(key string) bool {
	return kb.matchesKey(key, kb.Quit)
}

// IsHelp checks if the key is a help key
func (kb *KeyBindings) IsHelp(key string) bool {
	return kb.matchesKey(key, kb.Help)
}

// IsRefresh checks if the key is a refresh key
func (kb *KeyBindings) IsRefresh(key string) bool {
	return kb.matchesKey(key, kb.Refresh)
}

// IsPageSizeUp checks if the key increases page size
func (kb *KeyBindings) IsPageSizeUp(key string) bool {
	return kb.matchesKey(key, kb.PageSizeUp)
}

// IsPageSizeDown checks if the key decreases page size
func (kb *KeyBindings) IsPageSizeDown(key string) bool {
	return kb.matchesKey(key, kb.PageSizeDown)
}

// IsResetPage checks if the key resets to optimal page size
func (kb *KeyBindings) IsResetPage(key string) bool {
	return kb.matchesKey(key, kb.ResetPage)
}

// IsClearSort checks if the key clears sorting
func (kb *KeyBindings) IsClearSort(key string) bool {
	return kb.matchesKey(key, kb.ClearSort)
}

// GetSortColumn returns the column index for sorting, or -1 if not a sort key
func (kb *KeyBindings) GetSortColumn(key string) int {
	sortKeys := map[int][]string{
		0: kb.Sort1,
		1: kb.Sort2,
		2: kb.Sort3,
		3: kb.Sort4,
		4: kb.Sort5,
		5: kb.Sort6,
		6: kb.Sort7,
		7: kb.Sort8,
		8: kb.Sort9,
	}

	for colIndex, keys := range sortKeys {
		if kb.matchesKey(key, keys) {
			return colIndex
		}
	}
	return -1
}
