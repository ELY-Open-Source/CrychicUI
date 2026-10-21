package display

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// StatusBadge is a pill with a status dot and text.
func StatusBadge(c *ui.Context, text string, live bool) ui.Element {
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1.5)).Padding(t.Space(1), t.Space(2)).Radius(999).
		Background(t.Surface).Border(1, t.Border)
	return root.Children(func() {
		dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999)
		if live {
			pulse := dot.Loop("pulse", 1200*time.Millisecond, ui.Linear)
			alpha := 0.35 + 0.65*pulse
			dot.Background(t.Success.Alpha(alpha))
		} else {
			dot.Background(t.TextMuted)
		}
		ui.Text(c, text).FontSize(t.Rem(0.85))
	})
}
