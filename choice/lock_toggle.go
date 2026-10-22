package choice

import "github.com/egoist/mygo/ui"

// LockToggle is a lock whose shackle lifts as it opens.
func LockToggle(c *ui.Context, label string, unlocked *bool) ui.Element {
	return iconToggle(c, "LockToggle", label, "lock", "lock-open", unlocked, 2)
}
