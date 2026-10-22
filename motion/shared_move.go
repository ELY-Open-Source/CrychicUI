package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SharedMoveOptions tune a SharedMove.
type SharedMoveOptions struct {
	// Back labels the button closing the detail; "Back" when empty.
	Back string
}

// shared is one SharedMove item's flight identity.
type shared struct {
	root uint64
	i    int
}

// SharedMove flies a list item into its detail and back.
func SharedMove(c *ui.Context, items []string, open *int, detail func(i int), o SharedMoveOptions) ui.Element {
	if len(items) == 0 || open == nil || detail == nil || *open < -1 || *open >= len(items) {
		panic("motion: SharedMove needs items, an open index from -1 and a detail")
	}
	for _, s := range items {
		if s == "" {
			panic("motion: SharedMove items need text")
		}
	}
	if o.Back == "" {
		o.Back = "Back"
	}
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2))
	id := root.ID()
	fly := func(i int) ui.ElementTransition {
		return ui.ElementTransition{Animation: &ui.Snappy, Match: shared{id, i}}
	}
	return root.Children(func() {
		if *open < 0 {
			for i, s := range items {
				c.Key(i)
				row := ui.ButtonBase(c).Label(s).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).Transition(fly(i))
				bg := t.Surface
				if row.Hovered() {
					bg = t.SurfaceHover
				}
				kit.Press(row.Background(bg)).Children(func() { ui.Text(c, s) })
				if kit.Clicked(row, sound.Open) {
					*open = i
					root.SetChanged()
					kit.Log("SharedMove", "opened", "index", i)
				}
			}
			return
		}
		i := *open
		card := kit.Card(c.Key("detail")).Padding(t.Space(4)).Gap(t.Space(3)).AlignItems(ui.Start).Transition(fly(i))
		card.Children(func() {
			ui.Text(c, items[i]).Bold().FontSize(t.Rem(1.25))
			detail(i)
			if kit.Clicked(ui.Button(c, o.Back), sound.Close) {
				*open = -1
				root.SetChanged()
				kit.Log("SharedMove", "closed", "index", i)
			}
		})
	})
}
