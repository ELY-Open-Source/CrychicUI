package project

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TeamWeek is one person's weekly loads; 1 is full capacity.
type TeamWeek struct {
	Person string
	Loads  []float64
}

// TeamLoadOptions word the card; Weeks names each column.
type TeamLoadOptions struct {
	Title, Week string
	Weeks       []string
}

// TeamLoad shows weeks by person shaded.
func TeamLoad(c *ui.Context, weeks []TeamWeek, o TeamLoadOptions) ui.Element {
	def(&o.Title, "Team load")
	def(&o.Week, "W%d")
	if len(weeks) == 0 || len(weeks[0].Loads) == 0 || len(weeks[0].Loads) > 8 {
		panic(fmt.Sprintf("project: TeamLoad %q needs people and 1 to 8 weeks", o.Title))
	}
	n := len(weeks[0].Loads)
	if o.Weeks == nil {
		for i := range n {
			o.Weeks = append(o.Weeks, fmt.Sprintf(o.Week, i+1))
		}
	}
	if len(o.Weeks) != n {
		panic(fmt.Sprintf("project: TeamLoad %q has %d week names for %d weeks", o.Title, len(o.Weeks), n))
	}
	names := make([]string, len(weeks))
	for i, week := range weeks {
		if len(week.Loads) != n {
			panic(fmt.Sprintf("project: TeamLoad %q gives %q %d weeks, not %d", o.Title, week.Person, len(week.Loads), n))
		}
		for _, load := range week.Loads {
			if load < 0 || load > 2 || !finite(load) {
				panic(fmt.Sprintf("project: TeamLoad %q has load %v for %q", o.Title, load, week.Person))
			}
		}
		names[i] = week.Person
	}
	unique("TeamLoad", names)
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "users", o.Title, nil)
		cell := t.Space(12)
		view := ui.ScrollHorizontal(c).FillWidth()
		seen := view.LayoutBounds().W
		kit.Settle(c, view, []float32{seen}, func() []float32 { return []float32{view.LayoutBounds().W} })
		view.Children(func() {
			grid := ui.Column(c).Width(max(seen, t.Space(32)+float32(n)*(cell+t.Space(1)))).Gap(t.Space(1)).Padding(t.Space(3))
			grid.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
					ui.Box(c).Width(t.Space(26)).Shrink(0)
					for i, name := range o.Weeks {
						c.Key(i)
						ui.Row(c).Grow(1).Basis(0).MinWidth(cell).Justify(ui.Center).Children(func() {
							kit.Muted(c, name).FontSize(t.Rem(0.65)).Bold().MinWidth(0).SingleLine()
						})
					}
				})
				for i, week := range weeks {
					c.Key(week.Person)
					row := ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center)
					kit.Stagger(c, grid, row, i).Children(func() {
						ui.Row(c).Width(t.Space(26)).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							display.Avatar(c, week.Person, nil).Size(t.Space(6), t.Space(6)).FontSize(t.Rem(0.5))
							ui.Text(c, week.Person).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
						})
						for j, load := range week.Loads {
							c.Key(j)
							ink, alpha := t.Accent, 0.08+0.42*float32(load)
							if load > 1 {
								ink, alpha = t.Danger, 0.3+0.2*float32(load-1)
							}
							box := ui.Box(c).Grow(1).Basis(0).MinWidth(cell).Height(t.Space(8)).Center().Radius(t.Space(1))
							box.Background(box.AnimateColor("bg", ink.Alpha(alpha), ui.Smooth)).Children(func() {
								kit.Numeric(c, fmt.Sprintf("%.0f%%", math.Round(load*100))).FontSize(t.Rem(0.62)).Bold().MinWidth(0).SingleLine()
							})
						}
					})
				}
			})
		})
	})
	return root
}
