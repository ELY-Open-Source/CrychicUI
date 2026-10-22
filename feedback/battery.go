package feedback

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// Battery fills to a level; charging pulses a bolt.
func Battery(c *ui.Context, label string, level float64, charging bool) ui.Element {
	if label == "" {
		panic("feedback: Battery needs a label")
	}
	v := fraction("Battery", level)
	t := c.Theme()
	w, h, line := t.Space(8), t.Space(4), t.Space(0.375)
	tint := t.Success
	switch {
	case charging:
		tint = t.Accent
	case v <= 0.2:
		tint = t.Danger
	}
	root := ui.Row(c).AlignSelf(ui.Start).Role(ui.RoleMeter).Label(label).Range(0, 1, level)
	return root.Children(func() {
		body := ui.Box(c).Size(w, h).Radius(t.Space(1)).Border(line, t.TextMuted).Padding(t.Space(0.5)).Center()
		body.Children(func() {
			fill := ui.Box(c).Absolute().Top(t.Space(0.5)).Bottom(t.Space(0.5)).Start(t.Space(0.5)).Radius(t.Space(0.5))
			fill.Background(fill.AnimateColor("tint", tint, quick)).Width((w - 2*line - t.Space(1)) * fill.AnimateTo("v", v, ui.Smooth))
			if charging {
				bolt := ui.Icon(c, icon.Must("zap")).Size(16, 16).TextColor(t.Text)
				bolt.Opacity(0.4 + 0.6*bolt.Loop("pulse", 1200*time.Millisecond, ui.Bounce(ui.EaseInOut)))
			}
		})
		ui.Box(c).Size(t.Space(0.75), h/2).Radius(0, t.Space(0.5), t.Space(0.5), 0).Background(t.TextMuted)
	})
}
