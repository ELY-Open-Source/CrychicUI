package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LoadingButton shows a spinner and disables while busy.
func LoadingButton(c *ui.Context, label string, busy bool) ui.Element {
	t := c.Theme()
	spring := ui.Snappy
	b := textBase(c, label).Disabled(busy).Border(1, t.Border).Background(t.Surface).
		Transition(ui.ElementTransition{Size: true, Animation: &spring})
	b.Children(func() {
		if busy {
			spin := ui.Icon(c, icon.Must("loader-circle")).Size(16, 16)
			spin.Rotate(spin.Loop("spin", time.Second, ui.Linear) * 360)
		}
		ui.Text(c, label).SingleLine()
	})
	kit.Clicked(b, sound.Tap)
	return b
}
