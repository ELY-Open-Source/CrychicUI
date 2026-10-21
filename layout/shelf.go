package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ShelfOptions tune a Shelf.
type ShelfOptions struct {
	// Action labels the link; "See all" when empty.
	Action string
}

// Shelf is a titled sideways row with a see-all link.
func Shelf(c *ui.Context, title string, o ShelfOptions, content func()) (root, seeAll ui.Element) {
	if o.Action == "" {
		o.Action = "See all"
	}
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			ui.Text(c, title).Bold().Grow(1)
			seeAll = ui.ButtonBase(c).Gap(t.Space(1)).Radius(t.Radius).Padding(t.Space(1))
			fg := t.Accent
			if seeAll.Hovered() {
				fg = t.AccentHover
			}
			kit.Press(seeAll).TextColor(seeAll.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
				ui.Text(c, o.Action)
				arrow := ui.Icon(c, icon.Must("chevron-right")).Size(16, 16)
				arrow.Translate(arrow.AnimateTo("nudge", 2*unit(seeAll.Hovered()), ui.Snappy), 0)
			})
			if kit.Clicked(seeAll, sound.Tap) {
				kit.Log("Shelf", "see all")
			}
		})
		ui.ScrollHorizontal(c).FillWidth().Children(func() {
			ui.Row(c).Gap(t.Space(3)).Padding(t.Space(1), 0).Children(content)
		})
	})
	return root, seeAll
}
