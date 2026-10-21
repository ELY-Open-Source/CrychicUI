package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Breakpoint shows compact under 40 rem, regular under 64.
func Breakpoint(c *ui.Context, compact, regular, wide func()) ui.Element {
	t := c.Theme()
	w, _ := c.Size()
	key, fill := "wide", wide
	switch {
	case w < t.Rem(40):
		key, fill = "compact", compact
	case w < t.Rem(64):
		key, fill = "regular", regular
	}
	return ui.Box(c).FillWidth().Children(func() {
		slot := ui.Column(c).Key(key).FillWidth().Transition(ui.ElementTransition{Exit: &ui.Motion{}})
		kit.Arrive(slot, ui.Pose{Fade: 1}, 0).Children(fill)
	})
}
