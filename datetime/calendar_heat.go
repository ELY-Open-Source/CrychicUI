package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// CalendarHeatOptions tune a CalendarHeat.
type CalendarHeatOptions struct {
	// Caption formats a hovered day's count and date.
	Caption string
}

// CalendarHeat shades day cells by count, one column per week.
func CalendarHeat(c *ui.Context, label string, from time.Time, counts []int, o CalendarHeatOptions) ui.Element {
	if label == "" || from.IsZero() || len(counts) == 0 {
		panic("datetime: CalendarHeat needs a label, a first day and counts")
	}
	top := 0
	for i, n := range counts {
		if n < 0 {
			panic(fmt.Sprintf("datetime: CalendarHeat count %d is %d", i, n))
		}
		top = max(top, n)
	}
	def(&o.Caption, "%d on %s")
	t := c.Theme()
	lead := (int(from.Weekday()) - int(c.Region().FirstWeekday) + 7) % 7
	cols := (lead + len(counts) + 6) / 7
	s, g := t.Space(2.5), t.Space(0.5)
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(1.5))
	return root.Children(func() {
		heat := ui.Box(c).Size(float32(cols)*(s+g)-g, 7*(s+g)-g).Role(ui.RoleImage).Label(label)
		k := heat.Keyframes("crychic.sweep", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		hover := -1
		if x, y, over := heat.PointerPosition(); over {
			if i := int(x/(s+g))*7 + int(y/(s+g)) - lead; i >= 0 && i < len(counts) {
				hover = i
			}
		}
		heat.Draw(func(p *ui.Painter, r ui.Rect) {
			for i, n := range counts {
				col, row := (lead+i)/7, (lead+i)%7
				a := min(1, max(0, k*float32(cols)-float32(col)))
				if a == 0 {
					continue
				}
				fill := t.Border
				if n > 0 {
					fill = t.Accent.Alpha(0.25 + 0.75*float32(n)/float32(top))
				}
				if i == hover {
					fill = t.Text
				}
				p.Fill(ui.Rect{X: r.X + float32(col)*(s+g), Y: r.Y + float32(row)*(s+g), W: s, H: s}, fill.Alpha(a), t.Space(0.5))
			}
		})
		caption := ""
		if hover >= 0 {
			caption = fmt.Sprintf(o.Caption, counts[hover],
				i18n.FormatDate(locale(c), from.AddDate(0, 0, hover), i18n.DateOptions{Skeleton: "yMMMd"}))
		}
		kit.Muted(c, caption).FontSize(t.Rem(0.85)).MinHeight(t.Space(5))
	})
}
