package fx

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// GradientBackdrop fills behind content from accent to surface, turning slowly.
func GradientBackdrop(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	turn := 360 * phase(c, root, 20*time.Second)
	return root.LinearGradient(ui.LinearGradient{From: t.Accent.Mix(t.Surface, 0.7), To: t.Surface, Angle: 135 + turn,
		Oklab: true}).Children(content)
}
