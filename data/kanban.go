package data

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Lane is a titled column of cards.
type Lane struct {
	Title string
	Cards []string
}

type kanbanDrag struct {
	id   uint64
	card string
}

// Kanban drags cards between lanes; Alt+arrows move them.
func Kanban(c *ui.Context, label string, lanes *[]Lane) ui.Element {
	if label == "" || lanes == nil || len(*lanes) == 0 {
		panic(fmt.Sprintf("data: Kanban %q needs lanes", label))
	}
	var titles, cards []string
	for _, l := range *lanes {
		titles = append(titles, l.Title)
		cards = append(cards, l.Cards...)
	}
	unique("Kanban", label, titles)
	if len(cards) > 0 {
		unique("Kanban", label, cards)
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Role(ui.RoleGroup).Label(label)
	id := root.ID()
	moved := ui.Local(root, "moved", func() string { return "" })
	move := func(card string, lane, at int) {
		from, i := 0, 0
		for l, ln := range *lanes {
			if k := slices.Index(ln.Cards, card); k >= 0 {
				from, i = l, k
			}
		}
		next := slices.Clone(*lanes)
		for l := range next {
			next[l].Cards = slices.Clone(next[l].Cards)
		}
		next[from].Cards = slices.Delete(next[from].Cards, i, i+1)
		if from == lane && at > i {
			at--
		}
		at = min(max(at, 0), len(next[lane].Cards))
		if from == lane && at == i {
			return
		}
		next[lane].Cards = slices.Insert(next[lane].Cards, at, card)
		*lanes = next
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("Kanban", "moved", "lane", lane, "to", at)
	}
	side := 1
	if c.Direction() == ui.RightToLeft {
		side = -1
	}
	slide := func(card string) ui.ElementTransition {
		return ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{Scale: 0.96}, Position: true, Match: kanbanDrag{id, card}}
	}
	var cols []func()
	for li, ln := range *lanes {
		cols = append(cols, func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Text(c, ln.Title).Bold().Grow(1)
				kit.Numeric(c, fmt.Sprint(len(ln.Cards))).TextColor(t.TextMuted)
			})
			zone := ui.Column(c).FillWidth().MinHeight(t.Space(16)).Padding(0, 0, t.Space(8), 0).Gap(t.Space(2)).Radius(t.Radius).Role(ui.RoleList).Label(ln.Title)
			_, over := ui.DragOver[kanbanDrag](zone)
			zone.Background(zone.AnimateColor("bg", t.Accent.Alpha(0.08*unit(over)), ui.Snappy))
			took := false
			zone.Children(func() {
				for ci, card := range ln.Cards {
					c.Key(card)
					b := ui.Box(c).FillWidth().Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius).Background(t.Background).
						Border(1, t.Border).Shadow(0, 1, 3, 0, kit.Shade(0.06)).Focusable().Label(card).
						Cursor(ui.CursorGrab).Drag(kanbanDrag{id, card}).Transition(slide(card))
					if *moved == card {
						b.AutoFocus()
						*moved = ""
					}
					if b.Dragging() {
						b.Opacity(0.4)
					}
					b.Children(func() { ui.Text(c, card) })
					if d, ok := ui.Drop[kanbanDrag](b); ok && d.id == id {
						took = true
						move(d.card, li, ci)
					}
					step := func(lane, at int) {
						if lane >= 0 && lane < len(*lanes) {
							*moved = card
							move(card, lane, at)
						}
					}
					switch {
					case b.Shortcut(ui.Alt, ui.KeyUp):
						step(li, ci-1)
					case b.Shortcut(ui.Alt, ui.KeyDown):
						step(li, ci+2)
					case b.Shortcut(ui.Alt, ui.KeyLeft):
						step(li-side, ci)
					case b.Shortcut(ui.Alt, ui.KeyRight):
						step(li+side, ci)
					}
				}
			})
			if d, ok := ui.Drop[kanbanDrag](zone); ok && d.id == id && !took {
				move(d.card, li, len(ln.Cards))
			}
		})
	}
	return root.Children(func() { layout.Board(c, cols) })
}
