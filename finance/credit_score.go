package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CreditScoreOptions are the credit band labels.
type CreditScoreOptions struct{ Poor, Fair, Good, Excellent string }

// CreditScore shows a score on an arc with band.
func CreditScore(c *ui.Context, score int, o CreditScoreOptions) ui.Element {
	if score < 300 || score > 850 {
		panic("finance: CreditScore must be 300–850")
	}
	def(&o.Poor, "Poor")
	def(&o.Fair, "Fair")
	def(&o.Good, "Good")
	def(&o.Excellent, "Excellent")
	t := c.Theme()
	band, ink := o.Poor, t.Danger
	switch {
	case score >= 750:
		band, ink = o.Excellent, t.Success
	case score >= 650:
		band, ink = o.Good, t.Accent
	case score >= 550:
		band, ink = o.Fair, t.Warning
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		arc := ui.Box(c).Size(t.Space(32), t.Space(16))
		progress := arc.AnimateTo("score", float32(score-300)/550, ui.Smooth)
		arc.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			var path ui.Path
			path.Circle(r.X+r.W/2, r.Y+r.H, r.H-t.Space(1)).Trim(0.5, 1)
			p.StrokePath(&path, t.Space(2), t.SurfaceHover)
			p.StrokePath(path.Trim(0.5, 0.5+0.5*progress), t.Space(2), ink)
		}))
		kit.Numeric(c, fmt.Sprint(score)).Bold().FontSize(t.Rem(1.5))
		ui.Text(c, band).TextColor(ink)
	})
	return root
}
