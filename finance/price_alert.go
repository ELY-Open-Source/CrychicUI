package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AlertDirection selects the alert threshold direction.
type AlertDirection uint8

const (
	Above AlertDirection = iota
	Below
)

// PriceAlertOptions are the alert labels.
type PriceAlertOptions struct{ Above, Below, Triggered string }

// PriceAlert triggers above or below a price; cue Notify.
func PriceAlert(c *ui.Context, direction AlertDirection, target, current float64, currency string, active *bool, o PriceAlertOptions) ui.Element {
	if direction > Below || target <= 0 || current < 0 || !finite(target) || !finite(current) || currency == "" || active == nil {
		panic("finance: PriceAlert has invalid data")
	}
	def(&o.Above, "Above")
	def(&o.Below, "Below")
	def(&o.Triggered, "Price alert triggered")
	t := c.Theme()
	triggered := direction == Above && current >= target || direction == Below && current <= target
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
		Background(t.Accent.Alpha(0.12 * unit(triggered))).Role(ui.RoleGroup)
	root.Children(func() {
		name := o.Above
		if direction == Below {
			name = o.Below
		}
		ui.Text(c, name).Grow(1)
		money(c, target, currency)
		money(c, current, currency)
		if triggered {
			ui.Text(c, o.Triggered).TextColor(t.Warning)
		}
	})
	if triggered && !*active {
		*active = true
		root.SetChanged()
		sound.Play(sound.Notify)
		kit.Log("PriceAlert", "triggered")
	}
	return root
}
