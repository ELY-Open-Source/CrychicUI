package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var wipe = replay(600*time.Millisecond, ui.EaseInOut)

// Wipe reveals text by a moving clip; run replays it.
func Wipe(c *ui.Context, run uint64, text string) ui.Element {
	if text == "" {
		panic("motion: Wipe needs text")
	}
	root := ui.Box(c).AlignSelf(ui.Start)
	k := root.Keyframes("crychic.wipe", wipe, ui.KeyframeOptions{Run: run})
	return root.Children(func() {
		ui.Text(c, text).NoWrap().Opacity(0).Role(ui.RoleNone)
		ui.Box(c).Absolute().Top(0).Bottom(0).Start(0).WidthPercent(100 * k).ClipX().Children(func() {
			ui.Text(c, text).NoWrap()
		})
	})
}
