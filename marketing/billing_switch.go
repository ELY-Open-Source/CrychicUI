package marketing

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BillingSwitchOptions word the switch; Save takes the percent saved.
type BillingSwitchOptions struct{ Monthly, Yearly, Save, Currency, Per, Billed string }

// BillingSwitch flips monthly and yearly billing; cue Toggle.
func BillingSwitch(c *ui.Context, monthly, perYear float64, yearly *bool, o BillingSwitchOptions) ui.Element {
	if yearly == nil || !(monthly >= 0.01 && monthly <= 1e9) || !(perYear >= 0.01 && perYear <= 1e10) {
		panic(fmt.Sprintf("marketing: BillingSwitch needs a yearly state and prices of at least 0.01, got %v and %v", monthly, perYear))
	}
	def(&o.Monthly, "Monthly")
	def(&o.Yearly, "Yearly")
	def(&o.Save, "Save %d%%")
	def(&o.Currency, "$")
	def(&o.Per, "/ month")
	def(&o.Billed, "%s billed yearly")
	t := c.Theme()
	// Whole cents, so 10 and 96 save exactly 20.
	year := 12 * int64(math.Round(monthly*100))
	save := int(max(0, (year-int64(math.Round(perYear*100)))*100/year))
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			pill := ui.Row(c).Shrink(0).Padding(t.Space(1)).Radius(999).Background(t.SurfaceHover)
			x := pill.AnimateTo("thumb", unit(*yearly), ui.Snappy)
			pill.Draw(func(p *ui.Painter, r ui.Rect) {
				pad := t.Space(1)
				w := (r.W - 2*pad) / 2
				p.Fill(ui.Rect{X: r.X + pad + x*w, Y: r.Y + pad, W: w, H: r.H - 2*pad}, t.Background, 999)
			})
			pill.Children(func() {
				for i, name := range [2]string{o.Monthly, o.Yearly} {
					on := (i == 1) == *yearly
					b := kit.Press(ui.ButtonBase(c).Label(name).Checked(on).Width(t.Space(24)).Padding(t.Space(1.5), t.Space(2)).Radius(999))
					b.Children(func() {
						ui.Text(c, name).FontSize(t.Rem(0.8)).Bold().TextColor(t.TextMuted.Mix(t.Text, unit(on || b.Hovered()))).
							FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
					})
					if b.Clicked() && !on {
						*yearly = i == 1
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("BillingSwitch", "switched", "yearly", *yearly)
					}
				}
			})
			if save >= 1 {
				c.Key("save")
				chip := kit.Chip(c, fmt.Sprintf(o.Save, save), t.Success)
				chip.Opacity(chip.AnimateTo("lit", 0.45+0.55*unit(*yearly), ui.Smooth))
			}
		})
		amount := monthly
		if *yearly {
			amount = perYear / 12
		}
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
			kit.Numeric(c, o.Currency+price(c, math.Round(amount*100)/100)).FontSize(t.Rem(2.2)).Bold().Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
			kit.Muted(c, o.Per).FontSize(t.Rem(0.82)).Shrink(0).Padding(0, 0, t.Space(1.25), 0)
		})
		note := ""
		if *yearly {
			note = fmt.Sprintf(o.Billed, o.Currency+price(c, perYear))
		}
		c.Key("note")
		kit.Muted(c, note).FontSize(t.Rem(0.76)).FillWidth().TextAlign(ui.Center).Wrap().ContentTransition(ui.ContentTransition{})
	})
	return root
}
