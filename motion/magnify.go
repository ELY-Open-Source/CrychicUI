package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MagnifyItem is one icon of a Magnify row.
type MagnifyItem struct {
	Icon, Label string
}

// Magnify scales icons by pointer distance; a click picks one.
func Magnify(c *ui.Context, items []MagnifyItem, picked *int) ui.Element {
	if len(items) == 0 || picked == nil {
		panic("motion: Magnify needs items and a picked pointer")
	}
	for _, it := range items {
		if it.Label == "" {
			panic("motion: Magnify items need labels")
		}
		icon.Must(it.Icon)
	}
	t := c.Theme()
	size, gap, pad := t.Space(10), t.Space(2), t.Space(2)
	n := float32(len(items))
	width := 2 + 2*pad + n*size + (n-1)*gap
	root := ui.Row(c).AlignSelf(ui.Start).Gap(gap).AlignItems(ui.End).Padding(pad).Radius(t.Radius).Background(t.Surface).Border(1, t.Border)
	px, _, over := root.PointerPosition()
	rtl := c.Direction() == ui.RightToLeft
	return root.Children(func() {
		for i, it := range items {
			c.Key(i)
			center := 1 + pad + float32(i)*(size+gap) + size/2
			if rtl {
				center = width - center
			}
			s := float32(1)
			if over {
				s += 0.6 * max(0, 1-abs(px-center)/(2.5*size))
			}
			b := ui.ButtonBase(c).Label(it.Label).Size(size, size).Radius(t.Radius).Center().TransformOrigin(0.5, 1)
			k := b.AnimateTo("scale", s, ui.Snappy)
			b.Scale(k, k).Children(func() { ui.Icon(c, icon.Must(it.Icon)).Size(16, 16) })
			if kit.Clicked(b, sound.Tap) {
				*picked = i
				root.SetChanged()
				kit.Log("Magnify", "picked", "index", i)
			}
		}
	})
}
