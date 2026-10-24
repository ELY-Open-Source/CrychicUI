package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SavingsGoalOptions word the card; Icon is a Lucide name.
type SavingsGoalOptions struct{ Icon, Target, Of, Left, Reached string }

// SavingsGoal shows saved progress toward a target date.
func SavingsGoal(c *ui.Context, name string, saved, target float64, date, currency string, o SavingsGoalOptions) ui.Element {
	if name == "" || saved < 0 || target <= 0 || !finite(saved) || !finite(target) || date == "" || currency == "" {
		panic(fmt.Sprintf("finance: SavingsGoal %q needs saved of zero or more, a positive target and a date", name))
	}
	def(&o.Icon, "piggy-bank")
	def(&o.Target, "Target %s")
	def(&o.Of, "of %s")
	def(&o.Left, "%s to go")
	def(&o.Reached, "Goal reached")
	share := min(1, saved/target)
	done := saved >= target
	t := c.Theme()
	ink := t.Accent.Mix(t.Success, unit(done))
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, o.Icon, ink)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).Bold().FillWidth().SingleLine()
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					ui.Icon(c, icon.Must("calendar")).Size(12, 12)
					kit.Muted(c, fmt.Sprintf(o.Target, date)).FontSize(t.Rem(0.75)).Shrink(0)
				})
			})
			kit.Numeric(c, fmt.Sprintf("%.0f%%", 100*share)).Bold().FontSize(t.Rem(1.1)).TextColor(ink).Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
		})
		ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.End).Children(func() {
			money(c, saved, currency).Bold().FontSize(t.Rem(1.5)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
			kit.Muted(c, fmt.Sprintf(o.Of, moneyText(c, target, currency))).FontSize(t.Rem(0.8)).Padding(0, 0, t.Space(0.5), 0).Shrink(0)
		})
		bar := ui.Box(c).FillWidth().Height(t.Space(2.5)).Radius(999).Background(t.Surface)
		fill := bar.AnimateTo("saved", float32(share), ui.Smooth)
		tint := bar.AnimateColor("ink", ink, ui.Snappy)
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			if fill > 0 {
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*fill), H: r.H}, tint, 999)
			}
		}))
		if done {
			c.Key("reached")
			kit.Arrive(kit.Chip(c, o.Reached, t.Success).AlignSelf(ui.Start), ui.Pose{Scale: 0.8, Fade: 1}, 0)
		} else {
			c.Key("left")
			kit.Muted(c, fmt.Sprintf(o.Left, moneyText(c, target-saved, currency))).FontSize(t.Rem(0.8)).Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
		}
	})
	return root
}
