package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Resource is one room and its events.
type Resource struct {
	Name   string
	Events []Event
}

// ResourceRowsOptions set the title and hour window (default 8-20).
type ResourceRowsOptions struct {
	Title    string
	From, To int
}

// ResourceRows shows rooms as rows on a time axis.
func ResourceRows(c *ui.Context, resources []Resource, o ResourceRowsOptions) ui.Element {
	if len(resources) == 0 {
		panic("schedule: ResourceRows needs resources")
	}
	if o.From == 0 && o.To == 0 {
		o.From, o.To = 8, 20
	}
	if o.From < 0 || o.To > 24 || o.From >= o.To {
		panic(fmt.Sprintf("schedule: ResourceRows window %d-%d is invalid", o.From, o.To))
	}
	def(&o.Title, "Rooms")
	t := c.Theme()
	n := float32(o.To - o.From)
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "door-open", o.Title, nil)
		list := ui.Column(c).FillWidth().Padding(t.Space(2), t.Space(3), t.Space(3), t.Space(3)).Gap(t.Space(1.5))
		list.Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Box(c).Width(t.Space(22)).Shrink(0)
				ui.Row(c).Grow(1).MinWidth(0).Children(func() {
					for h := o.From; h < o.To; h++ {
						ui.Box(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
							kit.Muted(c, fmt.Sprintf("%02d", h)).FontSize(t.Rem(0.6))
						})
					}
				})
			})
			for i, resource := range resources {
				if resource.Name == "" {
					panic("schedule: ResourceRows room must not be empty")
				}
				validEvents("ResourceRows", resource.Events)
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Row(c).Width(t.Space(22)).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
						glyph(c, "door-closed", 14)
						ui.Text(c, resource.Name).FontSize(t.Rem(0.8)).Bold().TextColor(t.Text).SingleLine().Grow(1).MinWidth(0)
					})
					track := ui.Box(c).Grow(1).MinWidth(0).Height(t.Space(10)).Radius(t.Space(1)).Clip().Background(t.SurfaceHover.Alpha(0.5))
					w := track.LayoutBounds().W
					kit.Settle(c, track, []float32{w}, func() []float32 { return []float32{track.LayoutBounds().W} })
					track.Draw(func(p *ui.Painter, r ui.Rect) {
						for h := 1; h < o.To-o.From; h++ {
							x := r.X + r.W*float32(h)/n
							p.Line(x, r.Y, x, r.Y+r.H, 1, t.Border.Alpha(0.7))
						}
					})
					track.Children(func() {
						for _, e := range resource.Events {
							day := time.Date(e.Start.Year(), e.Start.Month(), e.Start.Day(), 0, 0, 0, 0, e.Start.Location())
							lo := max(0, min(1, (hours(e.Start, day)-float32(o.From))/n))
							hi := max(0, min(1, (hours(e.End, day)-float32(o.From))/n))
							if hi <= lo {
								continue
							}
							c.Key(e.ID)
							b := ui.Column(c).Absolute().Left(w*lo + 1).Width(max(w*(hi-lo)-2, 0)).Top(t.Space(1)).Bottom(t.Space(1))
							eventBox(c, e, e.Start.Location(), b, false)
						}
					})
				})
			}
		})
	})
	return root
}
