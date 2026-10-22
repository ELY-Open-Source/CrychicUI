package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NotFoundOptions tune a NotFound.
type NotFoundOptions struct {
	// Code, Title, Text and Home default to English 404 copy.
	Code, Title, Text, Home string
}

// NotFound is a code, message and home action rising in.
func NotFound(c *ui.Context, o NotFoundOptions) (root, home ui.Element) {
	if o.Code == "" {
		o.Code = "404"
	}
	if o.Title == "" {
		o.Title = "Page not found"
	}
	if o.Text == "" {
		o.Text = "The page you asked for is not here."
	}
	if o.Home == "" {
		o.Home = "Go home"
	}
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Center).AlignItems(ui.Center).Gap(t.Space(2)).Padding(t.Space(6)).
		Role(ui.RoleStatus).Label(o.Title)
	kit.Arrive(root, ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	root.Children(func() {
		code := ui.Text(c, o.Code).FontSize(t.Rem(4)).Bold().TextColor(t.Accent.Alpha(0.85))
		kit.Arrive(code, ui.Pose{Scale: 0.8, Fade: 1}, 0)
		ui.Text(c, o.Title).Bold().FontSize(t.Rem(1.15))
		kit.Muted(c, o.Text).TextAlign(ui.Center)
		home = button(c, o.Home, true).Margin(t.Space(2), 0, 0, 0)
		if kit.Clicked(home, sound.Tap) {
			kit.Log("NotFound", "home")
		}
	})
	return root, home
}
