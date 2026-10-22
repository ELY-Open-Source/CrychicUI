package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CollapsibleNavOptions tune a CollapsibleNav.
type CollapsibleNavOptions struct {
	// Collapse and Expand name the toggle button.
	Collapse, Expand string
}

// CollapsibleNav folds a side nav to icons.
func CollapsibleNav(c *ui.Context, label string, list []NavItem, active *int, open *bool, o CollapsibleNavOptions) ui.Element {
	items("CollapsibleNav", label, list, active)
	if open == nil {
		panic("navigation: CollapsibleNav needs an open state")
	}
	def(&o.Collapse, "Collapse")
	def(&o.Expand, "Expand")
	t := c.Theme()
	narrow, wide := t.Space(16), t.Space(56)
	root := ui.Column(c).Gap(t.Space(0.5)).Padding(t.Space(2)).Clip().FocusGroup(ui.Vertical).Label(label)
	k := root.AnimateTo("open", unit(*open), ui.Bouncy)
	id := root.ID()
	root.Width(narrow + (wide-narrow)*k).Children(func() {
		name, glyph := o.Collapse, "panel-left-close"
		if !*open {
			name, glyph = o.Expand, "panel-left-open"
		}
		b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(10), t.Space(10)).Center().Radius(t.Radius)).
			Expanded(*open)
		b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
		})
		if kit.Clicked(b, sound.Toggle) {
			*open = !*open
			root.SetChanged()
			kit.Log("CollapsibleNav", "toggled", "open", *open)
		}
		if navList(c, "CollapsibleNav", id, list, active, 1-k) {
			root.SetChanged()
		}
	})
	return root
}
