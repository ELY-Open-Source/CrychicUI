package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Stock is a SKU, a name, a count and capacity.
type Stock struct {
	SKU, Name  string
	Count, Max int
}

// InventoryOptions tune an Inventory.
type InventoryOptions struct {
	// Low turns bars red; zero gives 0.2.
	Low float64
	// Of takes the count and capacity.
	Of string
}

// Inventory shows SKU rows with stock bars colored by level.
func Inventory(c *ui.Context, label string, stock []Stock, o InventoryOptions) ui.Element {
	if label == "" || len(stock) == 0 || !finite64(o.Low) || o.Low < 0 || o.Low >= 0.5 {
		panic(fmt.Sprintf("data: Inventory %q needs stock and Low below 0.5", label))
	}
	var skus []string
	for _, s := range stock {
		if s.Max <= 0 || s.Count < 0 || s.Count > s.Max {
			panic(fmt.Sprintf("data: Inventory %q counts need 0 to a positive Max", label))
		}
		skus = append(skus, s.SKU)
	}
	unique("Inventory", label, skus)
	if o.Low == 0 {
		o.Low = 0.2
	}
	def(&o.Of, "%d / %d")
	t := c.Theme()
	root := frame(c, label).Dividers(1, t.Border)
	return root.Children(func() {
		for _, s := range stock {
			c.Key(s.SKU)
			frac := float64(s.Count) / float64(s.Max)
			ink := t.Success
			switch {
			case frac <= o.Low:
				ink = t.Danger
			case frac <= 2*o.Low:
				ink = t.Warning
			}
			hovered(c, false).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
				kit.Muted(c, s.SKU).FontFeatures("tnum").Width(t.Space(24)).Shrink(0).SingleLine()
				ui.Text(c, s.Name).Grow(1).SingleLine()
				ui.Row(c).Width(t.Space(36)).Shrink(0).Children(func() { bar(c, s.Name, float32(frac), ink) })
				kit.Numeric(c, fmt.Sprintf(o.Of, s.Count, s.Max)).Width(t.Space(16)).Shrink(0).TextAlign(ui.End)
			})
		}
	})
}
