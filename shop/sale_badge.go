package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SaleBadgeOptions are the sale badge labels.
type SaleBadgeOptions struct{ Format string }

// SaleBadge shows a discount badge; it wobbles in.
func SaleBadge(c *ui.Context, percent int, o SaleBadgeOptions) ui.Element {
	if percent <= 0 || percent > 100 {
		panic(fmt.Sprintf("shop: SaleBadge percent %d is invalid", percent))
	}
	def(&o.Format, "-%d%%")
	t := c.Theme()
	badge := ui.Text(c, fmt.Sprintf(o.Format, percent)).Padding(t.Space(1), t.Space(2)).Radius(999).
		Background(t.Danger).TextColor(t.AccentText).Bold()
	kit.Arrive(badge, ui.Pose{Scale: 0.7, Fade: 1}, 0)
	badge.Rotate(-8 * badge.Keyframes("wobble", []ui.Keyframe{{Value: 1}, {At: 200 * time.Millisecond, Value: 0, Ease: ui.EaseOut}}))
	return badge
}
