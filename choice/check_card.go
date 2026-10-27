package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CheckCard is a selectable card with animated border and check.
func CheckCard(c *ui.Context, title, body string, checked *bool) ui.Element {
	if checked == nil || title == "" {
		panic("choice: " + who("CheckCard", title) + " needs a title and state")
	}
	t := c.Theme()
	root := ui.CheckboxBase(c, checked).Label(title).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Radius(t.Radius)
	on := *checked
	edge, fill := t.Border, t.Background
	if on {
		edge, fill = t.Accent, t.Accent.Alpha(0.06)
	} else if root.Hovered() {
		edge = t.Border.Mix(t.Text, 0.25)
	}
	root.Background(root.AnimateColor("fill", fill, ui.Snappy)).Border(1, root.AnimateColor("edge", edge, ui.Snappy))
	root.AlignItems(ui.Start).Children(func() {
		ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, title).Bold()
			if body != "" {
				kit.Muted(c, body)
			}
		})
		box(c, root, on, false)
	})
	if kit.Changed(root, sound.Toggle) {
		kit.Log("CheckCard", "toggled", "on", *checked)
	}
	return root
}
