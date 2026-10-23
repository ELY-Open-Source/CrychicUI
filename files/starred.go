package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StarredOptions are the star labels.
type StarredOptions struct{ Star, Unstar string }

// Starred toggles starred files; cue Toggle.
func Starred(c *ui.Context, label string, list []FileItem, starred *[]bool, o StarredOptions) ui.Element {
	validFiles("Starred", label, list)
	if starred == nil || len(*starred) != len(list) {
		panic("files: Starred needs one flag per file")
	}
	def(&o.Star, "Star")
	def(&o.Unstar, "Unstar")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, item := range list {
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon.Must(itemGlyph(item))).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, item.Name).SingleLine().Grow(1)
				name := o.Star + " " + item.Name
				ink := t.TextMuted
				if (*starred)[i] {
					name = o.Unstar + " " + item.Name
					ink = t.Warning
				}
				b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(7), t.Space(7)).Center().Radius(999).
					TextColor(ink))
				b.Children(func() { ui.Icon(c, icon.Must("star")).Size(16, 16).ContentTransition(ui.ContentTransition{}) })
				if b.Clicked() {
					(*starred)[i] = !(*starred)[i]
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Starred", "toggled", "index", i, "starred", (*starred)[i])
				}
			})
		}
	})
	return root
}
