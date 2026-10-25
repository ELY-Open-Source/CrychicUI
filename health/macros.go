package health

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Macro is grams eaten of a nutrient against a goal.
type Macro struct{ Grams, Goal int }

// MacrosOptions word the card; Names are protein, carbs and fat.
type MacrosOptions struct {
	Title, Of, Share string
	Names            [3]string
}

// Macros shows protein, carbs and fat with an energy split.
func Macros(c *ui.Context, macros [3]Macro, o MacrosOptions) ui.Element {
	def(&o.Title, "Macros")
	energy := 0.0
	kcal := [3]float64{4, 4, 9}
	for i, m := range macros {
		if m.Grams < 0 || m.Goal < 1 || m.Grams > 1e5 || m.Goal > 1e5 {
			panic(fmt.Sprintf("health: Macros %q macro %d needs grams of 0 to 100,000 and a goal", o.Title, i))
		}
		energy += float64(m.Grams) * kcal[i]
	}
	def(&o.Of, "%dg of %dg")
	def(&o.Share, "%d%%")
	for i, n := range [...]string{"Protein", "Carbs", "Fat"} {
		def(&o.Names[i], n)
	}
	t := c.Theme()
	inks := [3]ui.Color{t.Accent, t.Warning, t.Danger.Mix(t.Warning, 0.4)}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Title).FontSize(t.Rem(0.9)).Bold().FillWidth().MinWidth(0).SingleLine()
		split := ui.Row(c).FillWidth().Height(t.Space(3)).Gap(t.Space(0.5))
		split.Children(func() {
			for i, m := range macros {
				if energy == 0 || m.Grams == 0 {
					continue
				}
				share := float64(m.Grams) * kcal[i] / energy
				c.Key(i)
				seg := ui.Box(c).Grow(float32(share)).Basis(0).FillHeight().Radius(999).Background(inks[i])
				kit.Arrive(seg, ui.Pose{Fade: 1}, 0)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(2.5))
		list.Children(func() {
			for i, m := range macros {
				pct := 0
				if energy > 0 {
					pct = int(math.Round(float64(m.Grams) * kcal[i] / energy * 100))
				}
				col := ui.Column(c).FillWidth().Gap(t.Space(1))
				kit.Stagger(c, list, col, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(inks[i])
						ui.Text(c, o.Names[i]).FontSize(t.Rem(0.82)).Bold().Grow(1).MinWidth(0).SingleLine()
						kit.Muted(c, fmt.Sprintf(o.Share, pct)).FontSize(t.Rem(0.7)).Shrink(0)
						kit.Numeric(c, fmt.Sprintf(o.Of, m.Grams, m.Goal)).FontSize(t.Rem(0.76)).Bold().Shrink(0)
					})
					over := m.Grams > m.Goal
					ink := inks[i]
					if over {
						ink = t.Danger
					}
					meter(c, float32(min(1, float64(m.Grams)/float64(m.Goal))), ink)
				})
			}
		})
	})
	return root
}
