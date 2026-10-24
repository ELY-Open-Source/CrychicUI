package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MindMapOptions are the mind map labels.
type MindMapOptions struct{ Add string }

// MindMap shows nodes branching from a center; cue Pop.
func MindMap(c *ui.Context, center string, branches []string, added *int, o MindMapOptions) ui.Element {
	if center == "" || added == nil || *added < -1 || *added > len(branches) {
		panic("project: MindMap has invalid data")
	}
	def(&o.Add, "Add branch")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, center).Padding(t.Space(2), t.Space(3)).Radius(999).Background(t.Accent).TextColor(t.AccentText).Bold()
		ui.Row(c).Wrap().Gap(t.Space(2)).Justify(ui.Center).Children(func() {
			for i, branch := range branches {
				if branch == "" {
					panic("project: MindMap branch must not be empty")
				}
				c.Key(i)
				ui.Text(c, branch).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Background(t.Accent.Alpha(0.15))
			}
		})
		b := kit.Press(ui.ButtonBase(c).Label(o.Add).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
		b.Children(func() { ui.Text(c, o.Add).SingleLine() })
		if b.Clicked() {
			*added = len(branches)
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("MindMap", "added", "count", len(branches))
		}
	})
	return root
}
