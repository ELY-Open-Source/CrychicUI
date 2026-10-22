package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FacetOption is an option and its count.
type FacetOption struct {
	Label string
	Count int
}

// Facet is a titled group of options.
type Facet struct {
	Title   string
	Options []FacetOption
}

// Facets checks options; chosen holds "title/label".
func Facets(c *ui.Context, label string, facets []Facet, chosen *ui.Selection[string]) ui.Element {
	if label == "" || len(facets) == 0 || chosen == nil {
		panic(fmt.Sprintf("data: Facets %q needs facets and a choice", label))
	}
	var titles []string
	for _, f := range facets {
		titles = append(titles, f.Title)
		var names []string
		for _, op := range f.Options {
			if op.Count < 0 {
				panic(fmt.Sprintf("data: Facets %q counts must not be negative", label))
			}
			names = append(names, op.Label)
		}
		named("Facets", label, names)
	}
	named("Facets", label, titles)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for _, f := range facets {
			c.Key(f.Title)
			ui.Column(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
				kit.Muted(c, f.Title).Bold().FontSize(t.Rem(0.75))
				for _, op := range f.Options {
					key := f.Title + "/" + op.Label
					on := chosen.Has(key)
					c.Key(op.Label)
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						if checkbox(c, op.Label, &on, false).Changed() {
							toggle(chosen, key)
							root.SetChanged()
							sound.Play(sound.Toggle)
							kit.Log("Facets", "toggled", "on", on)
						}
						ui.Text(c, op.Label).Grow(1).SingleLine()
						motion.RollingDigits(c, int64(op.Count)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
					})
				}
			})
		}
	})
	return root
}
