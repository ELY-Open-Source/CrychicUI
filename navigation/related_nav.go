package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RelatedNav is a titled list of links after content.
func RelatedNav(c *ui.Context, title string, links []string, picked *int) ui.Element {
	if title == "" || len(links) == 0 || picked == nil || *picked < -1 || *picked >= len(links) {
		panic(fmt.Sprintf("navigation: RelatedNav %q needs links and a pick or -1", title))
	}
	t := c.Theme()
	nudge := t.Space(1)
	if c.Direction() == ui.RightToLeft {
		nudge = -nudge
	}
	root := ui.Column(c).Width(t.Space(80)).Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Label(title)
	return root.Children(func() {
		kit.Muted(c, title).Bold().FontSize(t.Rem(0.8)).Padding(t.Space(1), t.Space(3))
		for i, l := range links {
			c.Key(i)
			b := rowBase(c, l)
			k := b.AnimateTo("nudge", unit(b.Hovered()), ui.Snappy)
			b.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, l).Grow(1).SingleLine().TextColor(t.Accent)
					ui.Icon(c, icon.Must("arrow-right")).Size(16, 16).TextColor(t.Accent).Translate(k*nudge, 0)
				})
			})
			if kit.Clicked(b, sound.Tap) {
				*picked = i
				root.SetChanged()
				kit.Log("RelatedNav", "picked", "index", i)
			}
		}
	})
}
