package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Navbar is a brand, links and actions; the underline slides.
func Navbar(c *ui.Context, brand string, links []string, active *int, actions func()) ui.Element {
	if brand == "" || len(links) == 0 || active == nil || !valid(*active, len(links)) {
		panic("navigation: Navbar needs a brand, links and an active index")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(6)).Padding(t.Space(1), t.Space(4)).
		Background(t.Background).Border(1, t.Border).Radius(t.Radius)
	id, shown := root.ID(), *active
	return root.Children(func() {
		ui.Text(c, brand).Bold().FontSize(t.Rem(1.1)).SingleLine()
		ui.Row(c).Grow(1).Gap(t.Space(1)).FocusGroup(ui.Horizontal).Label(brand).Children(func() {
			for i, l := range links {
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(l).Checked(i == shown).Padding(t.Space(2), t.Space(2), 0, t.Space(2)))
				fg := t.TextMuted
				if i == shown || b.Hovered() {
					fg = t.Text
				}
				b.TextColor(b.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
					ui.Column(c).Gap(t.Space(1.5)).AlignItems(ui.Stretch).Children(func() {
						ui.Text(c, l).SingleLine()
						line := ui.Box(c).Height(2).Radius(999)
						if i == shown {
							line.Background(t.Accent).Transition(glide(id))
						}
					})
				})
				if kit.Clicked(b, sound.Tap) && i != *active {
					*active = i
					root.SetChanged()
					kit.Log("Navbar", "picked", "index", i)
				}
			}
		})
		if actions != nil {
			ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(actions)
		}
	})
}
