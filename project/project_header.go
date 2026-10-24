package project

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProjectHeaderOptions are the project header labels.
type ProjectHeaderOptions struct{ Icon string }

// ProjectHeader shows icon, name, and tabs; cue Select.
func ProjectHeader(c *ui.Context, name string, tab *string, tabs []string, o ProjectHeaderOptions) ui.Element {
	if name == "" || tab == nil || len(tabs) == 0 {
		panic("project: ProjectHeader has invalid data")
	}
	def(&o.Icon, "folder")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must(o.Icon)).Size(t.Space(6), t.Space(6)).TextColor(t.Accent)
			ui.Text(c, name).Bold()
		})
		ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
			for _, item := range tabs {
				if item == "" {
					panic("project: ProjectHeader tab must not be empty")
				}
				b := kit.Press(ui.ButtonBase(c).Label(item).Padding(t.Space(1), t.Space(2)).Radius(999).
					Border(1, t.Border).Checked(*tab == item))
				b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*tab == item)), ui.Snappy))
				b.Children(func() { ui.Text(c, item).FontSize(t.Rem(0.75)) })
				if b.Clicked() && *tab != item {
					*tab = item
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("ProjectHeader", "tab", "tab", *tab)
				}
			}
		})
	})
	return root
}
