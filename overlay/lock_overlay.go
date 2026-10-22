package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LockOverlayOptions tune a LockOverlay.
type LockOverlayOptions struct {
	Text, Unlock string
}

// LockOverlay blurs content under a locked panel until unlocked.
func LockOverlay(c *ui.Context, title string, locked *bool, o LockOverlayOptions, content func()) ui.Element {
	if title == "" || locked == nil || content == nil {
		panic(fmt.Sprintf("overlay: LockOverlay %q needs locked and content", title))
	}
	def(&o.Unlock, "Unlock")
	t := c.Theme()
	root := ui.Box(c)
	cover(c, root, *locked, t.Space(1), content, func() {
		card(c, title, o.Text, "lock", func() {
			if primary(c, o.Unlock, t.Accent, t.AccentText).Clicked() {
				*locked = false
				root.SetChanged()
				sound.Play(sound.Success)
			}
		})
	})
	moved(root, *locked, "LockOverlay")
	return root
}
