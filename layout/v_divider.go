package layout

import "github.com/egoist/mygo/ui"

// VDivider is a vertical hairline growing from center.
func VDivider(c *ui.Context) ui.Element {
	line := ui.Box(c).FillHeight().Width(1).Background(c.Theme().Border)
	return line.Scale(1, grow(line))
}
