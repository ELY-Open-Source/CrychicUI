package account

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Resource is one metered account resource.
type Resource struct {
	Name, Unit  string
	Used, Limit float64
}

// UsageOptions are the usage value formats.
type UsageOptions struct{ Format string }

// Usage shows usage meters by resource.
func Usage(c *ui.Context, resources []Resource, o UsageOptions) ui.Element {
	if len(resources) == 0 {
		panic("account: Usage needs resources")
	}
	def(&o.Format, "%s/%s %s")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, resource := range resources {
			if resource.Name == "" || resource.Unit == "" || resource.Used < 0 || resource.Limit <= 0 ||
				!finite(resource.Used) || !finite(resource.Limit) {
				panic("account: invalid resource")
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).Children(func() {
					ui.Text(c, resource.Name).Grow(1)
					used := strconv.FormatFloat(resource.Used, 'g', -1, 64)
					limit := strconv.FormatFloat(resource.Limit, 'g', -1, 64)
					kit.Numeric(c, fmt.Sprintf(o.Format, used, limit, resource.Unit)).TextColor(t.TextMuted)
				})
				bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
				w := bar.AnimateTo("usage", float32(min(1, resource.Used/resource.Limit)), ui.Smooth)
				bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, t.Accent, 999)
				}))
			})
		}
	})
	return root
}
