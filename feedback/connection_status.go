package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Connection is online, reconnecting or offline.
type Connection uint8

// The three connection states.
const (
	Online Connection = iota
	Reconnecting
	Offline
)

// ConnectionStatusOptions tune a ConnectionStatus; empty copy reads English.
type ConnectionStatusOptions struct {
	Online, Reconnecting, Offline string
}

// ConnectionStatus is a colored dot and text; reconnecting pulses.
func ConnectionStatus(c *ui.Context, state Connection, o ConnectionStatusOptions) ui.Element {
	if o.Online == "" {
		o.Online = "Online"
	}
	if o.Reconnecting == "" {
		o.Reconnecting = "Reconnecting"
	}
	if o.Offline == "" {
		o.Offline = "Offline"
	}
	t := c.Theme()
	var text string
	var tint ui.Color
	switch state {
	case Online:
		text, tint = o.Online, t.Success
	case Reconnecting:
		text, tint = o.Reconnecting, t.Warning
	case Offline:
		text, tint = o.Offline, t.TextMuted
	default:
		panic("feedback: unknown Connection")
	}
	root := ui.Row(c).Gap(t.Space(2)).AlignSelf(ui.Start).Role(ui.RoleStatus).Label(text)
	return root.Children(func() {
		d := t.Space(2)
		dot := ui.Box(c).Size(d, d).Radius(999)
		dot.Background(dot.AnimateColor("tint", tint, quick))
		if state == Reconnecting {
			dot.Opacity(0.35 + 0.65*dot.Loop("pulse", time.Second, ui.Bounce(ui.EaseInOut)))
		}
		ui.Text(c, text).FontSize(t.Rem(0.85)).ContentTransition(ui.ContentTransition{})
	})
}
