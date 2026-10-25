package health

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Dish is a logged meal with its nutrients in grams.
type Dish struct {
	Name, Kind, At      string
	Calories            int
	Protein, Carbs, Fat int
	Items               []string
	Photo               ui.ImageSource
}

// MealOptions word the card; Names are protein, carbs and fat.
type MealOptions struct {
	Kcal, Grams string
	Names       [3]string
}

// Meal shows a meal with its photo, energy and nutrients.
func Meal(c *ui.Context, d Dish, o MealOptions) ui.Element {
	if strings.TrimSpace(d.Name) == "" || d.Calories < 0 || d.Calories > 1e5 || d.Protein < 0 || d.Carbs < 0 || d.Fat < 0 ||
		d.Protein > 1e4 || d.Carbs > 1e4 || d.Fat > 1e4 {
		panic(fmt.Sprintf("health: Meal %q needs a name and nutrients of 0 or more", d.Name))
	}
	def(&o.Kcal, "%s kcal")
	def(&o.Grams, "%dg")
	for i, n := range [...]string{"Protein", "Carbs", "Fat"} {
		def(&o.Names[i], n)
	}
	t := c.Theme()
	inks := [3]ui.Color{t.Accent, t.Warning, t.Danger.Mix(t.Warning, 0.4)}
	grams := [3]int{d.Protein, d.Carbs, d.Fat}
	root := card(c).Clip().Role(ui.RoleGroup)
	root.Children(func() {
		band := ui.Box(c).FillWidth().Height(t.Space(30)).Clip().Background(t.Warning.Alpha(0.14))
		band.Children(func() {
			if d.Photo != nil {
				ui.Image(c, d.Photo).Fill().Fit(ui.Cover)
			} else {
				ui.Box(c).Fill().Center().TextColor(t.Warning).Children(func() { glyph(c, "utensils", 36) })
			}
			if d.Kind != "" {
				ui.Box(c).Absolute().Left(t.Space(3)).Top(t.Space(3)).Children(func() { kit.Chip(c, d.Kind, t.Warning) })
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3.5)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, d.Name).FontSize(t.Rem(0.95)).Bold().FillWidth().Wrap()
					if d.At != "" {
						kit.Muted(c, d.At).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
					}
				})
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Shrink(0).TextColor(t.Danger).Children(func() {
					glyph(c, "flame", 14)
					kit.Numeric(c, fmt.Sprintf(o.Kcal, number(c, d.Calories))).FontSize(t.Rem(0.85)).Bold().Shrink(0)
				})
			})
			chips := ui.Row(c).FillWidth().Gap(t.Space(2))
			chips.Children(func() {
				for i, g := range grams {
					cell := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(0.5)).Padding(t.Space(2)).Radius(t.Radius).
						Background(inks[i].Alpha(0.1))
					kit.Stagger(c, chips, cell, i).Children(func() {
						kit.Numeric(c, fmt.Sprintf(o.Grams, g)).FontSize(t.Rem(0.9)).Bold().TextColor(inks[i]).FillWidth().MinWidth(0).SingleLine()
						kit.Muted(c, o.Names[i]).FontSize(t.Rem(0.66)).FillWidth().MinWidth(0).SingleLine()
					})
				}
			})
			for _, it := range d.Items {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(t.Space(1.25), t.Space(1.25)).Shrink(0).Radius(999).Background(t.TextMuted)
					ui.Text(c, it).FontSize(t.Rem(0.78)).Grow(1).MinWidth(0).Wrap()
				})
			}
		})
	})
	return root
}
