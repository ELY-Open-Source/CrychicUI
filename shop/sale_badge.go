package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SaleBadgeOptions word the badge, "%d" taking the percent.
type SaleBadgeOptions struct{ Format string }

// SaleBadge is a tagged discount badge that wobbles in.
func SaleBadge(c *ui.Context, percent int, o SaleBadgeOptions) ui.Element {
	if percent <= 0 || percent > 100 {
		panic(fmt.Sprintf("shop: SaleBadge percent %d is outside 1..100", percent))
	}
	def(&o.Format, "-%d%%")
	t := c.Theme()
	text := fmt.Sprintf(o.Format, percent)
	badge := ui.Row(c).AlignSelf(ui.Start).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.75), t.Space(2.5)).
		Radius(999).Background(t.Danger).TextColor(t.AccentText).Role(ui.RoleStatus).Label(text)
	kit.Arrive(badge, ui.Pose{Scale: 0.6, Fade: 1}, 0)
	wobble := badge.Keyframes("crychic.wobble", []ui.Keyframe{
		{Value: -10}, {At: 140 * time.Millisecond, Value: 8, Ease: ui.EaseOut}, {At: 280 * time.Millisecond, Value: -4, Ease: ui.EaseInOut},
		{At: 420 * time.Millisecond, Value: 0, Ease: ui.EaseOut},
	})
	return badge.Rotate(wobble).Children(func() {
		ui.Icon(c, icon.Must("tag")).Size(12, 12)
		ui.Text(c, text).FontSize(t.Rem(0.8)).Bold().Shrink(0)
	})
}
