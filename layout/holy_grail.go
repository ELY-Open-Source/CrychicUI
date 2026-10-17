package layout

import "github.com/egoist/mygo/ui"

// HolyGrailSlots fill a HolyGrail; only Main is required.
type HolyGrailSlots struct {
	Header, Nav, Main, Aside, Footer func()
}

// HolyGrail lays out header, nav, main, aside and footer.
func HolyGrail(c *ui.Context, s HolyGrailSlots) ui.Element {
	if s.Main == nil {
		panic("layout: HolyGrail needs Main")
	}
	t := c.Theme()
	slot := func(key string, fn func()) ui.Element {
		e := ui.Column(c).Key(key).Gap(t.Space(2)).Padding(t.Space(3))
		return e.Transition(ui.ElementTransition{Enter: &ui.Motion{Collapse: true}, Exit: &ui.Motion{Collapse: true}}).
			Children(fn)
	}
	return ui.Column(c).Fill().Children(func() {
		if s.Header != nil {
			slot("header", s.Header).Background(t.Surface)
		}
		ui.Row(c).Grow(1).AlignItems(ui.Stretch).Children(func() {
			if s.Nav != nil {
				slot("nav", s.Nav).Width(t.Rem(12))
			}
			slot("main", s.Main).Grow(1).Basis(0)
			if s.Aside != nil {
				slot("aside", s.Aside).Width(t.Rem(12))
			}
		})
		if s.Footer != nil {
			slot("footer", s.Footer).Background(t.Surface)
		}
	})
}
