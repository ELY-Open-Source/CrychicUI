package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Quarter is one roadmap quarter and initiatives.
type Quarter struct {
	Label string
	Items []string
}

// RoadmapOptions mark the Current quarter with Now.
type RoadmapOptions struct{ Current, Now string }

// Roadmap shows quarters with initiatives.
func Roadmap(c *ui.Context, quarters []Quarter, o RoadmapOptions) ui.Element {
	if len(quarters) == 0 {
		panic("project: Roadmap needs quarters")
	}
	labels := make([]string, len(quarters))
	for i, quarter := range quarters {
		if len(quarter.Items) == 0 {
			panic(fmt.Sprintf("project: Roadmap quarter %q needs initiatives", quarter.Label))
		}
		unique("Roadmap", quarter.Items)
		labels[i] = quarter.Label
	}
	unique("Roadmap", labels)
	def(&o.Now, "Now")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		n := 0
		for _, quarter := range quarters {
			now := quarter.Label == o.Current
			c.Key(quarter.Label)
			col := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1.5)).Padding(t.Space(2)).Radius(t.Radius)
			col.Background(col.AnimateColor("bg", t.SurfaceHover.Alpha(0.5).Mix(t.Accent.Alpha(0.08), unit(now)), ui.Snappy))
			col.Border(1, col.AnimateColor("line", t.Border.Alpha(0).Mix(t.Accent.Alpha(0.4), unit(now)), ui.Snappy))
			col.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, quarter.Label).FontSize(t.Rem(0.85)).Bold().TextColor(col.AnimateColor("ink", t.Text.Mix(t.Accent, unit(now)), ui.Snappy)).
						Grow(1).MinWidth(0).SingleLine()
					if now {
						kit.Chip(c, o.Now, t.Accent)
					} else {
						kit.Chip(c, fmt.Sprint(len(quarter.Items)), t.TextMuted)
					}
				})
				items := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				items.Children(func() {
					for _, item := range quarter.Items {
						ink := labelInk(t, item)
						c.Key(item)
						tile := kit.Lift(kit.Card(c).FillWidth().Row().Gap(t.Space(2)).AlignItems(ui.Stretch).Padding(t.Space(2)))
						kit.Stagger(c, root, tile, n).Children(func() {
							ui.Box(c).Width(t.Space(0.75)).Shrink(0).Radius(999).Background(ink)
							ui.Text(c, item).FontSize(t.Rem(0.75)).Bold().Grow(1).MinWidth(0).Wrap()
						})
						n++
					}
				})
			})
		}
	})
	return root
}
