package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ButtonGroup joins buttons and selects one with a tap cue.
func ButtonGroup(c *ui.Context, items []string, selected *int) ui.Element {
	if len(items) == 0 || selected == nil || *selected < 0 || *selected >= len(items) {
		panic("button: ButtonGroup needs items and a valid selection")
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start)
	return root.Children(func() {
		for i, item := range items {
			c.Key(i)
			b := textBase(c, item)
			switch {
			case len(items) == 1:
				b.Radius(t.Radius).Border(1, t.Border)
			case i == 0:
				b.Radius(t.Radius, 0, 0, t.Radius).Border(1, t.Border)
			case i == len(items)-1:
				b.Radius(0, t.Radius, t.Radius, 0).BorderWidth(1, 0, 1, 1).BorderColor(t.Border)
			default:
				b.Radius(0).BorderWidth(1, 0, 1, 1).BorderColor(t.Border)
			}
			bg, fg := t.Background, t.Text
			if i == *selected {
				bg, fg = t.Accent.Alpha(0.14), t.Accent
			} else if b.Pressed() {
				bg = t.SurfacePressed
			} else if b.Hovered() {
				bg = t.SurfaceHover
			}
			b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(b.AnimateColor("fg", fg, ui.Snappy)).
				Children(func() { ui.Text(c, item).SingleLine() })
			if kit.Clicked(b, sound.Tap) && i != *selected {
				*selected = i
				root.SetChanged()
				kit.Log("ButtonGroup", "selected", "index", i)
			}
		}
	})
}
