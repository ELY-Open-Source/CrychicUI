package learn

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Module is a course module with lessons done of total.
type Module struct {
	Title       string
	Done, Total int
}

// ProgressOptions word the card; Count takes done and total.
type ProgressOptions struct{ Title, Count, Lessons string }

// Progress shows course progress overall and by module.
func Progress(c *ui.Context, modules []Module, o ProgressOptions) ui.Element {
	def(&o.Title, "Course progress")
	if len(modules) == 0 {
		panic(fmt.Sprintf("learn: Progress %q needs modules", o.Title))
	}
	titles := make([]string, len(modules))
	done, total := 0, 0
	for i, m := range modules {
		if m.Total < 1 || m.Done < 0 || m.Done > m.Total || m.Total > math.MaxInt-total {
			panic(fmt.Sprintf("learn: Progress module %q needs 0 <= done <= total and a total that fits", m.Title))
		}
		done += m.Done
		total += m.Total
		titles[i] = m.Title
	}
	unique("Progress", titles)
	def(&o.Count, "%d of %d")
	def(&o.Lessons, "%s lessons done")
	t := c.Theme()
	current := -1
	for i, m := range modules {
		if m.Done < m.Total {
			current = i
			break
		}
	}
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(4)).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ink := t.Accent
			if done == total {
				ink = t.Success
			}
			dial(c, t.Space(20), float32(float64(done)/float64(total)), ink, func() {
				ui.Text(c, fmt.Sprintf("%d%%", percent(done, total))).FontSize(t.Rem(1)).Bold().Shrink(0)
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
				ui.Text(c, o.Title).FontSize(t.Rem(0.95)).Bold().FillWidth().MinWidth(0).SingleLine()
				kit.Muted(c, fmt.Sprintf(o.Lessons, fmt.Sprintf(o.Count, done, total))).FontSize(t.Rem(0.76)).FillWidth().MinWidth(0).Wrap()
			})
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3), t.Space(4), t.Space(4), t.Space(4))
		list.Children(func() {
			for i, m := range modules {
				full := m.Done == m.Total
				ink := t.TextMuted
				switch {
				case full:
					ink = t.Success
				case i == current:
					ink = t.Accent
				}
				c.Key(strconv.Quote(m.Title))
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center)
				kit.Stagger(c, list, row, i).Children(func() {
					badge := ui.Box(c).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(999).Background(ink.Alpha(0.14)).TextColor(ink)
					badge.Children(func() {
						c.Key(full)
						if full {
							kit.Arrive(glyph(c, "check", 14), ui.Pose{Scale: 0.4, Fade: 1}, 0)
						} else {
							ui.Text(c, fmt.Sprint(i+1)).FontSize(t.Rem(0.72)).Bold().Shrink(0)
						}
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.25)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, m.Title).FontSize(t.Rem(0.82)).Bold().Grow(1).MinWidth(0).SingleLine().
								TextColor(t.Text.Mix(t.TextMuted, unit(i > current && current >= 0)))
							kit.Muted(c, fmt.Sprintf(o.Count, m.Done, m.Total)).Font("monospace").FontSize(t.Rem(0.7)).Shrink(0).
								MaxWidth(t.Space(30)).SingleLine()
						})
						meter(c, float32(float64(m.Done)/float64(m.Total)), ink)
					})
				})
			}
		})
	})
	return root
}
