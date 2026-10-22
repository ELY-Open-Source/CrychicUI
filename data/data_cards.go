package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Entry is a titled record of key and value pairs.
type Entry struct {
	Title  string
	Fields []Pair
}

// DataCards lays records out as cards in a grid.
func DataCards(c *ui.Context, label string, entries []Entry, chosen *int) ui.Element {
	if label == "" || len(entries) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(entries) {
		panic(fmt.Sprintf("data: DataCards %q needs entries and a choice or -1", label))
	}
	var titles []string
	for _, e := range entries {
		titles = append(titles, e.Title)
	}
	unique("DataCards", label, titles)
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).Role(ui.RoleList).Label(label)
	shown := *chosen
	return root.Children(func() {
		for i, e := range entries {
			c.Key(e.Title)
			on := i == shown
			b := kit.Press(ui.ButtonBase(c).Label(e.Title).Width(t.Space(56)).Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).
				Background(t.Background).Shadow(0, 1, 3, 0, kit.Shade(0.06)).Checked(on)).Column().AlignItems(ui.Stretch)
			b.Translate(0, -t.Space(1)*b.AnimateTo("lift", unit(b.Hovered()), ui.Snappy)).
				Border(1, b.AnimateColor("edge", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
			kit.Arrive(b, ui.Pose{Y: t.Space(3), Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				ui.Text(c, e.Title).Bold()
				for k, f := range e.Fields {
					c.Key(k)
					ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
						kit.Muted(c, f.Key).Grow(1).SingleLine()
						ui.Text(c, f.Value).SingleLine()
					})
				}
			})
			if b.Clicked() {
				*chosen = i
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("DataCards", "chose", "index", i)
			}
		}
	})
}
