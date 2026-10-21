package display

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// IconViewOptions tune an IconView.
type IconViewOptions struct {
	// Size is the icon's edge in DIPs; zero is 16.
	Size float32
}

// IconView shows a Lucide icon with swapping content.
func IconView(c *ui.Context, name string, o IconViewOptions) ui.Element {
	if o.Size == 0 {
		o.Size = c.Theme().Space(4)
	}
	if !finite(float64(o.Size)) || o.Size <= 0 {
		panic("display: IconView size must be finite and positive")
	}
	return ui.Icon(c, icon.Must(name)).Size(o.Size, o.Size).ContentTransition(ui.ContentTransition{})
}
