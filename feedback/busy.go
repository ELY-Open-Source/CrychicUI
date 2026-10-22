package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// Busy is a spinner and text fading to done.
func Busy(c *ui.Context, busy bool, working, done string) ui.Element {
	if working == "" || done == "" {
		panic("feedback: Busy needs working and done text")
	}
	t := c.Theme()
	text := done
	if busy {
		text = working
	}
	root := ui.Row(c).Gap(t.Space(2)).AlignSelf(ui.Start).Role(ui.RoleStatus).Label(text)
	return root.Children(func() {
		ui.Box(c).Size(t.Space(4), t.Space(4)).Center().Children(func() {
			if busy {
				c.Key("spin")
				Spin(c, working, SpinOptions{Size: SpinSmall}).Transition(fade)
				return
			}
			c.Key("done")
			mark := ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Success)
			mark.Transition(fade)
		})
		ui.Text(c, text).TextColor(t.TextMuted).ContentTransition(ui.ContentTransition{})
	})
}
