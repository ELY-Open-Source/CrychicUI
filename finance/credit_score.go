package finance

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CreditScoreOptions word the gauge; Note sits under the band.
type CreditScoreOptions struct{ Title, Note, Poor, Fair, Good, Excellent string }

// CreditScore shows a score on a banded, sweeping arc.
func CreditScore(c *ui.Context, score int, o CreditScoreOptions) ui.Element {
	if score < 300 || score > 850 {
		panic(fmt.Sprintf("finance: CreditScore %d must be 300–850", score))
	}
	def(&o.Title, "Credit score")
	def(&o.Poor, "Poor")
	def(&o.Fair, "Fair")
	def(&o.Good, "Good")
	def(&o.Excellent, "Excellent")
	t := c.Theme()
	bands := []struct {
		from int
		ink  ui.Color
		name string
	}{{300, t.Danger, o.Poor}, {550, t.Warning, o.Fair}, {650, t.Accent, o.Good}, {750, t.Success, o.Excellent}}
	band := bands[0]
	for _, b := range bands {
		if score >= b.from {
			band = b
		}
	}
	at := func(v int) float32 { return float32(v-300) / 550 }
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "gauge", o.Title, nil)
		ui.Column(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(1)).Padding(t.Space(4), t.Space(3), t.Space(3), t.Space(3)).Children(func() {
			arc := ui.Box(c).Size(t.Space(56), t.Space(30))
			sweep := arc.AnimateTo("score", at(score), ui.Smooth)
			ink := arc.AnimateColor("ink", band.ink, ui.Snappy)
			arc.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				stroke := t.Space(3)
				cx, cy := r.X+r.W/2, r.Y+r.H-stroke/2
				rad := min(r.W/2, r.H) - stroke
				var ring ui.Path
				ring.Circle(cx, cy, rad)
				for i, b := range bands {
					end := 850
					if i+1 < len(bands) {
						end = bands[i+1].from
					}
					p.StrokePath(ring.Trim(0.5+0.5*at(b.from)+0.004, 0.5+0.5*at(end)-0.004), stroke, b.ink.Alpha(0.18))
				}
				p.StrokePath(ring.Trim(0.5, 0.5+0.5*sweep), stroke, ink)
				angle := math.Pi * float64(1+sweep)
				x, y := cx+rad*float32(math.Cos(angle)), cy+rad*float32(math.Sin(angle))
				d := stroke + t.Space(1.5)
				p.Fill(ui.Rect{X: x - d/2, Y: y - d/2, W: d, H: d}, t.Background, 999)
				p.Fill(ui.Rect{X: x - stroke/2, Y: y - stroke/2, W: stroke, H: stroke}, ink, 999)
			}))
			arc.Children(func() {
				ui.Column(c).FillWidth().FillHeight().AlignItems(ui.Center).Justify(ui.End).Gap(t.Space(0.5)).Children(func() {
					kit.Numeric(c, fmt.Sprint(score)).Bold().FontSize(t.Rem(2.2)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
					kit.Chip(c, band.name, band.ink)
				})
			})
			ui.Row(c).Width(t.Space(56)).Children(func() {
				kit.Muted(c, "300").FontSize(t.Rem(0.7)).Grow(1)
				kit.Muted(c, "850").FontSize(t.Rem(0.7)).Shrink(0)
			})
			if o.Note != "" {
				kit.Muted(c, o.Note).FontSize(t.Rem(0.75)).Shrink(0)
			}
		})
	})
	return root
}
