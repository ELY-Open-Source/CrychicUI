package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RailItem names an icon slot and its label.
type RailItem struct{ Icon, Label string }

// Rail is an icon column; the indicator slides.
func Rail(c *ui.Context, items []RailItem, active *int) ui.Element {
	if len(items) == 0 || active == nil || *active < 0 || *active >= len(items) {
		panic("layout: Rail needs items and an active index in range")
	}
	t := c.Theme()
	slot, gap, pad := t.Space(10), t.Space(1), t.Space(2)
	root := ui.Column(c).Width(slot + 2*pad).Gap(gap).Padding(pad).AlignItems(ui.Center)
	return root.Children(func() {
		mark := ui.Box(c).Absolute().Left(0).Width(3).Height(slot).Radius(999).Background(t.Accent)
		mark.Top(mark.AnimateTo("y", pad+float32(*active)*(slot+gap), ui.Snappy))
		for i, it := range items {
			b := ui.ButtonBase(c).Key(i).Size(slot, slot).Radius(t.Radius).Label(it.Label)
			bg, fg := t.Accent.Alpha(0), t.TextMuted
			if i == *active {
				bg, fg = t.Accent.Alpha(0.12), t.Accent
			} else if b.Hovered() {
				bg = t.SurfaceHover
			}
			b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(b.AnimateColor("fg", fg, ui.Snappy))
			kit.Press(b).Children(func() { ui.Icon(c, icon.Must(it.Icon)).Size(16, 16) })
			if kit.Clicked(b, sound.Tap) && i != *active {
				*active = i
				root.SetChanged()
				kit.Log("Rail", "selected", "index", i)
			}
		}
	})
}
