package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StarButtonOptions tune a StarButton.
type StarButtonOptions struct {
	// Label names the action.
	Label string
}

// StarButton fills and spins a star when toggled.
func StarButton(c *ui.Context, on *bool, o StarButtonOptions) ui.Element {
	if on == nil {
		panic("button: StarButton needs state")
	}
	if o.Label == "" {
		o.Label = "Star"
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(o.Label).Tooltip(o.Label).Size(t.Space(10), t.Space(10)).Center().
		Radius(t.Radius).Border(1, t.Border).Background(t.Surface).Checked(*on))
	root.Children(func() {
		glyph := icon.Must("star")
		if *on {
			glyph = icon.Filled("star")
		}
		star := ui.Icon(c, glyph).Size(16, 16)
		if *on {
			star.TextColor(t.Warning)
		}
		star.Rotate(star.AnimateTo("turn", 180*unit(*on), ui.Snappy)).ContentTransition(ui.ContentTransition{})
	})
	if root.Clicked() {
		*on = !*on
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("StarButton", "toggled", "on", *on)
	}
	return root
}
