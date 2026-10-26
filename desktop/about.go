package desktop

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AboutOptions describe the app; Format words the version.
type AboutOptions struct{ Icon, Version, Format, Copyright, Credits string }

// About is an app's about window.
func About(c *ui.Context, app string, o AboutOptions) ui.Element {
	def(&o.Icon, "app-window")
	def(&o.Format, "Version %s")
	if strings.TrimSpace(app) == "" {
		panic("desktop: About needs an app name")
	}
	t := c.Theme()
	root := window(c).Role(ui.RoleGroup).Label(app)
	return root.Children(func() {
		titleBar(c, "", true, nil)
		ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(6), t.Space(5)).Children(func() {
			icon := tile(c, o.Icon, t.Accent, t.Space(18)).Radius(t.Radius*3).Margin(0, 0, t.Space(2), 0)
			kit.Arrive(icon, ui.Pose{Scale: 0.6, Rotate: -12, Fade: 1}, 0)
			ui.Text(c, app).FontSize(t.Rem(1.2)).Bold().FillWidth().TextAlign(ui.Center).Wrap()
			if o.Version != "" {
				ui.Text(c, fmt.Sprintf(o.Format, o.Version)).Font("monospace").FontSize(t.Rem(0.76)).TextColor(t.TextMuted).
					FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
			}
			if o.Credits != "" {
				kit.Muted(c, o.Credits).FontSize(t.Rem(0.78)).FillWidth().MaxWidth(t.Space(80)).TextAlign(ui.Center).Wrap().Margin(t.Space(2), 0, 0, 0)
			}
			if o.Copyright != "" {
				ui.Text(c, o.Copyright).FontSize(t.Rem(0.68)).TextColor(t.TextMuted.Alpha(0.8)).FillWidth().TextAlign(ui.Center).Wrap()
			}
		})
	})
}
