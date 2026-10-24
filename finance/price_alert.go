package finance

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AlertDirection selects the alert threshold direction.
type AlertDirection uint8

// Alert directions.
const (
	Above AlertDirection = iota
	Below
)

// PriceAlertOptions word the alert; each takes a price.
type PriceAlertOptions struct{ Above, Below, Now, Triggered, Watching string }

// PriceAlert rings its bell when a price crosses; cue Notify.
func PriceAlert(c *ui.Context, direction AlertDirection, target, current float64, currency string, active *bool, o PriceAlertOptions) ui.Element {
	if direction > Below || target <= 0 || current < 0 || !finite(target) || !finite(current) || currency == "" || active == nil {
		panic(fmt.Sprintf("finance: PriceAlert needs a direction, a positive target and bound state, got %v", target))
	}
	def(&o.Above, "Price above %s")
	def(&o.Below, "Price below %s")
	def(&o.Now, "Now %s")
	def(&o.Triggered, "Price alert triggered")
	def(&o.Watching, "Watching")
	t := c.Theme()
	triggered := direction == Above && current >= target || direction == Below && current <= target
	near := current / target
	if direction == Below {
		near = 1
		if current > 0 {
			near = target / current
		}
	}
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.Warning.Alpha(0.1), unit(triggered)), ui.Snappy))
	ink := root.AnimateColor("ink", t.Accent.Mix(t.Warning, unit(triggered)), ui.Snappy)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			bell := ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(999).Background(ink.Alpha(0.14)).TextColor(ink)
			bell.Children(func() {
				glyph := "bell"
				if triggered {
					glyph = "bell-ring"
				}
				c.Key(triggered)
				ring := ui.Icon(c, icon.Must(glyph)).Size(18, 18)
				if triggered {
					swing := ring.Keyframes("crychic.ring", []ui.Keyframe{{Value: 0}, {At: 120 * time.Millisecond, Value: 18},
						{At: 240 * time.Millisecond, Value: -14}, {At: 360 * time.Millisecond, Value: 8}, {At: 480 * time.Millisecond, Value: 0}})
					ring.Rotate(swing)
				}
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				rule := o.Above
				if direction == Below {
					rule = o.Below
				}
				ui.Text(c, fmt.Sprintf(rule, moneyText(c, target, currency))).Bold().FontSize(t.Rem(0.9)).FillWidth().Wrap()
				kit.Muted(c, fmt.Sprintf(o.Now, moneyText(c, current, currency))).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
			})
			if triggered {
				c.Key("fired")
				kit.Arrive(kit.Chip(c, o.Triggered, t.Warning), ui.Pose{Scale: 0.8, Fade: 1}, 0)
			} else {
				c.Key("watching")
				kit.Chip(c, o.Watching, t.TextMuted)
			}
		})
		bar := ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(999).Background(t.Surface)
		fill := bar.AnimateTo("near", float32(min(1, near)), ui.Smooth)
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			if fill > 0 {
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*fill), H: r.H}, ink, 999)
			}
		}))
	})
	if triggered && !*active {
		*active = true
		root.SetChanged()
		sound.Play(sound.Notify)
		kit.Log("PriceAlert", "triggered")
	}
	return root
}
