package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FlashSaleOptions are the sale countdown labels.
type FlashSaleOptions struct{ Format, Ended string }

// FlashSale shows a sale header with countdown.
func FlashSale(c *ui.Context, title string, ends time.Time, o FlashSaleOptions) ui.Element {
	if title == "" || ends.IsZero() {
		panic("shop: FlashSale needs a title and end time")
	}
	def(&o.Format, "%02d:%02d:%02d")
	def(&o.Ended, "Sale ended")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius).
		Background(t.Danger).TextColor(t.AccentText)
	kit.Arrive(root, ui.Pose{Y: -t.Space(2), Fade: 1}, 0)
	left := max(0, ends.Sub(c.Now()))
	if left > 0 {
		c.After(min(time.Second, left))
	}
	root.Children(func() {
		ui.Text(c, title).Bold().Grow(1)
		text := o.Ended
		if left > 0 {
			hours := int(left / time.Hour)
			minutes := int(left/time.Minute) % 60
			seconds := int(left/time.Second) % 60
			text = fmt.Sprintf(o.Format, hours, minutes, seconds)
		}
		kit.Numeric(c, text).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
	})
	return root
}
