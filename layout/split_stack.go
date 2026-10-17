package layout

import "github.com/egoist/mygo/ui"

// SplitStack holds two panes stacked; *size is the first's height.
func SplitStack(c *ui.Context, size *float32, first, second func()) ui.Element {
	return split(c, "SplitStack", true, size, first, second)
}
