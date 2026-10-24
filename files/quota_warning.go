package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuotaWarningOptions word the banner.
type QuotaWarningOptions struct {
	// Format words fullness, "%.0f" taking the percent.
	Format, Title, Of, Upgrade, Upgraded string
}

// QuotaWarning is a near-full storage banner with Upgrade; cue Warning.
func QuotaWarning(c *ui.Context, used, total int64, upgraded *bool, o QuotaWarningOptions) (root, upgrade ui.Element) {
	if used < 0 || total <= 0 || used > total || upgraded == nil {
		panic("files: QuotaWarning needs usage within a positive total and bound state")
	}
	def(&o.Format, "%.0f%% full")
	def(&o.Title, "Storage almost full")
	def(&o.Of, "of")
	def(&o.Upgrade, "Upgrade")
	def(&o.Upgraded, "Plan upgraded")
	t := c.Theme()
	ink := t.Warning
	if *upgraded {
		ink = t.Success
	}
	root = ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).Radius(t.Radius).Role(ui.RoleStatus).Label(o.Title)
	root.Background(root.AnimateColor("bg", ink.Alpha(0.08), ui.Smooth)).Border(1, root.AnimateColor("line", ink.Alpha(0.4), ui.Smooth))
	kit.Arrive(root, ui.Pose{Y: -t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			glyph := "triangle-alert"
			if *upgraded {
				glyph = "circle-check"
			}
			ui.Box(c).Size(t.Space(9), t.Space(9)).Shrink(0).Center().Radius(999).Background(ink.Alpha(0.16)).TextColor(ink).
				Children(func() { ui.Icon(c, icon.Must(glyph)).Size(18, 18).ContentTransition(ui.ContentTransition{}) })
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				title := o.Title
				if *upgraded {
					title = o.Upgraded
				}
				ui.Text(c, title).Bold().ContentTransition(ui.ContentTransition{})
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					sizeText(c, used).FontSize(t.Rem(0.8))
					kit.Muted(c, o.Of).FontSize(t.Rem(0.8))
					sizeText(c, total).FontSize(t.Rem(0.8))
					kit.Muted(c, "·").FontSize(t.Rem(0.8))
					kit.Muted(c, fmt.Sprintf(o.Format, 100*float64(used)/float64(total))).FontSize(t.Rem(0.8)).Shrink(0)
				})
			})
			if !*upgraded {
				c.Key("upgrade")
				upgrade = kit.Press(ui.ButtonBase(c).Label(o.Upgrade).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius)).Transition(kit.Pop)
				upgrade.Background(upgrade.AnimateColor("bg", t.Warning.Mix(t.Text, 0.15*unit(upgrade.Hovered())), ui.Snappy)).TextColor(t.Background).
					Children(func() { ui.Text(c, o.Upgrade).Bold().Shrink(0) })
			}
		})
		progressBar(c, float64(used)/float64(total), ink)
	})
	if !*upgraded && upgrade.Clicked() {
		*upgraded = true
		root.SetChanged()
		sound.Play(sound.Warning)
		kit.Log("QuotaWarning", "upgraded")
	}
	return root, upgrade
}
