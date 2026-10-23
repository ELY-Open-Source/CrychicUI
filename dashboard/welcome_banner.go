package dashboard

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// WelcomeBanner shows a greeting and daily summary.
func WelcomeBanner(c *ui.Context, greeting, summary string) ui.Element {
	if greeting == "" || summary == "" {
		panic("dashboard: WelcomeBanner needs greeting and summary")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(4)).Radius(t.Radius).
		Background(t.Accent).TextColor(t.AccentText)
	kit.Arrive(root, ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, greeting).Bold().FontSize(t.Rem(1.4))
		ui.Text(c, summary).Wrap()
	})
	return root
}
