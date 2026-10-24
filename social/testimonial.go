package social

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Testimonial shows a quote, name, and role.
func Testimonial(c *ui.Context, quote, name, role string) ui.Element {
	if quote == "" || name == "" || role == "" {
		panic("social: Testimonial needs quote, name, and role")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(4)).Radius(t.Radius).Background(t.Surface)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Icon(c, icon.Must("quote")).Size(t.Space(6), t.Space(6)).TextColor(t.Accent)
		ui.Text(c, quote).Wrap()
		ui.Column(c).Children(func() {
			ui.Text(c, name).Bold()
			kit.Muted(c, role)
		})
	})
	return root
}
