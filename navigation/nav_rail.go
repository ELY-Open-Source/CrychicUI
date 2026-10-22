package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NavRail is a narrow rail: icons with labels under them.
func NavRail(c *ui.Context, label string, list []NavItem, active *int) ui.Element {
	items("NavRail", label, list, active)
	t := c.Theme()
	root := ui.Column(c).Width(t.Space(20)).Gap(t.Space(3)).Padding(t.Space(3), 0).AlignItems(ui.Center).
		FocusGroup(ui.Vertical).Label(label)
	id, shown := root.ID(), *active
	root.Children(func() {
		for i, it := range list {
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(it.Label).Checked(i == shown))
			fg := t.TextMuted
			if i == shown || b.Hovered() {
				fg = t.Text
			}
			b.TextColor(b.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
				ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					stack(c, i == shown, func() ui.Element {
						return ui.Box(c).Radius(999).Background(t.Accent.Alpha(0.16)).Transition(glide(id))
					}, func() ui.Element {
						return ui.Box(c).Size(t.Space(14), t.Space(8)).Center().Children(func() {
							ui.Icon(c, icon.Must(it.Icon)).Size(16, 16)
						})
					})
					ui.Text(c, it.Label).FontSize(t.Rem(0.75)).SingleLine()
				})
			})
			if kit.Clicked(b, sound.Select) && i != *active {
				*active = i
				root.SetChanged()
				kit.Log("NavRail", "picked", "index", i)
			}
		}
	})
	return root
}
