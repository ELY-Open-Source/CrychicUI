package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Testimonial shows a quote, name, and role.
func Testimonial(c *ui.Context, quote, name, role string) ui.Element {
	if quote == "" || name == "" || role == "" {
		panic(fmt.Sprintf("social: Testimonial %q needs quote, name and role", name))
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4))
	root.Children(func() {
		mark := ui.Box(c).Size(t.Space(10), t.Space(10)).Center().Radius(999).Background(t.Accent.Alpha(0.12)).TextColor(t.Accent)
		kit.Arrive(mark, ui.Pose{Scale: 0.5, Rotate: -20, Fade: 1}, 0).Children(func() { glyph(c, "quote", 20) })
		ui.Text(c, quote).FontSize(t.Rem(1.02)).Wrap().FillWidth()
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(3), 0, 0, 0).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			avatar(c, User{ID: name, Name: name}, t.Space(10))
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(0.88)).Bold().FillWidth().MinWidth(0).SingleLine()
				kit.Muted(c, role).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
			})
		})
	})
	return root
}
