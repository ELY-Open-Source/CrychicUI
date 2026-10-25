package health

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BMIOptions word the scale; Bands run under to obese.
type BMIOptions struct {
	Title string
	Bands [4]string
}

// BMI places body mass index on a 15–40 band scale.
func BMI(c *ui.Context, kg, meters float64, o BMIOptions) ui.Element {
	def(&o.Title, "BMI")
	if !(kg >= 2 && kg <= 500) || !(meters >= 0.4 && meters <= 2.6) {
		panic(fmt.Sprintf("health: BMI %q needs 2 to 500 kg and 0.4 to 2.6 m, got %v and %v", o.Title, kg, meters))
	}
	for i, b := range [...]string{"Underweight", "Healthy", "Overweight", "Obese"} {
		def(&o.Bands[i], b)
	}
	t := c.Theme()
	bmi := kg / (meters * meters)
	cuts := [5]float64{15, 18.5, 25, 30, 40}
	band := 3
	for i := 1; i < 4; i++ {
		if bmi < cuts[i] {
			band = i - 1
			break
		}
	}
	inks := [4]ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	pos := float32((math.Max(15, math.Min(40, bmi)) - 15) / 25)
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleMeter).Label(o.Title)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				kit.Numeric(c, fmt.Sprintf("%.1f", math.Floor(bmi*10)/10)).FontSize(t.Rem(1.6)).Bold().Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
			})
			c.Key(band)
			kit.Arrive(kit.Chip(c, o.Bands[band], inks[band]), ui.Pose{Scale: 0.7, Fade: 1}, 0)
		})
		scale := ui.Box(c).FillWidth().Height(t.Space(7))
		x := scale.AnimateTo("at", pos, ui.Bouncy)
		scale.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			bar := r.H * 0.4
			y := r.Y + r.H - bar
			for i := range 4 {
				a := float32((cuts[i] - 15) / 25)
				b := float32((cuts[i+1] - 15) / 25)
				p.Fill(ui.Rect{X: r.X + r.W*a + 1, Y: y, W: r.W*(b-a) - 2, H: bar}, inks[i].Alpha(0.75), bar/2)
			}
			mx := r.X + r.W*max(0, min(1, x))
			var tip ui.Path
			tip.MoveTo(mx-t.Space(1.5), r.Y).LineTo(mx+t.Space(1.5), r.Y).LineTo(mx, y-t.Space(0.5))
			p.FillPath(tip.Close(), t.Text)
		}))
		ui.Row(c).FillWidth().Children(func() {
			for i := range 4 {
				share := float32((cuts[i+1] - cuts[i]) / 25)
				kit.Muted(c, o.Bands[i]).FontSize(t.Rem(0.6)).Grow(share).Basis(0).MinWidth(0).TextAlign(ui.Center).SingleLine()
			}
		})
	})
	return root
}
