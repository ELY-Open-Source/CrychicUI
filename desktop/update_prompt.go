package desktop

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UpdatePromptOptions word the prompt; Title formats the app name.
type UpdatePromptOptions struct{ Icon, Title, Version, Install, Later string }

// UpdatePrompt offers a new version; Install cues Success, Later Tap.
func UpdatePrompt(c *ui.Context, app, version string, notes []string, o UpdatePromptOptions) (root, install, later ui.Element) {
	def(&o.Icon, "download")
	def(&o.Title, "A new version of %s is ready")
	def(&o.Version, "Version %s")
	def(&o.Install, "Install update")
	def(&o.Later, "Later")
	if strings.TrimSpace(app) == "" || strings.TrimSpace(version) == "" {
		panic(fmt.Sprintf("desktop: UpdatePrompt %q needs an app and a version", app))
	}
	t := c.Theme()
	root = card(c).Gap(t.Space(3)).Padding(t.Space(5)).Role(ui.RoleAlertDialog).Label(fmt.Sprintf(o.Title, app))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			kit.Arrive(tile(c, o.Icon, t.Accent, t.Space(11)).Radius(t.Radius*2), ui.Pose{Y: -t.Space(1), Scale: 0.8, Fade: 1}, 0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, fmt.Sprintf(o.Title, app)).FontSize(t.Rem(0.92)).Bold().FillWidth().Wrap()
				ui.Text(c, fmt.Sprintf(o.Version, version)).Font("monospace").FontSize(t.Rem(0.72)).TextColor(t.Accent).FillWidth().MinWidth(0).SingleLine()
			})
		})
		if len(notes) > 0 {
			list := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(3)).Radius(t.Radius).Background(t.SurfaceHover)
			list.Children(func() {
				for i, n := range notes {
					c.Key(i)
					row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start)
					kit.Stagger(c, list, row, i).Children(func() {
						glyph(c, "check", 14).TextColor(t.Success)
						ui.Text(c, n).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
					})
				}
			})
		}
		c.Key("actions")
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Justify(ui.End).Children(func() {
			later = kit.Press(ui.ButtonBase(c).Label(o.Later).Padding(t.Space(2), t.Space(4)).Radius(999).Border(1, t.Border))
			later.Background(later.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(later.Hovered())), ui.Snappy))
			later.Children(func() { ui.Text(c, o.Later).FontSize(t.Rem(0.8)).Bold().Shrink(0).MaxWidth(t.Space(36)).SingleLine() })
			install = kit.Press(ui.ButtonBase(c).Label(o.Install).Padding(t.Space(2), t.Space(4)).Radius(999))
			install.Background(install.AnimateColor("bg", t.Accent.Mix(t.Text, 0.12*unit(install.Hovered())), ui.Snappy)).TextColor(t.AccentText)
			install.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, "download", 14)
					ui.Text(c, o.Install).FontSize(t.Rem(0.8)).Bold().Shrink(0).MaxWidth(t.Space(44)).SingleLine()
				})
			})
		})
		if later.Clicked() {
			sound.Play(sound.Tap)
			kit.Log("UpdatePrompt", "later")
		}
		if install.Clicked() {
			sound.Play(sound.Success)
			kit.Log("UpdatePrompt", "install")
		}
	})
	return root, install, later
}
