package form

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Facet is titled options with their result counts.
type Facet struct {
	Title   string
	Options []string
	Counts  []int
}

// Filter is one chosen option of a facet.
type Filter struct{ Facet, Option string }

// FilterPanelOptions tune a FilterPanel.
type FilterPanelOptions struct {
	// Clear is the clear button's copy.
	Clear string
}

// FilterPanel checks facet options with counts and clears them.
func FilterPanel(c *ui.Context, facets []Facet, selected *[]Filter, o FilterPanelOptions) ui.Element {
	if len(facets) == 0 || selected == nil {
		panic("form: FilterPanel needs facets and a selection")
	}
	for _, f := range facets {
		if f.Title == "" || len(f.Options) == 0 || len(f.Counts) != len(f.Options) {
			panic(fmt.Sprintf("form: FilterPanel facet %q needs options with counts", f.Title))
		}
	}
	def(&o.Clear, "Clear filters")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4))
	return root.Children(func() {
		for _, f := range facets {
			kit.Arrive(ui.Column(c).Key(f.Title).FillWidth().Gap(t.Space(2)), ui.Pose{Fade: 1}, 0).Children(func() {
				ui.Text(c, f.Title).FontWeight(600)
				for j, opt := range f.Options {
					pick := Filter{f.Title, opt}
					on := slices.Contains(*selected, pick)
					ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
						if choice.Check(c, opt, &on).Grow(1).Changed() {
							if on {
								*selected = append(*selected, pick)
							} else {
								*selected = slices.DeleteFunc(*selected, func(p Filter) bool { return p == pick })
							}
							root.SetChanged()
							kit.Log("FilterPanel", "filters", "count", len(*selected))
						}
						kit.Muted(c, fmt.Sprint(f.Counts[j]))
					})
				}
			})
		}
		if len(*selected) > 0 {
			clear := button.Ghost(c, o.Clear).Transition(ui.ElementTransition{Exit: &ui.Motion{}})
			if kit.Arrive(clear, ui.Pose{Fade: 1}, 0).Clicked() {
				*selected = nil
				root.SetChanged()
				kit.Log("FilterPanel", "cleared")
			}
		}
	})
}
