package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CartIconOptions are the cart labels.
type CartIconOptions struct{ Label string }

// CartIcon shows a bag with a rolling count.
func CartIcon(c *ui.Context, count int, o CartIconOptions) ui.Element {
	if count < 0 {
		panic("shop: CartIcon count must not be negative")
	}
	def(&o.Label, "Cart")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(o.Label)
	root.Children(func() {
		ui.Icon(c, icon.Must("shopping-bag")).Size(t.Space(6), t.Space(6))
		kit.Numeric(c, fmt.Sprint(count)).Bold()
	})
	return root
}
