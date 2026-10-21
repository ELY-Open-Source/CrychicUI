package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SortDirection is the order a SortButton selects.
type SortDirection int

// Sort directions.
const (
	SortNone SortDirection = iota
	SortAscending
	SortDescending
)

// SortButton cycles none, ascending and descending.
func SortButton(c *ui.Context, label string, direction *SortDirection) ui.Element {
	if direction == nil || *direction < SortNone || *direction > SortDescending {
		panic("button: SortButton needs a valid direction")
	}
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface)
	root.Children(func() {
		ui.Text(c, label).SingleLine()
		name := "arrow-up-down"
		if *direction != SortNone {
			name = "arrow-up"
		}
		arrow := ui.Icon(c, icon.Must(name)).Size(16, 16).ContentTransition(ui.ContentTransition{})
		arrow.Rotate(arrow.AnimateTo("turn", 180*unit(*direction == SortDescending), ui.Snappy))
	})
	if root.Clicked() {
		*direction = (*direction + 1) % 3
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("SortButton", "changed", "direction", *direction)
	}
	return root
}
