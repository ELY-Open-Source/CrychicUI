package choice

import "github.com/egoist/mygo/ui"

// MuteToggle is a speaker that mutes and unmutes.
func MuteToggle(c *ui.Context, label string, muted *bool) ui.Element {
	return iconToggle(c, "MuteToggle", label, "volume-2", "volume-x", muted, 0)
}
