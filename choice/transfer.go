package choice

import (
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TransferOptions tune a Transfer; empty titles read Available, Chosen.
type TransferOptions struct {
	// LeftTitle and RightTitle head the lists.
	LeftTitle, RightTitle string
}

type transferMatch struct {
	id   uint64
	item string
}

// Transfer is two lists whose items fly across when clicked.
func Transfer(c *ui.Context, left, right *[]string, o TransferOptions) ui.Element {
	if left == nil || right == nil {
		panic("choice: Transfer needs both lists")
	}
	distinct("Transfer", append(slices.Clone(*left), *right...))
	if o.LeftTitle == "" {
		o.LeftTitle = "Available"
	}
	if o.RightTitle == "" {
		o.RightTitle = "Chosen"
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch)
	id := root.ID()
	list := func(title string, from, to *[]string, arrow string) {
		kit.Card(c).Grow(1).Basis(0).Padding(t.Space(1)).Gap(t.Space(0.5)).Children(func() {
			ui.Row(c).Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2.5)).Children(func() {
				ui.Text(c, title).Bold().Grow(1)
				kit.Numeric(c, strconv.Itoa(len(*from))).TextColor(t.TextMuted)
			})
			moved := ""
			for _, item := range *from {
				c.Key(item)
				row := ui.ButtonBase(c).Label(item).Justify(ui.Start).Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2.5)).
					Radius(t.Radius).Transition(ui.ElementTransition{Animation: &ui.Snappy, Match: transferMatch{id, item}})
				bg := t.Background
				if row.Hovered() {
					bg = t.SurfaceHover
				}
				row.Background(row.AnimateColor("bg", bg, quick)).Children(func() {
					ui.Text(c, item).SingleLine().Grow(1)
					ui.Icon(c, icon.Must(arrow)).Size(16, 16).TextColor(t.TextMuted).Opacity(row.AnimateTo("hint", unit(row.Hovered()), quick))
				})
				if kit.Clicked(row, sound.Drop) {
					moved = item
				}
			}
			if moved != "" {
				*from = slices.DeleteFunc(slices.Clone(*from), func(s string) bool { return s == moved })
				*to = append(slices.Clone(*to), moved)
				root.SetChanged()
				kit.Log("Transfer", "moved", "left", len(*left), "right", len(*right))
			}
		})
	}
	return root.Children(func() {
		list(o.LeftTitle, left, right, "arrow-right")
		list(o.RightTitle, right, left, "arrow-left")
	})
}
