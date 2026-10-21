package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThemeToggleOptions tune a ThemeToggle.
type ThemeToggleOptions struct {
	// Label names the action.
	Label string
}

// ThemeToggle swaps sun and moon with a turn.
func ThemeToggle(c *ui.Context, dark *bool, o ThemeToggleOptions) ui.Element {
	if dark == nil {
		panic("button: ThemeToggle needs state")
	}
	if o.Label == "" {
		o.Label = "Toggle theme"
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(o.Label).Tooltip(o.Label).AlignSelf(ui.Start).Size(t.Space(10), t.Space(10)).Center().
		Radius(999).Background(t.Surface).Border(1, t.Border).Checked(*dark))
	root.Children(func() {
		name := "moon"
		if *dark {
			name = "sun"
		}
		glyph := ui.Icon(c, icon.Must(name)).Size(16, 16).ContentTransition(ui.ContentTransition{})
		glyph.Rotate(glyph.AnimateTo("turn", 180*unit(*dark), ui.Snappy))
	})
	if root.Clicked() {
		*dark = !*dark
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("ThemeToggle", "toggled", "dark", *dark)
	}
	return root
}
