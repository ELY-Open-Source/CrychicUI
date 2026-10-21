package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RefreshButton spins its icon while refreshing.
func RefreshButton(c *ui.Context, label string, refreshing bool) ui.Element {
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface)
	root.Children(func() {
		glyph := ui.Icon(c, icon.Must("refresh-cw")).Size(16, 16)
		if refreshing {
			glyph.Rotate(glyph.Loop("spin", time.Second, ui.Linear) * 360)
		}
		ui.Text(c, label).SingleLine()
	})
	if root.Clicked() {
		sound.Play(sound.Tap)
	}
	return root
}
