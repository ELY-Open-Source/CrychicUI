package feedback

import "github.com/egoist/mygo/ui"

// Skeleton is a placeholder block with a moving sheen.
func Skeleton(c *ui.Context) ui.Element {
	t := c.Theme()
	return sheen(c, ui.Box(c).FillWidth().Height(t.Space(4)), t.Radius)
}
