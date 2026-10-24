package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BalanceOptions word the card; Change sits under the amount.
type BalanceOptions struct{ Label, Show, Hide, Mask, Change string }

// Balance shows rolling digits with a hide toggle; cue Toggle.
func Balance(c *ui.Context, amount float64, currency string, hidden *bool, o BalanceOptions) ui.Element {
	if hidden == nil || !finite(amount) || currency == "" {
		panic(fmt.Sprintf("finance: Balance needs a finite amount, a currency and bound state, got %v", amount))
	}
	def(&o.Label, "Available balance")
	def(&o.Show, "Show balance")
	def(&o.Hide, "Hide balance")
	def(&o.Mask, "••••••")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	var b ui.Element
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			tile(c, "wallet", t.Accent)
			kit.Muted(c, o.Label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
			name, glyph := o.Hide, "eye-off"
			if *hidden {
				name, glyph = o.Show, "eye"
			}
			b = kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(9), t.Space(9)).Center().Radius(999).TextColor(t.TextMuted))
			b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
			b.Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{}) })
		})
		if *hidden {
			c.Key("mask")
			kit.Arrive(ui.Text(c, o.Mask).Bold().FontSize(t.Rem(2)), ui.Pose{Scale: 0.9, Fade: 1}, 0)
		} else {
			c.Key("amount")
			value := money(c, amount, currency).Bold().FontSize(t.Rem(2)).ContentTransition(ui.ContentTransition{Numeric: true})
			if amount < 0 {
				value.TextColor(t.Danger)
			}
			kit.Arrive(value, ui.Pose{Scale: 0.9, Fade: 1}, 0)
		}
		if o.Change != "" {
			kit.Muted(c, o.Change).FontSize(t.Rem(0.8)).FillWidth().Wrap()
		}
	})
	if b.Clicked() {
		*hidden = !*hidden
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Balance", "visibility", "hidden", *hidden)
	}
	return root
}
