package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BottomTabs are mobile tabs; the active icon bounces.
func BottomTabs(c *ui.Context, label string, list []NavItem, active *int) ui.Element {
	items("BottomTabs", label, list, active)
	t := c.Theme()
	root := ui.Row(c).FillWidth().Padding(t.Space(1.5), t.Space(2)).Background(t.Background).Border(1, t.Border).
		FocusGroup(ui.Horizontal).Label(label)
	shown := *active
	root.Children(func() {
		for i, it := range list {
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(it.Label).Grow(1).Basis(0).Padding(t.Space(1.5), 0).Checked(i == shown))
			fg := t.TextMuted
			if i == shown {
				fg = t.Accent
			}
			b.TextColor(b.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
				ui.Column(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					ic := ui.Icon(c, icon.Must(it.Icon)).Size(16, 16)
					k := ic.AnimateTo("lift", unit(i == shown), ui.Bouncy)
					ic.Translate(0, -t.Space(0.5)*k).Scale(1+0.15*k, 1+0.15*k)
					ui.Text(c, it.Label).FontSize(t.Rem(0.75)).SingleLine()
				})
			})
			if kit.Clicked(b, sound.Select) && i != *active {
				*active = i
				root.SetChanged()
				kit.Log("BottomTabs", "picked", "index", i)
			}
		}
	})
	return root
}
