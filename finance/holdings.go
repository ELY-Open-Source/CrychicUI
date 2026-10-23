package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Holding is one market position.
type Holding struct {
	Symbol        string
	Shares, Price float64
	Change        float64
}

// Holdings shows positions with value and change.
func Holdings(c *ui.Context, positions []Holding, currency string) ui.Element {
	if len(positions) == 0 || currency == "" {
		panic("finance: Holdings needs positions and currency")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, position := range positions {
			if position.Symbol == "" || position.Shares < 0 || position.Price < 0 || !finite(position.Shares) ||
				!finite(position.Price) || !finite(position.Change) {
				panic("finance: Holdings has an invalid position")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, position.Symbol).Bold().Grow(1)
				money(c, position.Shares*position.Price, currency)
				ink := t.Success
				if position.Change < 0 {
					ink = t.Danger
				}
				kit.Numeric(c, fmt.Sprintf("%+.2f%%", position.Change)).TextColor(ink)
			})
		}
	})
	return root
}
