package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// Logo is a mark and wordmark; the mark turns in.
func Logo(c *ui.Context, mark, wordmark string) ui.Element {
	t := c.Theme()
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		glyph := ui.Icon(c, icon.Must(mark)).Size(t.Space(5), t.Space(5)).TextColor(t.Accent)
		turn := glyph.Keyframes("turn", []ui.Keyframe{{Value: -90}, {At: 200 * time.Millisecond, Value: 0, Ease: ui.EaseOut}})
		glyph.Rotate(turn)
		ui.Text(c, wordmark).Bold().FontSize(t.Rem(1.1))
	})
}
