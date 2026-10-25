package project

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// WorkloadPerson is a load; 1 is full capacity.
type WorkloadPerson struct {
	Name, Role string
	Load       float64
}

// WorkloadOptions word the card.
type WorkloadOptions struct{ Title, Over, Balanced string }

// Workload shows people with load bars.
func Workload(c *ui.Context, people []WorkloadPerson, o WorkloadOptions) ui.Element {
	def(&o.Title, "Workload")
	def(&o.Over, "%d over capacity")
	def(&o.Balanced, "Balanced")
	if len(people) == 0 {
		panic(fmt.Sprintf("project: Workload %q needs people", o.Title))
	}
	names := make([]string, len(people))
	over := 0
	for i, person := range people {
		if person.Load < 0 || person.Load > 2 || !finite(person.Load) {
			panic(fmt.Sprintf("project: Workload %q has load %v for %q", o.Title, person.Load, person.Name))
		}
		names[i] = person.Name
		if person.Load > 1 {
			over++
		}
	}
	unique("Workload", names)
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "gauge", o.Title, func() {
			if over > 0 {
				kit.Chip(c, fmt.Sprintf(o.Over, over), t.Danger)
			} else {
				kit.Chip(c, o.Balanced, t.Success)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3))
		list.Children(func() {
			for i, person := range people {
				c.Key(person.Name)
				ink := t.Accent
				switch {
				case person.Load > 1:
					ink = t.Danger
				case person.Load >= 0.85:
					ink = t.Warning
				}
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center)
				kit.Stagger(c, list, row, i).Children(func() {
					display.Avatar(c, person.Name, nil).Size(t.Space(8), t.Space(8)).FontSize(t.Rem(0.65))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, person.Name).FontSize(t.Rem(0.8)).Bold().MinWidth(0).SingleLine()
							kit.Muted(c, person.Role).FontSize(t.Rem(0.7)).Grow(1).MinWidth(0).SingleLine()
							kit.Numeric(c, fmt.Sprintf("%.0f%%", math.Round(person.Load*100))).FontSize(t.Rem(0.75)).Bold().
								TextColor(row.AnimateColor("ink", ink, ui.Snappy)).Shrink(0)
						})
						progress(c, min(1, person.Load), ink)
					})
				})
			}
		})
	})
	return root
}
