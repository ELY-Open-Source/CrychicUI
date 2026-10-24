package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RiskMeterOptions word the scale; Levels name each step.
type RiskMeterOptions struct {
	Title, Low, High, Of string
	Levels               [5]string
}

// RiskMeter marks risk one through five on a colored scale.
func RiskMeter(c *ui.Context, level int, o RiskMeterOptions) ui.Element {
	if level < 1 || level > 5 {
		panic(fmt.Sprintf("finance: RiskMeter level %d must be 1–5", level))
	}
	def(&o.Title, "Risk level")
	def(&o.Low, "Low risk")
	def(&o.High, "High risk")
	def(&o.Of, "%d of 5")
	for i, name := range [...]string{"Very low", "Low", "Moderate", "High", "Very high"} {
		def(&o.Levels[i], name)
	}
	t := c.Theme()
	inks := [...]ui.Color{t.Success, t.Success.Mix(t.Warning, 0.5), t.Warning, t.Warning.Mix(t.Danger, 0.5), t.Danger}
	ink := inks[level-1]
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "shield", ink)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.75))
				ui.Text(c, o.Levels[level-1]).Bold().FontSize(t.Rem(1.1)).FillWidth().SingleLine().ContentTransition(ui.ContentTransition{})
			})
			kit.Chip(c, fmt.Sprintf(o.Of, level), ink)
		})
		scale := ui.Box(c).FillWidth().Height(t.Space(6))
		at := scale.AnimateTo("level", float32(level-1), ui.Bouncy)
		lit := make([]ui.Color, 5)
		for i := range lit {
			lit[i] = scale.AnimateColor(fmt.Sprint("seg", i), t.SurfaceHover.Mix(inks[i], unit(i < level)), ui.Snappy)
		}
		marker := scale.AnimateColor("marker", ink, ui.Snappy)
		scale.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			gap := t.Space(1)
			w := (r.W - 4*gap) / 5
			y := r.Y + r.H - t.Space(2)
			for i, col := range lit {
				p.Fill(ui.Rect{X: r.X + float32(i)*(w+gap), Y: y, W: w, H: t.Space(2)}, col, 999)
			}
			cx := r.X + at*(w+gap) + w/2
			var pin ui.Path
			pin.MoveTo(cx-t.Space(1.25), r.Y).LineTo(cx+t.Space(1.25), r.Y).LineTo(cx, r.Y+t.Space(2)).Close()
			p.FillPath(&pin, marker)
		}))
		ui.Row(c).FillWidth().Children(func() {
			kit.Muted(c, o.Low).FontSize(t.Rem(0.75)).Grow(1)
			kit.Muted(c, o.High).FontSize(t.Rem(0.75)).Shrink(0)
		})
	})
	return root
}
