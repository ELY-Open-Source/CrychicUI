package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UnsubscribeOptions word the form and the done note.
type UnsubscribeOptions struct{ Title, Confirm, Done string }

// Unsubscribe confirms a selected reason; cue Success.
func Unsubscribe(c *ui.Context, reasons []string, reason *string, confirmed *bool, o UnsubscribeOptions) (root, confirm ui.Element) {
	if len(reasons) == 0 || reason == nil || confirmed == nil {
		panic("mail: Unsubscribe needs reasons and state")
	}
	unique("Unsubscribe", reasons)
	def(&o.Title, "Why are you leaving?")
	def(&o.Confirm, "Unsubscribe")
	def(&o.Done, "You are unsubscribed")
	t := c.Theme()
	root = card(c)
	validReason := false
	for _, item := range reasons {
		validReason = validReason || *reason == item
	}
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(3), t.Space(1), t.Space(3)).Children(func() {
			tile(c, "mail-x", t.Danger)
			ui.Text(c, o.Title).FontSize(t.Rem(0.9)).Bold().Grow(1).MinWidth(0).Wrap()
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, item := range reasons {
				c.Key(item)
				on := *reason == item
				b := row(c, item, on).Disabled(*confirmed)
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						ring := ui.Box(c).Size(t.Space(4.5), t.Space(4.5)).Shrink(0).Center().Radius(999)
						ring.Border(1, ring.AnimateColor("ring", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
						dot := ring.AnimateTo("dot", unit(on), ui.Bouncy)
						ring.Children(func() {
							ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(999).Background(t.Accent).Scale(dot, dot)
						})
						ui.Text(c, item).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
					})
				})
				if b.Clicked() && !*confirmed && !on {
					*reason = item
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Unsubscribe", "reason", "index", i)
				}
			}
		})
		ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3), t.Space(3), t.Space(3)).Children(func() {
			if *confirmed {
				note := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).
					Background(t.Success.Alpha(0.12)).TextColor(t.Success)
				kit.Arrive(note, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
					glyph(c, "circle-check", 16)
					ui.Text(c, o.Done).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).Wrap()
				})
			}
			c.Key("confirm")
			confirm = kit.Press(ui.ButtonBase(c).Label(o.Confirm).FillWidth().Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
				TextColor(t.AccentText).Disabled(!validReason || *confirmed))
			confirm.Background(confirm.AnimateColor("bg", t.Danger.Mix(t.Text, 0.12*unit(confirm.Hovered())), ui.Snappy))
			confirm.Children(func() { ui.Text(c, o.Confirm).FontSize(t.Rem(0.85)).Bold() })
		})
	})
	if confirm.Clicked() && validReason && !*confirmed {
		*confirmed = true
		confirm.SetChanged()
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("Unsubscribe", "confirmed")
	}
	return root, confirm
}
