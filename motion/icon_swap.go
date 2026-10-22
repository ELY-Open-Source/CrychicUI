package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// IconSwap shows the named Lucide icon, swapping in on change.
func IconSwap(c *ui.Context, name string) ui.Element {
	d := c.Theme().Space(6)
	return ui.Icon(c, icon.Must(name)).Size(d, d).ContentTransition(ui.ContentTransition{})
}
