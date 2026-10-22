package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// DateBadge shows month over day on a leaf that flips.
func DateBadge(c *ui.Context, day time.Time) ui.Element {
	if day.IsZero() {
		panic("datetime: DateBadge needs a day")
	}
	t := c.Theme()
	loc := locale(c)
	w, h := t.Space(14), t.Space(16)
	root := ui.Box(c).Size(w, h).Shrink(0).Role(ui.RoleText).
		Label(i18n.FormatDate(loc, day, i18n.DateOptions{Skeleton: "yMMMd"}))
	return root.Children(func() {
		y, m, d := day.Date()
		leaf := ui.Column(c).Key(y*10000+int(m)*100+d).Size(w, h).Radius(t.Radius).Background(t.Background).
			Border(1, t.Border).Shadow(0, 1, 2, 0, kit.Shade(0.12)).Clip().TransformOrigin(0.5, 0)
		k := leaf.Keyframes("crychic.flip", []ui.Keyframe{{Value: 0}, {At: 320 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		leaf.RotateX(-90 * (1 - k)).Perspective(t.Space(100)).Children(func() {
			ui.Box(c).FillWidth().Padding(t.Space(0.5)).Center().Background(t.Accent).TextColor(t.AccentText).Children(func() {
				ui.Text(c, i18n.Months(loc, i18n.Abbreviated, true)[m-1]).FontSize(t.Rem(0.75)).Bold()
			})
			ui.Box(c).FillWidth().Grow(1).Center().Children(func() {
				kit.Numeric(c, i18n.FormatDate(loc, day, i18n.DateOptions{Skeleton: "d"})).FontSize(t.Rem(1.75)).Bold()
			})
		})
	})
}
