package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RemovableTag is a tag whose cross collapses it out.
func RemovableTag(c *ui.Context, text string, visible *bool) ui.Element {
	if visible == nil {
		panic("display: RemovableTag needs visibility")
	}
	t := c.Theme()
	spring := ui.Snappy
	root := ui.Box(c).AlignSelf(ui.Start).Clip().Transition(ui.ElementTransition{Size: true, Animation: &spring})
	return root.Children(func() {
		if !*visible {
			return
		}
		row := ui.Row(c).Key("tag").Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).
			Background(t.Surface).Border(1, t.Border).Transition(ui.ElementTransition{Exit: &ui.Motion{}})
		row.Children(func() {
			ui.Text(c, text).FontSize(t.Rem(0.85))
			close := ui.ButtonBase(c).Label("Remove " + text).TextColor(t.TextMuted)
			kit.Press(close).Children(func() { ui.Text(c, "×") })
			if kit.Clicked(close, sound.Delete) {
				*visible = false
				root.SetChanged()
				kit.Log("RemovableTag", "removed")
			}
		})
	})
}
