package display

import "github.com/egoist/mygo/ui"

// PresenceStatus is a person's availability.
type PresenceStatus int

// Presence states.
const (
	PresenceOffline PresenceStatus = iota
	PresenceOnline
	PresenceAway
	PresenceBusy
)

// AvatarStatus adds an animated presence dot to an avatar.
func AvatarStatus(c *ui.Context, name string, image *ui.Bitmap, status PresenceStatus) ui.Element {
	if status < PresenceOffline || status > PresenceBusy {
		panic("display: AvatarStatus status out of range")
	}
	t := c.Theme()
	color := t.TextMuted
	switch status {
	case PresenceOnline:
		color = t.Success
	case PresenceAway:
		color = t.Warning
	case PresenceBusy:
		color = t.Danger
	}
	root := ui.Box(c).AlignSelf(ui.Start)
	return root.Children(func() {
		avatarCircle(c, name, image, t.Space(10))
		dot := ui.Box(c).Absolute().Bottom(0).Right(0).Size(t.Space(2.5), t.Space(2.5)).Radius(999).Border(t.Space(0.5), t.Background)
		dot.Background(dot.AnimateColor("presence", color, ui.Snappy))
	})
}
