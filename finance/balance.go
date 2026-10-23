package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BalanceOptions are the balance visibility labels.
type BalanceOptions struct{ Show, Hide, Mask string }

// Balance hides and rolls balance digits; cue Toggle.
func Balance(c *ui.Context, amount float64, currency string, hidden *bool, o BalanceOptions) ui.Element {
	if hidden == nil || !finite(amount) || currency == "" {
		panic("finance: Balance has invalid data")
	}
	def(&o.Show, "Show balance")
	def(&o.Hide, "Hide balance")
	def(&o.Mask, "••••")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		if *hidden {
			ui.Text(c, o.Mask).Bold().FontSize(t.Rem(1.4)).Grow(1)
		} else {
			money(c, amount, currency).Bold().FontSize(t.Rem(1.4)).Grow(1).ContentTransition(ui.ContentTransition{Numeric: true})
		}
		name := o.Hide
		glyph := "eye-off"
		if *hidden {
			name, glyph = o.Show, "eye"
		}
		b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(8), t.Space(8)).Center().Radius(999))
		b.Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{}) })
		if b.Clicked() {
			*hidden = !*hidden
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("Balance", "visibility", "hidden", *hidden)
		}
	})
	return root
}
