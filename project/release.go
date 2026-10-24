package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReleaseOptions are the release labels.
type ReleaseOptions struct{ Ship string }

// Release shows version items and ship; cue Success.
func Release(c *ui.Context, version string, items []string, shipped *bool, o ReleaseOptions) ui.Element {
	if version == "" || len(items) == 0 || shipped == nil {
		panic("project: Release has invalid data")
	}
	def(&o.Ship, "Ship")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, version).Bold().Grow(1)
			b := kit.Press(ui.ButtonBase(c).Label(o.Ship).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).
				Background(t.Accent).TextColor(t.AccentText).Disabled(*shipped))
			b.Children(func() { ui.Text(c, o.Ship).SingleLine() })
			if b.Clicked() && !*shipped {
				*shipped = true
				root.SetChanged()
				sound.Play(sound.Success)
				kit.Log("Release", "shipped")
			}
		})
		for _, item := range items {
			if item == "" {
				panic("project: Release item must not be empty")
			}
			ui.Text(c, item).FontSize(t.Rem(0.85)).Wrap()
		}
	})
	return root
}
