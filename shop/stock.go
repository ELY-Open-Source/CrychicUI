package shop

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StockOptions word stock states; Low may take the count.
type StockOptions struct{ In, Low, Out string }

// Stock is a colored pill: in stock, low or out.
func Stock(c *ui.Context, count int, o StockOptions) ui.Element {
	if count < 0 {
		panic(fmt.Sprintf("shop: Stock count %d must not be negative", count))
	}
	def(&o.In, "In stock")
	def(&o.Low, "Low stock")
	def(&o.Out, "Out of stock")
	t := c.Theme()
	text, ink := o.In, t.Success
	switch {
	case count == 0:
		text, ink = o.Out, t.Danger
	case count < 5:
		text, ink = o.Low, t.Warning
		if strings.Contains(o.Low, "%d") {
			text = fmt.Sprintf(o.Low, count)
		}
	}
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(0.75), t.Space(2.5)).Radius(999).Role(ui.RoleStatus).Label(text)
	root.Background(root.AnimateColor("bg", ink.Alpha(0.12), ui.Snappy)).TextColor(root.AnimateColor("ink", ink, ui.Snappy))
	kit.Arrive(root, ui.Pose{Fade: 1, Scale: 0.9}, 0)
	root.Children(func() {
		ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
		ui.Text(c, text).FontSize(t.Rem(0.8)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
	})
	return root
}
