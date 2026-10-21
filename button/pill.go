package button

import "github.com/egoist/mygo/ui"

// Pill is a fully rounded solid button.
func Pill(c *ui.Context, label string) ui.Element {
	return Solid(c, label).Radius(999)
}
