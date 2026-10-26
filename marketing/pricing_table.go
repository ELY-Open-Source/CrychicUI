package marketing

import (
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PricingTableOptions word the table; Yearly shows yearly prices.
type PricingTableOptions struct {
	PlanWords
	Yearly bool
}

// PricingTable sets plans side by side; cue Tap.
func PricingTable(c *ui.Context, plans []Plan, picked *string, o PricingTableOptions) ui.Element {
	names := make([]string, len(plans))
	for i, p := range plans {
		validPlan("PricingTable", p)
		names[i] = p.Name
	}
	unique("PricingTable", names)
	if len(plans) == 0 || picked == nil || *picked != "" && !slices.Contains(names, *picked) {
		panic("marketing: PricingTable needs plans and a picked name among them")
	}
	o.fill()
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Stretch).Role(ui.RoleGroup)
	root.Children(func() {
		for i, p := range plans {
			c.Key(strconv.Quote(p.Name))
			ui.Column(c).Grow(1).Basis(t.Space(40)).MinWidth(0).Children(func() {
				_, choose := planCard(c, p, o.Yearly, p.Name == *picked, o.PlanWords)
				if choose.Clicked() && p.Name != *picked {
					*picked = p.Name
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("PricingTable", "picked", "index", i)
				}
			})
		}
	})
	return root
}
