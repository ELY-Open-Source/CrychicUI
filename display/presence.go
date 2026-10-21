package display

import "github.com/egoist/mygo/ui"

// Presence is a status dot and caller's label.
func Presence(c *ui.Context, status PresenceStatus, label string) ui.Element {
	if status < PresenceOffline || status > PresenceBusy || label == "" {
		panic("display: Presence needs a status and label")
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
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
		dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999)
		dot.Background(dot.AnimateColor("color", color, ui.Snappy))
		ui.Text(c, label).ContentTransition(ui.ContentTransition{})
	})
}
