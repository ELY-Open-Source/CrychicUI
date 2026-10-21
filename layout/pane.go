package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Pane folds content under a title; chevron turns.
func Pane(c *ui.Context, title string, open *bool, content func()) ui.Element {
	if open == nil {
		panic("layout: Pane needs open")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Border(1, t.Border).Radius(t.Radius).Clip().
		Transition(ui.ElementTransition{Size: true, Animation: &collapseSpring})
	return root.Children(func() {
		head := ui.ButtonBase(c).FillWidth().Justify(ui.Start).Gap(t.Space(2)).Padding(t.Space(2), t.Space(3)).
			Expanded(*open)
		bg := t.Background
		if head.Hovered() {
			bg = t.SurfaceHover
		}
		head.Background(head.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
			chevron := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16)
			chevron.Rotate(chevron.AnimateTo("turn", -90*unit(!*open), ui.Snappy))
			ui.Text(c, title).Bold()
		})
		if kit.Clicked(head, sound.Toggle) {
			*open = !*open
			root.SetChanged()
			kit.Log("Pane", "toggled", "open", *open)
		}
		Collapse(c, *open, func() {
			ui.Divider(c)
			ui.Column(c).Padding(t.Space(3)).Gap(t.Space(2)).Children(content)
		})
	})
}
