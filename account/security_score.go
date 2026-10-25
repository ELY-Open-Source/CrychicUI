package account

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SecurityStep is one account security recommendation.
type SecurityStep struct {
	Title string
	Done  bool
}

// SecurityScoreOptions word the card; Done takes done and total.
type SecurityScoreOptions struct{ Label, Out, Strong, Fair, Weak, Done string }

// SecurityScore shows a score ring with steps.
func SecurityScore(c *ui.Context, score int, steps []SecurityStep, o SecurityScoreOptions) ui.Element {
	if score < 0 || score > 100 || len(steps) == 0 {
		panic(fmt.Sprintf("account: SecurityScore %q needs a 0-100 score and steps", o.Label))
	}
	def(&o.Label, "Security score")
	def(&o.Out, "/ 100")
	def(&o.Strong, "Strong")
	def(&o.Fair, "Fair")
	def(&o.Weak, "At risk")
	def(&o.Done, "%d of %d done")
	t := c.Theme()
	band, ink := o.Weak, t.Danger
	switch {
	case score >= 80:
		band, ink = o.Strong, t.Success
	case score >= 50:
		band, ink = o.Fair, t.Warning
	}
	done := 0
	for _, s := range steps {
		if s.Title == "" {
			panic(fmt.Sprintf("account: SecurityScore %q step needs a title", o.Label))
		}
		if s.Done {
			done++
		}
	}
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(4)).Children(func() {
			dial := ui.Box(c).Size(t.Space(22), t.Space(22)).Shrink(0).Center()
			grow := dial.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 800 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			fill := grow * dial.AnimateTo("score", float32(score)/100, ui.Smooth)
			col := dial.AnimateColor("ink", ink, ui.Snappy)
			dial.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-t.Space(1)
				p.StrokePath(arc(cx, cy, rad, 1), t.Space(1.75), t.SurfaceHover)
				if fill > 0 {
					p.StrokePath(arc(cx, cy, rad, fill), t.Space(1.75), col)
				}
			}))
			dial.Children(func() {
				ui.Column(c).AlignItems(ui.Center).Children(func() {
					kit.Numeric(c, fmt.Sprint(score)).FontSize(t.Rem(1.5)).Bold().Shrink(0)
					kit.Muted(c, o.Out).FontSize(t.Rem(0.7)).Shrink(0)
				})
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
				ui.Text(c, o.Label).FontSize(t.Rem(0.95)).Bold().Wrap().FillWidth()
				c.Key(band)
				kit.Arrive(kit.Chip(c, band, ink).AlignSelf(ui.Start), ui.Pose{Scale: 0.7, Fade: 1}, 0)
				kit.Muted(c, fmt.Sprintf(o.Done, done, len(steps))).FontSize(t.Rem(0.78))
			})
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(0, t.Space(2), t.Space(2), t.Space(2)).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
		list.Children(func() {
			for i, step := range steps {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				kit.Stagger(c, list, row, i).Children(func() {
					c.Key(step.Done)
					mark, tint := "circle", t.TextMuted
					if step.Done {
						mark, tint = "circle-check", t.Success
					}
					kit.Arrive(glyph(c, mark, 18).TextColor(tint), ui.Pose{Scale: 0.5, Fade: 1}, 0)
					title := ui.Text(c, step.Title).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
					if step.Done {
						title.TextColor(t.TextMuted)
					}
				})
			}
		})
	})
	return root
}

// arc runs clockwise from the top for share.
func arc(cx, cy, r, share float32) *ui.Path {
	var p ui.Path
	n := 64
	for i := 0; i <= n; i++ {
		a := float64(share)*2*math.Pi*float64(i)/float64(n) - math.Pi/2
		x, y := cx+r*float32(math.Cos(a)), cy+r*float32(math.Sin(a))
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	return &p
}
