package finance

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TransferOptions word the card; SentTo takes amount and account.
type TransferOptions struct{ Title, From, To, Amount, Confirm, Invalid, Done, SentTo string }

// Transfer moves an amount between accounts on confirm; cue Send.
func Transfer(c *ui.Context, from, to string, amount *string, currency string, confirmed *bool, o TransferOptions) (root, confirm ui.Element) {
	if from == "" || to == "" || amount == nil || currency == "" || confirmed == nil {
		panic("finance: Transfer has invalid data")
	}
	def(&o.Title, "Transfer")
	def(&o.From, "From")
	def(&o.To, "To")
	def(&o.Amount, "Amount")
	def(&o.Confirm, "Confirm")
	def(&o.Invalid, "Enter a positive amount")
	def(&o.Done, "Sent")
	def(&o.SentTo, "%s sent to %s")
	t := c.Theme()
	account := func(label, name string, ink ui.Color) {
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1.5)).Padding(t.Space(2.5)).Radius(t.Radius).Background(t.Surface).
			Children(func() {
				kit.Muted(c, label).FontSize(t.Rem(0.7))
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					tile(c, "landmark", ink).Size(t.Space(8), t.Space(8))
					ui.Text(c, name).Bold().FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
				})
			})
	}
	root = kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	value, err := 0.0, error(nil)
	root.Children(func() {
		kit.CardHead(c, "arrow-left-right", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				account(o.From, from, t.Accent)
				ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Border(1, t.Border).Children(func() {
					ui.Icon(c, icon.Must("arrow-right")).Size(14, 14)
				})
				account(o.To, to, t.Success)
			})
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				kit.Muted(c, currency).Bold().Shrink(0)
				in := ui.TextInput(c, amount).Label(o.Amount).FontSize(t.Rem(1.2)).Grow(1).Disabled(*confirmed)
				if in.Changed() {
					root.SetChanged()
					kit.Log("Transfer", "edited", "runes", utf8.RuneCountInString(*amount))
				}
			})
			value, err = strconv.ParseFloat(strings.TrimSpace(*amount), 64)
			valid := err == nil && value > 0 && finite(value)
			if !valid && strings.TrimSpace(*amount) != "" {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.Danger).Children(func() {
					ui.Icon(c, icon.Must("circle-alert")).Size(14, 14)
					ui.Text(c, o.Invalid).FontSize(t.Rem(0.8)).Shrink(0)
				})
			}
			if *confirmed && valid {
				c.Key("sent")
				banner := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
					Background(t.Success.Alpha(0.1)).TextColor(t.Success)
				kit.Arrive(banner, ui.Pose{Scale: 0.95, Fade: 1}, 0).Children(func() {
					ui.Icon(c, icon.Must("circle-check")).Size(16, 16)
					ui.Text(c, fmt.Sprintf(o.SentTo, moneyText(c, value, currency), to)).Bold().FontSize(t.Rem(0.85)).
						Grow(1).MinWidth(0).Wrap()
				})
			}
			text, glyph := o.Confirm, "send"
			if *confirmed {
				text, glyph = o.Done, "check"
			}
			confirm = kit.Press(ui.ButtonBase(c).Label(o.Confirm).FillWidth().Center().Padding(t.Space(2.5)).Radius(t.Radius).
				TextColor(t.AccentText)).Disabled(!valid || *confirmed)
			confirm.Background(confirm.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(confirm.Hovered())).
				Mix(t.Success, unit(*confirmed)), ui.Snappy))
			confirm.Children(func() {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
					ui.Text(c, text).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
				})
			})
		})
	})
	if confirm.Clicked() && !*confirmed {
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
