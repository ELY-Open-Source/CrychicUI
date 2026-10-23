package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HealthRingOptions are the health labels.
type HealthRingOptions struct{ Label string }

// HealthRing shows overall health as a ring.
func HealthRing(c *ui.Context, score int, o HealthRingOptions) ui.Element {
	if score < 0 || score > 100 {
		panic("dashboard: HealthRing score must be 0–100")
	}
	def(&o.Label, "Health")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ring := ui.Box(c).Size(t.Space(16), t.Space(16))
		progress := ring.AnimateTo("health", float32(score)/100, ui.Smooth)
		ring.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			var path ui.Path
			path.Circle(r.X+r.W/2, r.Y+r.H/2, r.W/2-t.Space(1))
			p.StrokePath(&path, t.Space(2), t.SurfaceHover)
			p.StrokePath(path.Trim(0, progress), t.Space(2), t.Success)
		}))
		ui.Text(c, o.Label).Bold()
		kit.Numeric(c, fmt.Sprintf("%d%%", score)).FontSize(t.Rem(1.2))
	})
	return root
}
