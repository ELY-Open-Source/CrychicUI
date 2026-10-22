package choice

import "github.com/egoist/mygo/ui"

// VisibilityToggle is an eye that opens and closes.
func VisibilityToggle(c *ui.Context, label string, visible *bool) ui.Element {
	return iconToggle(c, "VisibilityToggle", label, "eye-off", "eye", visible, 0)
}
