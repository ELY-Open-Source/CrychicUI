package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SplitBillOptions are the split labels.
type SplitBillOptions struct{ People, Each, Decrease, Increase string }

// SplitBill splits a total among people; cue Tick.
func SplitBill(c *ui.Context, total float64, people *int, currency string, o SplitBillOptions) ui.Element {
	if total < 0 || !finite(total) || people == nil || *people < 1 || *people > 100 || currency == "" {
		panic("finance: SplitBill has invalid data")
	}
	def(&o.People, "People")
	def(&o.Each, "Per person")
	def(&o.Decrease, "Decrease")
	def(&o.Increase, "Increase")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.People).Grow(1)
			for _, action := range []struct {
				name, glyph string
				delta       int
			}{{o.Decrease, "minus", -1}, {o.Increase, "plus", 1}} {
				b := kit.Press(ui.ButtonBase(c).Label(action.name).Size(t.Space(6), t.Space(6)).Center().Radius(999).
					Disabled(action.delta < 0 && *people == 1 || action.delta > 0 && *people == 100))
				b.Children(func() { ui.Icon(c, icon.Must(action.glyph)).Size(16, 16) })
				if b.Clicked() && (action.delta < 0 && *people > 1 || action.delta > 0 && *people < 100) {
					*people += action.delta
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("SplitBill", "people", "value", *people)
				}
			}
			kit.Numeric(c, fmt.Sprint(*people))
		})
		ui.Row(c).FillWidth().Children(func() {
			ui.Text(c, o.Each).Grow(1)
			money(c, total/float64(*people), currency).Bold()
		})
	})
	return root
}
