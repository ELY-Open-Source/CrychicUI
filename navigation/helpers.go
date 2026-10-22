// Package navigation holds bars, tabs, menus and paths.
package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// mark matches one component's sliding indicator across items.
type mark struct{ id uint64 }

// glide flies the indicator of component id between items.
func glide(id uint64) ui.ElementTransition {
	return ui.ElementTransition{Animation: &ui.Snappy, Match: mark{id}}
}

// stack puts the indicator under face, one cell.
func stack(c *ui.Context, on bool, indicator func() ui.Element, face func() ui.Element) {
	ui.Grid(c).Columns(1).FillWidth().Children(func() {
		if on {
			c.Key("mark")
			indicator().ColumnStart(1).RowStart(1)
		}
		c.Key("face")
		face().ColumnStart(1).RowStart(1)
	})
}

// rowBase is a full-width button tinted while hovered.
func rowBase(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).FillWidth().Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius))
	bg := t.SurfaceHover.Alpha(0)
	if b.Hovered() {
		bg = t.SurfaceHover
	}
	return b.Background(b.AnimateColor("bg", bg, ui.Snappy))
}

// row is a rowBase showing label.
func row(c *ui.Context, label string) ui.Element {
	return rowBase(c, label).Children(func() { ui.Text(c, label).SingleLine().FillWidth() })
}

// valid reports whether i indexes n items.
func valid(i, n int) bool { return i >= 0 && i < n }

// def sets an empty copy field to its English default.
func def(field *string, english string) {
	if *field == "" {
		*field = english
	}
}

// pip is a page dot that widens while on.
func pip(c *ui.Context, e ui.Element, on bool) ui.Element {
	t := c.Theme()
	k := e.AnimateTo("on", unit(on), ui.Snappy)
	return e.Size(t.Space(2)+k*t.Space(4), t.Space(2)).Radius(999).Background(t.Border.Mix(t.Accent, k))
}

func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}

// arrowButton is a small icon button playing cue.
func arrowButton(c *ui.Context, name, glyph string, off bool, cue sound.Cue) bool {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius)).Disabled(off)
	b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(16, 16)
	})
	return kit.Clicked(b, cue)
}
