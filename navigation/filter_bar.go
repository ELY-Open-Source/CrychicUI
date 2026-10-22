package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FilterBarOptions tune a FilterBar.
type FilterBarOptions struct {
	// Clear names the clear-all button.
	Clear string
}

// FilterBar is filter chips with a clear-all.
func FilterBar(c *ui.Context, label string, options []string, selected *[]string, o FilterBarOptions) ui.Element {
	if selected == nil {
		panic(fmt.Sprintf("navigation: FilterBar %q needs a selection", label))
	}
	def(&o.Clear, "Clear all")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
	return root.Children(func() {
		if choice.FilterChips(c, label, options, selected).Changed() {
			root.SetChanged()
			kit.Log("FilterBar", "filtered", "count", len(*selected))
		}
		if len(*selected) == 0 {
			return
		}
		c.Key("clear")
		ui.Box(c).Shrink(0).Transition(kit.Enter).Children(func() {
			if button.Ghost(c, o.Clear).Clicked() {
				*selected = nil
				root.SetChanged()
				kit.Log("FilterBar", "cleared")
			}
		})
	})
}
