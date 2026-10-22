package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NavItem is a destination's label and Lucide icon.
type NavItem struct{ Label, Icon string }

// items checks a component's items and active index.
func items(component, label string, list []NavItem, active *int) {
	if label == "" || len(list) == 0 || active == nil || !valid(*active, len(list)) {
		panic(fmt.Sprintf("navigation: %s needs a label, items and an active index", component))
	}
	for i, it := range list {
		if it.Label == "" {
			panic(fmt.Sprintf("navigation: %s item %d needs a label", component, i))
		}
		icon.Must(it.Icon)
	}
}

// navList lays out items under a sliding highlight.
func navList(c *ui.Context, component string, id uint64, list []NavItem, active *int, narrow float32) (picked bool) {
	t := c.Theme()
	shown := *active
	for i, it := range list {
		c.Key(i)
		b := kit.Press(ui.ButtonBase(c).Label(it.Label).FillWidth().Radius(t.Radius).Checked(i == shown))
		if narrow > 0.5 {
			b.Tooltip(it.Label)
		}
		bg, fg := t.SurfaceHover.Alpha(0), t.TextMuted
		switch {
		case i == shown:
			fg = t.Accent
		case b.Hovered():
			bg, fg = t.SurfaceHover, t.Text
		}
		b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(b.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
			stack(c, i == shown, func() ui.Element {
				return ui.Box(c).Radius(t.Radius).Background(t.Accent.Alpha(0.12)).Transition(glide(id))
			}, func() ui.Element {
				return ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), t.Space(3)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must(it.Icon)).Size(16, 16)
					ui.Text(c, it.Label).SingleLine().Opacity(1 - narrow)
				})
			})
		})
		if kit.Clicked(b, sound.Select) && i != *active {
			*active, picked = i, true
			kit.Log(component, "picked", "index", i)
		}
	}
	return picked
}
