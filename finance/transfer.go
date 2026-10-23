package finance

import (
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TransferOptions are the transfer labels.
type TransferOptions struct{ From, To, Amount, Confirm, Invalid string }

// Transfer confirms an account transfer; cue Send.
func Transfer(c *ui.Context, from, to string, amount *string, currency string, confirmed *bool, o TransferOptions) (root, confirm ui.Element) {
	if from == "" || to == "" || amount == nil || currency == "" || confirmed == nil {
		panic("finance: Transfer has invalid data")
	}
	def(&o.From, "From")
	def(&o.To, "To")
	def(&o.Amount, "Amount")
	def(&o.Confirm, "Confirm")
	def(&o.Invalid, "Enter a positive amount")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			ui.Text(c, o.From).Grow(1)
			ui.Text(c, from)
			ui.Text(c, o.To).Grow(1)
			ui.Text(c, to)
		})
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			in := ui.TextInput(c, amount).Label(o.Amount).Grow(1)
			if in.Changed() {
				root.SetChanged()
			}
			confirm = kit.Press(ui.ButtonBase(c).Label(o.Confirm).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
				Background(t.Accent).TextColor(t.AccentText))
			confirm.Children(func() { ui.Text(c, o.Confirm).SingleLine() })
		})
		value, err := strconv.ParseFloat(strings.TrimSpace(*amount), 64)
		valid := err == nil && value > 0 && finite(value)
		confirm.Disabled(!valid)
		if !valid {
			ui.Text(c, o.Invalid).FontSize(t.Rem(0.8)).TextColor(t.Danger)
		}
	})
	if confirm.Clicked() {
		value, err := strconv.ParseFloat(strings.TrimSpace(*amount), 64)
		if err != nil || value <= 0 || !finite(value) {
			panic("finance: Transfer confirmed an invalid amount")
		}
		*confirmed = true
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("Transfer", "confirmed")
	}
	return root, confirm
}
