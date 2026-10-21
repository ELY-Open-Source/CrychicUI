package choice

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

var quick = ui.Animation{Duration: 120 * time.Millisecond}

// trigger styles a button that opens a popup below it.
func trigger(c *ui.Context, b ui.Element, open bool, content func()) ui.Element {
	t := c.Theme()
	face := t.Surface
	switch {
	case b.Pressed():
		face = t.SurfacePressed
	case b.Hovered():
		face = t.SurfaceHover
	}
	b.Justify(ui.SpaceBetween).MinWidth(t.Space(35)).Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
		Background(b.AnimateColor("face", face, quick)).Border(1, t.Border)
	return b.Children(func() {
		content()
		chevron := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16).TextColor(t.TextMuted)
		chevron.Rotate(chevron.AnimateTo("turn", 180*unit(open), ui.Snappy))
	})
}

// panel styles a popup panel that scales in.
func panel(c *ui.Context, p ui.Element) ui.Element {
	t := c.Theme()
	p.Margin(t.Space(1), 0, 0, 0).Padding(t.Space(1)).Radius(t.Radius+2).Background(t.Background).Border(1, t.Border).
		Shadow(0, 6, 20, 0, kit.Shade(0.18))
	return kit.Arrive(p, ui.Pose{Y: -4, Scale: 0.96, Fade: 1}, 0)
}

// entry styles a popup row whose check pops in.
func entry(c *ui.Context, item ui.Element, label string, chosen, hot bool) ui.Element {
	t := c.Theme()
	bg, fg := ui.Color{}, t.Text
	if hot {
		bg, fg = t.Accent, t.AccentText
	}
	item.Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius).
		Background(item.AnimateColor("bg", bg, quick)).TextColor(fg)
	return item.Children(func() {
		ui.Text(c, label).SingleLine().Grow(1)
		check := ui.Icon(c, icon.Must("check")).Size(16, 16)
		k := check.AnimateTo("on", unit(chosen), ui.Bouncy)
		check.Scale(k, k).Opacity(min(max(k, 0), 1))
	})
}
