package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ScrollTabsOptions tune ScrollTabs.
type ScrollTabsOptions struct {
	// Back and Forth name the arrow buttons.
	Back, Forth string
}

// ScrollTabs are overflowing tabs with arrow buttons.
func ScrollTabs(c *ui.Context, label string, names []string, active *int, o ScrollTabsOptions) (root, back, forth ui.Element) {
	tabs("ScrollTabs", label, names, active)
	def(&o.Back, "Scroll back")
	def(&o.Forth, "Scroll forward")
	t := c.Theme()
	root = ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(1))
	st := ui.Local(root, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	last := ui.Local(root, "last", func() int { return *active })
	used := []float32{st.X, st.MaxX}
	step, shown := t.Space(40), *active
	arrow := func(name, glyph string, off bool) ui.Element {
		b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius)).Disabled(off)
		return b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must(glyph)).Size(16, 16)
		})
	}
	root.Children(func() {
		back = arrow(o.Back, "chevron-left", st.X <= 0)
		if kit.Clicked(back, sound.Tap) {
			st.ScrollTo(max(0, st.X-step), 0, ui.Smooth)
			kit.Log("ScrollTabs", "scrolled", "dir", -1)
		}
		ui.ScrollHorizontal(c).Grow(1).TrackScroll(st).Children(func() {
			parts := ui.TabsBase(c, active, len(names))
			list := parts.List.Label(label).Gap(t.Space(1))
			id := list.ID()
			all := make([]ui.Element, len(names))
			list.Children(func() {
				for i, name := range names {
					c.Key(i)
					all[i] = underlined(c, parts.Tab(i), name, i == shown, id, nil)
				}
			})
			if kit.Changed(list, sound.Select) {
				root.SetChanged()
				kit.Log("ScrollTabs", "picked", "index", *active)
			}
			// Keys may pick a tab built before them.
			if *active != *last {
				all[*active].ScrollIntoView()
				*last = *active
			}
		})
		forth = arrow(o.Forth, "chevron-right", st.X >= st.MaxX)
		if kit.Clicked(forth, sound.Tap) {
			st.ScrollTo(min(st.MaxX, st.X+step), 0, ui.Smooth)
			kit.Log("ScrollTabs", "scrolled", "dir", 1)
		}
	})
	kit.Settle(c, root, used, func() []float32 { return []float32{st.X, st.MaxX} })
	return root, back, forth
}
