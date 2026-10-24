package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FlashSaleOptions word the banner; Hours, Minutes, Seconds name the tiles.
type FlashSaleOptions struct {
	Subtitle, Ended         string
	Hours, Minutes, Seconds string
}

// FlashSale is a sale banner counting down in rolling tiles.
func FlashSale(c *ui.Context, title string, ends time.Time, o FlashSaleOptions) ui.Element {
	if title == "" || ends.IsZero() {
		panic(fmt.Sprintf("shop: FlashSale %q needs a title and an end time", title))
	}
	def(&o.Ended, "Sale ended")
	def(&o.Hours, "hrs")
	def(&o.Minutes, "min")
	def(&o.Seconds, "sec")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius).
		Background(t.Danger).TextColor(t.AccentText).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: -t.Space(2), Fade: 1}, 0)
	now := c.Now()
	left := max(0, ends.Sub(now))
	// Whole seconds from Unix times: Sub saturates past 292 years.
	secs := ends.Unix() - now.Unix()
	frac := time.Duration((ends.Nanosecond() - now.Nanosecond() + 1e9) % 1e9)
	if ends.Nanosecond() > now.Nanosecond() {
		secs++
	}
	if left > 0 {
		if frac == 0 {
			frac = time.Second
		}
		c.After(min(left, frac+time.Millisecond))
	}
	root.Children(func() {
		ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(999).Background(t.AccentText.Alpha(0.18)).
			Children(func() { ui.Icon(c, icon.Must("zap")).Size(18, 18) })
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, title).Bold().FontSize(t.Rem(1.05)).FillWidth().SingleLine()
			if o.Subtitle != "" {
				ui.Text(c, o.Subtitle).FontSize(t.Rem(0.8)).TextColor(t.AccentText.Alpha(0.85)).FillWidth().Wrap()
			}
		})
		if left == 0 {
			ui.Text(c, o.Ended).Bold().Shrink(0).Padding(t.Space(1), t.Space(2.5)).Radius(999).Background(t.AccentText.Alpha(0.18))
			return
		}
		parts := [3]int64{secs / 3600, secs / 60 % 60, secs % 60}
		units := [3]string{o.Hours, o.Minutes, o.Seconds}
		ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Start).Children(func() {
			for i, v := range parts {
				c.Key(i)
				if i > 0 {
					ui.Text(c, ":").Bold().Padding(t.Space(1), 0, 0, 0)
				}
				ui.Column(c).Gap(t.Space(0.5)).AlignItems(ui.Center).Children(func() {
					kit.Numeric(c, fmt.Sprintf("%02d", v)).Bold().FontSize(t.Rem(1.05)).MinWidth(t.Space(9)).TextAlign(ui.Center).
						Padding(t.Space(1), t.Space(1.5)).Radius(t.Space(1.5)).Background(t.AccentText.Alpha(0.18)).
						ContentTransition(ui.ContentTransition{Numeric: true})
					ui.Text(c, units[i]).FontSize(t.Rem(0.65)).TextColor(t.AccentText.Alpha(0.85)).Shrink(0)
				})
			}
		})
	})
	return root
}
